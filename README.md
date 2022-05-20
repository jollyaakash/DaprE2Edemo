# Dapr pub/sub with Edge MQTT broker + SPIFFE identities running as Dapr component

#### Pre-requirements installed
- Helm version 3.6.3
- Kubectl
- docker
- Linux AMD64 environment

In this quickstart, you'll first create a kubernetes cluster of your choice. It can be AKS cluster or on prem k3-s/minikube etc. Assuming you have cluster installed,
what we will demo here is 2 publisher microservice (generating orders on different topics and at different rate) and a subscriber microservice (processing both topics) 
in Go code running as an application. These will be running in harmony with Dapr sidecar where Edge MQTT+SPIFFE component is translating the calls to Edge MQTT broker. But the user pub and sub
are not aware of anything and just use Dapr go sdk to talk to side car ubiquitously. 

Visit [this](https://docs.dapr.io/developing-applications/building-blocks/pubsub/) link for more information about Dapr and Pub-Sub.

### Install SPIFFE server/agent plus MQTT broker on the cluster

```
helm repo add spiffecharts https://jollyaakash.github.io/SpiffeChart
helm repo update
helm install spiffe spiffecharts/spiffe-identity --version 0.4.4
```

This will install edge spiffe server(deployment)+agent(service) with MQTT broker tied to it and it understands the spiffe identities. By default there are 5 identities that
identity manager provides - agent, mqttbroker, publisher, bulkpublisher, subscriber. This can be configured though.

```
# To-uninstall
helm uninstall spiffe
kubectl delete namespace spiffe
```
### Install DAPR with MQTTE4K component bits

```
helm repo add dapr https://dapr.github.io/helm-charts/
helm repo update
export DAPR_RELEASE_NAME=dapr
export DAPR_REGISTRY=dmqttdemo.azurecr.io/dapr #[These contain MQTTE4K bits]
export DAPR_NAMESPACE=dapr-system
helm install $DAPR_RELEASE_NAME dapr/dapr --namespace=$DAPR_NAMESPACE --wait --timeout 5m0s --set global.ha.enabled=false --set-string global.tag=dev-linux-amd64 --set-string global.registry=$DAPR_REGISTRY --set global.logAsJson=true --set global.daprControlPlaneOs=linux --set global.daprControlPlaneArch=amd64 --set dapr_placement.logLevel=debug --set dapr_sidecar_injector.sidecarImagePullPolicy=Always --set global.imagePullPolicy=Always --set global.mtls.enabled=true --set dapr_placement.cluster.forceInMemoryLog=true
```

### Install DAPR Pub sub MQTTE4K component
Visit [this](https://docs.dapr.io/developing-applications/building-blocks/pubsub/howto-publish-subscribe/#step-1-setup-the-pubsub-component) link for more information about Dapr and Pub-Sub component.

```
kubectl apply -f ./k8s_deploy/postdapr/components/iot_pubsubmqtt.yaml
```

If you look closely at Yaml file. It tells Dapr I want to use a Component of type pubsub.mqtte4k, and provides it metadata. The metadata tells it to talk to mqttbroker service installed above in SPIFFE helm chart and tells the sidecar component that the Spiffe workload API will be mounted at "/home/nonroot/sockets/workloadapi.sock" in Sidecar VolumeMount.

*Note* - This doesn't spin up Dapr sidecar. It only tells Dapr that if any app wants to use Dapr sdk to talk to a pubsub named - orderpubsub. It should create it with the required configuration in the yaml file.


### Run Publisher and Subscriber apps which purely uses Go dapr sdk
These apps purely uses Dapr sdk ( can use http too ) to talk to Dapr sidecar component defined above and the component takes the heavy lifting to actually going and getting a SPIFFE identity, using it to authenticate with E4K MQTT Broker which understands the identity and work.

In this example we install following components -
- Checkout - Publisher creating orders
- BulkCheckout - Publisher creating bulk orders
- Order-Processor - Subscriber processing the orders

```
kubectl apply -f ./k8s_deploy/postdapr/pubsub/sub_deployment.yaml
kubectl apply -f ./k8s_deploy/postdapr/pubsub/pub_deployment.yaml
kubectl apply -f ./k8s_deploy/postdapr/pubsub/bulkpub_deployment.yaml
```

#### Check logs of publishers and subscriber to see functioning
Get pods name for each of them and in parallel terminal, do *kubectl logs pod_name publisher -f* or *kubectl logs pod_name subscriber -f* respectively for publisher or subscriber. 

If you want to see what the Application code for Publisher and Subscriber is doing or build it yourself.  Check [this](https://github.com/jollyaakash/DaprE2Edemo/blob/main/gocode/pub_sub/sdk/README.md)


### To uninstall everything

```
kubectl delete deployment dapr-subscriber # for removing order-processor
kubectl delete deployment dapr-publisher # for removing checkout
kubectl delete deployment dapr-bulkpublisher # for removing bulkcheckout
```

```
# Uninstall dapr and component crd
kubectl delete components.dapr.io daprpubsub
helm uninstall dapr
```

```
# Uninstall spiffe+mqttbroker
helm uninstall spiffe
```

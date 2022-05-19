# Dapr pub/sub with Edge MQTT broker + SPIFFE identities running as Dapr component

#### Pre-requirements installed
- Helm version 3.6.3
- Kubectl
- docker

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
### Install SPIFFE server/agent plus MQTT broker on the cluster


### To uninstall everything

```
kubectl delete deployment dapr-subscriber # for removing order-processor
kubectl delete deployment dapr-publisher # for removing checkout
kubectl delete deployment dapr-bulkpublisher # for removing bulkcheckout
```

```
# Uninstall spiffe+mqttbroker
helm uninstall spiffe
```

# Dapr pub/sub

In this quickstart, you'll create a publisher microservice and a subscriber microservice to demonstrate how Dapr enables a publish-subcribe pattern. The publisher generates messages of a specific topic, while subscribers listen for messages of specific topics.

Visit [this](https://docs.dapr.io/developing-applications/building-blocks/pubsub/) link for more information about Dapr and Pub-Sub.

> **Note:** This example leverages the Dapr client SDK. 

This quickstart includes two publisher:

- Go client message generator `checkout` 
- Go client message generator `bulkcheckout` 

And one subscriber: 
 
- Go subscriber `order-processor`

### Run Go message subscriber with Dapr

1. Navigate to the directory and build dependencies: 

```
cd ./order-processor
go build app.go
docker build .
#Tag docker image and push it to your required ACR
```

2. Run the Go subscriber app with Dapr: 

```
cd ./k8s_deploy/dapr_deploy/pubandusb/
Use the sub_deployment.yaml file to update the image with the correct image
```

<!-- END_STEP -->

### Run Go message publisher with Dapr

1. Navigate to the directory and install dependencies: 

```
cd ./checkout
go build app.go
docker build .
#Tag docker image and push it to your required ACR
```

```
cd ./bulkcheckout
go build app.go
docker build .
#Tag docker image and push it to your required ACR
```

```
cd ./k8s_deploy/dapr_deploy/pubandusb/
Use the pub_deployment.yaml file to update the image with the correct image for 'checkout'
Use the bulkpub_deployment.yaml file to update the image with the correct image for 'bulkcheckout'
```
#### To delete deployments/pods from cluster

```
kubectl delete deployment dapr-subscriber # for removing order-processor
kubectl delete deployment dapr-publisher # for removing checkout
kubectl delete deployment dapr-bulkpublisher # for removing bulkcheckout
```

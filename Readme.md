### Run mosquitto locally using this command

docker run -d -p 1883:1883 -p 9001:9001 -v /home/aakashjolly/Documents/mosquitto_dapr/mosquitto/config/mosquitto.conf:/mosquitto/config/mosquitto.conf -v /home/aakashjolly/Documents/mosquitto_dapr/mosquitto/data:/mosquitto/data -v /home/aakashjolly/Documents/mosquitto_dapr/mosquitto/log:/mosquitto/log eclipse-mosquitto

### Modify component and build it again

https://github.com/dapr/components-contrib/blob/master/docs/developing-component.md


cd $GOPATH/src

#### Clone dapr
mkdir -p github.com/dapr/dapr
git clone https://github.com/dapr/dapr.git github.com/dapr/dapr

#### Clone component-contrib
mkdir -p github.com/dapr/components-contrib
git clone https://github.com/dapr/components-contrib.git github.com/dapr/components-contrib


#### Add new component in relevant component directory Test it using
make test

#### Then go to directory path - 
....../go/src/github.com/dapr/dapr and modify module with local module

go mod edit -replace github.com/dapr/components-contrib=../components-contrib

#### Login to ACR and Build

az acr login -n loginto_your_choice_of_acr
export DAPR_REGISTRY=name.acr.io/dapr/
export DAPR_TAG=dev
make build-linux
make docker-build
make docker-deploy-k8s

### E2E publisher and subscriber -

https://github.com/dapr/quickstarts/tree/master/pub_sub/go/sdk

#### Subscriber
dapr run --app-port 6001 --app-id order-processor --app-protocol http --dapr-http-port 3501 -- go run app.go

#### Publisher
dapr run --app-id checkout --app-protocol http --dapr-http-port 3500 -- go run app.go


#### Build Dapr docker images locally and use them to deploy dapr in K8s environment 

https://github.com/dapr/dapr/blob/master/docs/development/developing-dapr.md


## To build single node broker -
cd repos/iotedge-broker/data-plane/

rustup target add x86_64-unknown-linux-musl
cargo build target=x86_64-unknown-linux-musl

Binary is present at repos/iotedge-broker/data-plane/target/x86_64-unknown-linux-musl/debug
## API definition
The API for the go client is defined by the api located in [this geo-grpc repo](https://github.com/geo-grpc/api). The Protobuf files and gRPC files are already compiled there. So you only need to import those files in your program as in this test and sample application.

## Install protoc

also, remember that you `GOROOT` and `GOPATH` need to be defined (https://stackoverflow.com/a/34896844/445372):
```bash
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin

go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```

Build protobuf
https://jbrandhorst.com/post/go-protobuf-tips/
```bash
protoc -I proto/ proto/epl/protobuf/geometry.proto --go_out=$GOPATH/src
protoc -I proto/ proto/epl/grpc/geometry_operators.proto --go_out=plugins=grpc:$GOPATH/src
```
## Testing againt geometry service
```bash
docker run -p 8980:8980 -d --name=temp-c echoparklabs/geometry-service-java:8-jre-slim
go test test/geometry_test.go -v


## Running example in Minikube
```bash
cd geometry-client-go
minikube start
eval $(minikube docker-env)
docker build -t go-client:latest .
kubectl create -f geometry-service.yml
kubectl create -f go-api.yml
minikube service geom-api --url
```

Then curl the output from the above `minikube service geom-api --url` command

testing:
https://github.com/2tvenom/go-test-teamcity


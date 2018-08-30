## Install protoc

also, remember that you `GOROOT` and `GOPATH` need to be defined (https://stackoverflow.com/a/34896844/445372):
```bash
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin

go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```

Build protobuf
```bash
protoc -I proto/ proto/epl/grpc/geometry/geometry_operators.proto --go_out=plugins=grpc:./
```



## Running example in Minikube
```bash
cd geometry-client-go
minikube start
docker build -t go-client:latest .
kubectl create -f geometry-service.yml
kubectl create -f geom-api.yml
minikube service geom-api --url
```

Then curl the output from the above `minikube service geom-api --url` command




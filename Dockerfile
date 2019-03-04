FROM golang:1.8

COPY ./ /go/src/github.com/geo-grpc/geometry-client-go
WORKDIR /go/src/github.com/geo-grpc/geometry-client-go

WORKDIR /go/src/github.com/geo-grpc/geometry-client-go/sample

RUN go get -d -v ./...

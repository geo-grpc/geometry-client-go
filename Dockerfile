FROM golang:1.12.0-stretch

WORKDIR /go/src/github.com/geo-grpc/geometry-client-go
COPY ./ /go/src/github.com/geo-grpc/geometry-client-go

WORKDIR /go/src/github.com/geo-grpc/geometry-client-go/sample

RUN go get -v ./...
RUN go install -v ./...

WORKDIR /go/src/github.com/geo-grpc/geometry-client-go
FROM golang:1.12.0-stretch

WORKDIR /go/src/github.com/geo-grpc/geometry-client-go
COPY ./ /go/src/github.com/geo-grpc/geometry-client-go

WORKDIR /go/src/github.com/geo-grpc/geometry-client-go/sample

ENV GO111MODULE=on
RUN go build

WORKDIR /go/src/github.com/geo-grpc/geometry-client-go
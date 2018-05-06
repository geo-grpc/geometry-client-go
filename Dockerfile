FROM golang:1.8

WORKDIR /go/src/geometry-client-go
COPY ./ ./

WORKDIR /go/src/geometry-client-go/sample

RUN go get -d -v ./...

RUN go build
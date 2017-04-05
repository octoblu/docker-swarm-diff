FROM golang:1.8
MAINTAINER Octoblu, Inc. <docker@octoblu.com>

WORKDIR /go/src/github.com/octoblu/docker-swarm-diff
COPY . /go/src/github.com/octoblu/docker-swarm-diff

RUN env CGO_ENABLED=0 go build -o docker-swarm-diff -a -ldflags '-s' .

CMD ["./docker-swarm-diff"]

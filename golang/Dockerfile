FROM golang:1.20

ENV GO111MODULE="on"
WORKDIR /go/src

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

FROM golang:1.20

ENV GO111MODULE="on"
WORKDIR /go/src

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
# RUN go install github.com/cosmtrek/air@v1.42.0 # curl経由でerrorが発生するなら固定verでインストールする
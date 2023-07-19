FROM golang:1.20

ENV GO111MODULE="on"  
RUN mkdir -p /go/src
WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download 
COPY . /go/src
RUN go build -o app
CMD ./app

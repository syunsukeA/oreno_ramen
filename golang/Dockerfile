FROM golang:1.20

RUN mkdir -p /go/src
WORKDIR /go/src
COPY . /go/src
RUN go build -o app
CMD ./app

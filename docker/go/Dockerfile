FROM golang:1.9.4-alpine3.7
WORKDIR /go/src/github.com/PumpkinSeed/incrmntr

ENV GO15VENDOREXPERIMENT=0

ADD ./ /go/src/github.com/PumpkinSeed/incrmntr
RUN apk add --no-cache gcc musl-dev
RUN go install github.com/PumpkinSeed/incrmntr

ENTRYPOINT go run examples/api/main.go
FROM golang:1.10.1-alpine3.7

RUN apk update
RUN apk add git

WORKDIR /go/src/github.com/therealfakemoot/genesis
ADD . .
RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT /go/bin/genesis

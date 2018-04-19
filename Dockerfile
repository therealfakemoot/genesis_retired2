FROM alpine:latest

RUN apk update
RUN apk add go
RUN apk add git

WORKDIR /go/src/github.com/therealfakemoot/genesis
ADD . .
RUN go get -d -v ./...
RUN go install -v ./...

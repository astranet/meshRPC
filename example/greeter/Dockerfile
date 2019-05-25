FROM golang:1.12-alpine as builder

RUN apk add --no-cache git

ADD . /gopath/src/github.com/astranet/meshRPC
ENV GOPATH=/gopath
RUN go get github.com/astranet/meshRPC/example/greeter

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /gopath/bin/greeter /usr/local/bin/

EXPOSE 11999
ENTRYPOINT ["greeter"]

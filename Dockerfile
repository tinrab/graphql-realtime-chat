FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make git ca-certificates
RUN go get -u -v golang.org/x/vgo

WORKDIR /go/src/github.com/tinrab/graphql-realtime-chat
COPY go.mod .
COPY vendor vendor
COPY server server
COPY main.go .
RUN vgo build -o /go/bin/app .

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]

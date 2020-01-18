FROM golang:1.13-alpine3.11 AS build
RUN apk --no-cache add clang gcc g++ make git ca-certificates

WORKDIR /go/src/github.com/tinrab/graphql-realtime-chat
COPY go.mod go.sum vendor main.go ./
COPY server server
RUN go build -o /go/bin/app .

FROM alpine:3.11
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]

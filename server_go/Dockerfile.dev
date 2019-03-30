FROM golang:alpine AS build-env
ADD src /go/src
RUN apk update && apk add git
RUN go get github.com/lib/pq
RUN go get -u github.com/go-redis/redis
RUN go build -o bin/httpserver httpserver

FROM alpine
WORKDIR /app
COPY --from=build-env /go/bin/httpserver /app
ENTRYPOINT ./httpserver 
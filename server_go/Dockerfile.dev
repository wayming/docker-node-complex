FROM golang:alpine AS build-env
ADD src  /go/src
RUN go build -o bin/httpserver httpserver

FROM alpine
WORKDIR /app
COPY --from=build-env /go/bin/httpserver /app
ENTRYPOINT ./httpserver 
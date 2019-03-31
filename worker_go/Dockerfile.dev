FROM golang:alpine AS build-env
ADD src /go/src
RUN apk update && apk add git
RUN go get -u github.com/go-redis/redis
RUN go build -o /go/bin/fibcalc fibcalc

FROM alpine
WORKDIR /app
COPY --from=build-env /go/bin/fibcalc /app
ENTRYPOINT ./fibcalc
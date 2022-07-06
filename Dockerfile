# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
# Download necessary Go modules
RUN go mod download

COPY basic/*.go ./

RUN go build -o /docker-gs-ping

##
## Deploy
##
FROM alpine:3.10

WORKDIR /

COPY --from=build /docker-gs-ping /docker-gs-ping

EXPOSE 8080

ENTRYPOINT ["/docker-gs-ping"]
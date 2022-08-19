# syntax=docker/dockerfile:1
FROM golang:1.18-alpine

RUN mkdir /var/tmp/file_storage

RUN mkdir /app
WORKDIR /app

COPY go.mod .
# COPY go.sum .
RUN go mod download

COPY . .

RUN go build /app/cmd/http/main.go
RUN mv /app/main /usr/local/bin/

EXPOSE 8080

CMD export $(cat /app/.docker_env) && main

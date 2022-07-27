# syntax=docker/dockerfile:1
FROM golang:1.18-alpine

RUN mkdir /app
WORKDIR /app

COPY go.mod /app
# COPY go.sum .
RUN go mod download

COPY . /app

RUN go build /app/cmd/main.go
RUN mv /app/main /usr/local/bin/

EXPOSE 8080

CMD ["main"]

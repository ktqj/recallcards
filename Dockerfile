FROM golang:1.18-alpine as build
WORKDIR /go/src/app/
COPY . /go/src/app/
RUN go generate ./... && go mod download && go build -o /bin/recallcards_web ./cmd/http/main.go

##################

FROM alpine:3.6
EXPOSE 8080

ENV JSON_STORAGE_DIR /file_storage
ENV TEMPLATES_DIR /templates

RUN mkdir ${JSON_STORAGE_DIR}
RUN mkdir ${TEMPLATES_DIR}

COPY /pkg/web/templates ${TEMPLATES_DIR}
COPY --from=build /bin/recallcards_web /recallcards_web

CMD ["/recallcards_web"]



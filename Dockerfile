# syntax=docker/dockerfile:1
FROM scratch
ADD ./bin/sample_webapp sample_webapp
EXPOSE 8080
ENTRYPOINT ["/sample_webapp"]
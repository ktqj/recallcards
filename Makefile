compile:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/sample_webapp ./cmd/main.go 

build-image:
	docker build -t sample_webapp .

docker-run:
	docker run -p 6060:8080 --name localtest --rm sample_webapp
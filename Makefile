compile:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/recallcards ./cmd/main.go 

build-image:
	docker build -t recallcards .

docker-run: build-image
	docker run -p 6060:8080 --name localtest --rm recallcards
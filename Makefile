include .env

compile:
	go build -o ./bin/recallcards ./cmd/http/main.go

run: compile
	./bin/recallcards

build-image:
	docker build -t recallcards .

docker-run: build-image
	docker run --init -p 6060:8080 -v $$(pwd)/tmp/mem_storage/:/var/tmp/mem_storage --name localtest --rm recallcards

create-tmp:
	mkdir -p ./tmp/mem_storage

up: create-tmp build-image
	docker-compose up

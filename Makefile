include ./.env
export

create-tmp:
	mkdir -p ./tmp/file_storage

clean:
	go fmt ./...

lint:
	golangci-lint run -v

################

compile-http:
	go build -o ./bin/recallcards_web ./cmd/http/main.go

web: create-tmp compile-http
	./bin/recallcards_web

################

compile-cli:
	go build -o ./bin/recallcards_cli ./cmd/cli/cli.go

cli: create-tmp compile-cli
	./bin/recallcards_cli

#################

compile-debug:
	go build -o ./bin/recallcards_debug ./cmd/debug/debug.go

debug: create-tmp compile-debug
	./bin/recallcards_debug

#################
test:
	go test -v -count=1  ./...

################

build-image:
	docker build -t recallcards .

docker-run: build-image
	docker run --init -p 6060:8080 -v $$(pwd)/tmp/file_storage/:/var/tmp/file_storage --name localtest --rm recallcards

up: create-tmp build-image
	docker-compose up

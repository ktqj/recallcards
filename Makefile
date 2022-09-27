include ./.env
export

create-tmp:
	mkdir -p ./tmp/file_storage

clean:
	go fmt ./...

lint:
	golangci-lint run -c ./golangci-lint.yaml ./...

################

compile-http:
	go generate ./...
	go build -o ./bin/recallcards_web ./cmd/http/main.go

web: create-tmp compile-http
	./bin/recallcards_web

################

compile-cli:
	go generate ./...
	go build -o ./bin/recallcards_cli ./cmd/cli/cli.go

cli: create-tmp compile-cli
	./bin/recallcards_cli

#################

compile-debug:
	go generate ./...
	go build -o ./bin/recallcards_debug ./cmd/debug/debug.go

debug: create-tmp compile-debug
	./bin/recallcards_debug

#################
test:
	go test -v --cover -count=1  ./...

bench:
	go test -v -bench=. -benchmem ./...

################

build-web-image:
	docker build -t recallcards_web .

docker-run-web: build-web-image
	docker run --init -p 6060:8080 -v $$(pwd)/tmp/file_storage/:/file_storage --name localtest --rm recallcards_web

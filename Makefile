include .env

# set up
create-tmp:
	mkdir -p ./tmp/mem_storage

# local work
compile:
	go build -o ./bin/recallcards ./cmd/http/main.go

run: create-tmp compile
	./bin/recallcards

test:
	go test -v -count=1  ./...

# docker work
build-image:
	docker build -t recallcards .

docker-run: build-image
	docker run --init -p 6060:8080 -v $$(pwd)/tmp/mem_storage/:/var/tmp/mem_storage --name localtest --rm recallcards



up: create-tmp build-image
	docker-compose up

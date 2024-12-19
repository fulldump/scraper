

run:
	go run .

build:
	go build -o ./bin/scraper .

.PHONY: deps
deps:
	go get -t -u ./...
	go mod tidy
	go mod vendor


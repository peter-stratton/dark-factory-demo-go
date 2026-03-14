.PHONY: build test lint run clean

build:
	go build -o bin/bookmarks-api .

test:
	go test ./...

lint:
	golangci-lint run

run:
	go run .

clean:
	rm -rf bin/

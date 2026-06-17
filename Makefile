.PHONY: run test build tidy lint

run:
	go run main.go

build:
	go build -o bin/coach main.go

test:
	go test ./... -v -count=1

tidy:
	go mod tidy

lint:
	gofmt -w .
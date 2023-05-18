.PHONY: build
build:
	go build -o anti-brute-force main.go

port ?= 8080

.PHONY: run
run:
	go run main.go -p ${port}

.PHONY: test
test:
	go test -race ./...


.PHONY: lint
lint:
	golangci-lint run
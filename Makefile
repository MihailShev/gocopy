default: dev

dev:
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
.PHONY: dev

test:
	go test -v ./cmd/gocopy
.PHONY: test

lint:
	golangci-lint run
.PHONY: lint

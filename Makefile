default: install

install:
	go install ./src/gocopy.go
.PHONY: install

test:
	go test -v ./src
.PHONY: test
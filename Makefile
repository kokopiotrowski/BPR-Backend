BINARY_NAME=stockx
LINT_BIN := $(GOPATH)/bin/golangci-lint

build:
	@$(LINT_BIN) run
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go

run:
	./${BINARY_NAME}-linux

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows
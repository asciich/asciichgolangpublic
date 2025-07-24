REPO_ROOT := $(shell git rev-parse --show-toplevel)
BIN_DIR := ${REPO_ROOT}/bin

PHONY: test
test:
	go test -timeout 20m ./... -v

PHONY: tidy
tidy:
	go mod tidy

PHONY: build
build:
	@echo "Build all packages"
	go build ./...

	@echo "Build asciichgolangpublic"
	@mkdir -p "$(BIN_DIR)"
	go build -o "$(BIN_DIR)/asciichgolangpublic" cmd/asciichgolangpublic/asciichgolangpublic.go

	@echo "Build finished"
	
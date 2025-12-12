REPO_ROOT := $(shell git rev-parse --show-toplevel)
BIN_DIR := ${REPO_ROOT}/bin
BIN_PATH := $(BIN_DIR)/asciichgolangpublic

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

	@echo "Build $(BIN_PATH)"
	@mkdir -p "$(BIN_DIR)"
	go build -o "$(BIN_PATH)" cmd/asciichgolangpublic/asciichgolangpublic.go

	@echo "Build finished"
	
PHONY: install
install: build
	@echo "Install started."

	$(BIN_PATH) install --verbose --binary-name=asciichgolangpublic

	@echo "Install finished."

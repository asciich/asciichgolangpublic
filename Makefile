REPO_ROOT := $(shell git rev-parse --show-toplevel)
BIN_DIR := ${REPO_ROOT}/bin
BIN_PATH := $(BIN_DIR)/asciichgolangpublic

# UID and GID of the current user.
# If make is started with sudo, use the user who called sudo instead of root.
USER_UID := $(or $(SUDO_UID),$(shell id -u))
USER_GID := $(or $(SUDO_GID),$(shell id -g))

.PHONY: test
test:
	go test -timeout 20m ./... -v

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build:
	@echo "Build all packages"
	go build ./...

	@echo "Build $(BIN_PATH)"
	@mkdir -p "$(BIN_DIR)"
	CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o "$(BIN_PATH)" cmd/asciichgolangpublic/main.go

	@echo "Build finished"

.PHONY: install
install: build
	@echo "Install started."

	$(BIN_PATH) install --verbose --binary-name=asciichgolangpublic

	@echo "Install finished."

.PHONY: fix_access_rights
fix_access_rights:
	@echo "Fix access rights started for UID '$(USER_UID)' and GID '$(USER_GID)'."

	sudo chown -R $(USER_UID):$(USER_GID) .
	go clean -cache -testcache
	git ls-files -z '*.go' '*.yml' '*.yaml' '*.md' '*.html' '*.txt' | xargs -0 chmod 644

	@echo "Fix access rights finished."

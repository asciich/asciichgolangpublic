# ----------------------------------------------------------------------
#   Makefile – Build, vet, test, install, and housekeeping
# ----------------------------------------------------------------------
REPO_ROOT := $(shell git rev-parse --show-toplevel)
BIN_DIR   := ${REPO_ROOT}/bin
BIN_PATH  := $(BIN_DIR)/asciichgolangpublic

# UID/GID of the user that invoked make (handles sudo correctly)
USER_UID := $(or $(SUDO_UID),$(shell id -u))
USER_GID := $(or $(SUDO_GID),$(shell id -g))

# ----------------------------------------------------------------------
#   High‑level targets
# ----------------------------------------------------------------------
.PHONY: test
test:
	@echo "Running tests..."
	go test -timeout 20m ./... -v

.PHONY: tidy
tidy:
	@echo "Tidying go.mod…"
	go mod tidy

# ----------------------------------------------------------------------
#   Vet – type‑check all packages (source + test files)
# ----------------------------------------------------------------------
.PHONY: vet
vet:
	@echo "=== Vet all packages (source + test files) ==="
	go vet ./...

# ----------------------------------------------------------------------
#   Build – compile sources and the main binary
# ----------------------------------------------------------------------
.PHONY: build
build: vet
	@echo "=== Build source packages (normal binary) ==="
	go build ./...

	@echo "=== Build the main binary ==="
	@mkdir -p "$(BIN_DIR)"
	@CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" \
		-o "$(BIN_PATH)" cmd/asciichgolangpublic/main.go
	@echo "Build finished."

# ----------------------------------------------------------------------
#   Install – install the built binary
# ----------------------------------------------------------------------
.PHONY: install
install: build
	@echo "Install started."
	$(BIN_PATH) install --verbose --binary-name=asciichgolangpublic
	@echo "Install finished."

# ----------------------------------------------------------------------
#   Fix access rights – adjust ownership and permissions
# ----------------------------------------------------------------------
.PHONY: fix_access_rights
fix_access_rights:
	@echo "Fix access rights started for UID '$(USER_UID)' and GID '$(USER_GID)'."
	sudo chown -R $(USER_UID):$(USER_GID) .
	go clean -cache -testcache
	git ls-files -z '*.go' '*.yml' '*.yaml' '*.md' '*.html' '*.txt' | \
		xargs -0 chmod 644
	@echo "Fix access rights finished."
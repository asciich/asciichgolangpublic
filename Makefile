PHONY: test
test:
	go test -timeout 20m ./... -v

PHONY: tidy
tidy:
	go mod tidy

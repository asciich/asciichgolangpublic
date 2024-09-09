PHONY: test
test:
	go test -timeout 20m ./... -v

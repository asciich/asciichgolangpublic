PHONY: test
test:
	go test -timeout 15m ./... -v

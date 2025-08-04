.PHONY: test
test:
	go test ./...  -v  -coverpkg=./... -coverprofile profile.cov
.PHONY: lint
lint:
	golangci-lint run

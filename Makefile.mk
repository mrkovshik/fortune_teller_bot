test:
	go test ./...  -v  -coverpkg=./... -coverprofile profile.cov

lint:
	golangci-lint run

lint:
	@go fmt ./...
	@go vet ./...

test:
	@go test ./...
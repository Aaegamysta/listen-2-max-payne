run:
	go run ./cmd/listen2maxpayne

lint:
	gofumpt -w .
	golangci-lint cache clean
	golangci-lint run ./...

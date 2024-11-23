build:
	@go build

run: build
	@./bin/fs

test:
	@go test ./...
BIN_DIR=bin
BIN_ARTIFACT=fs
BIN_PATH=$(BIN_DIR)/$(BIN_ARTIFACT)
PEER1_OPTIONS=--listen=:8000 --api-listen=:9000 --db=/tmp/peer1.db --file-storage-path=/tmp/peer1_fs_storage
PEER2_OPTIONS=--listen=:8001 --api-listen=:9001 --db=/tmp/peer2.db --file-storage-path=/tmp/peer2_fs_storage --bootstrap=localhost:8000 --test-storage=true

build:
	@go build -o $(BIN_PATH)

run-with-defaults: build
	$(BIN_PATH)

run-peer-1: build
	$(BIN_PATH) $(PEER1_OPTIONS)

run-peer-2: build
	$(BIN_PATH) $(PEER2_OPTIONS)

test:
	@go test ./...
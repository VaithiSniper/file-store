BIN_DIR=bin
BIN_ARTIFACT=fs
BIN_PATH=$(BIN_DIR)/$(BIN_ARTIFACT)

WILDCARD_IP := "0.0.0.0"
LOCAL_IP := $(shell ip route get 1 | awk '{print $$7; exit}')

LOG_LEVEL=info

PEER1_OPTIONS=--name=alpha-node-1 --listen=$(LOCAL_IP):8000 --api-listen=$(WILDCARD_IP):9000 --db=/tmp/peer1.db --file-storage-path=/tmp/peer1_fs_storage --log-level=$(LOG_LEVEL)
PEER2_OPTIONS=--name=beta-node-1 --listen=$(LOCAL_IP):8001 --api-listen=$(WILDCARD_IP):9001 --db=/tmp/peer2.db --file-storage-path=/tmp/peer2_fs_storage --bootstrap=$(LOCAL_IP):8000 --test-storage=true --log-level=$(LOG_LEVEL)

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
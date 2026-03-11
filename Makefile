# Ensure protoc finds Go plugins (protoc-gen-go, protoc-gen-go-grpc)
GOBIN := $(shell go env GOPATH)/bin
export PATH := $(GOBIN):$(PATH)

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=module=grpc_starbuckscoffee \
		--go-grpc_out=. --go-grpc_opt=module=grpc_starbuckscoffee \
		proto/coffeeshop.proto

# Seed Firestore (emulator). Start emulators first: firebase emulators:start
.PHONY: seed
seed:
	FIRESTORE_EMULATOR_HOST=localhost:8080 GOOGLE_CLOUD_PROJECT=grpc-starbucks-coffee go run ./cmd/seed

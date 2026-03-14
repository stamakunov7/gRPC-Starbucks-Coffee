# Ensure protoc finds Go plugins (protoc-gen-go, protoc-gen-go-grpc)
GOBIN := $(shell go env GOPATH)/bin
export PATH := $(GOBIN):$(PATH)

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=module=grpc_starbuckscoffee \
		--go-grpc_out=. --go-grpc_opt=module=grpc_starbuckscoffee \
		proto/coffeeshop.proto

# Start Firestore emulator with persistence (data saved to emulator-data/ on exit, restored on start).
# Run in one terminal and leave it running. First time: run "make seed" in another terminal once.
.PHONY: emulators
emulators:
	firebase emulators:start --only firestore,ui --import=./emulator-data --export-on-exit=./emulator-data

# Seed Firestore emulator (categories + drinks). Run once after first "make emulators", or when you want to reset data.
.PHONY: seed
seed:
	FIRESTORE_EMULATOR_HOST=localhost:8080 GOOGLE_CLOUD_PROJECT=grpc-starbucks-coffee go run ./cmd/seed

# Kill any process listening on port 9001 (old gRPC server). Run if "make run-server" says address already in use.
.PHONY: kill-server
kill-server:
	@lsof -ti :9001 | xargs kill -9 2>/dev/null || true

# Run gRPC server (uses Firestore emulator). Stops any existing server on :9001 first so you always get the latest code.
.PHONY: run-server
run-server: kill-server
	FIRESTORE_EMULATOR_HOST=localhost:8080 GOOGLE_CLOUD_PROJECT=grpc-starbucks-coffee go run ./server

# Web UI (v0-generated Next.js app). First time: make install-ui; then make run-ui
.PHONY: install-ui
install-ui:
	cd web && npm install

.PHONY: run-ui
run-ui:
	cd web && npm run dev

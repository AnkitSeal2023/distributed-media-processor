API_DIR := services/api
API_BIN := $(API_DIR)/api

AUTH_DIR := services/auth
AUTH_BIN := $(AUTH_DIR)/auth

UPLOAD_DIR := services/upload
UPLOAD_BIN := $(UPLOAD_DIR)/upload

VIDEO_STATUS_DIR := services/upload/status
VIDEO_STATUS_BIN := $(VIDEO_STATUS_DIR)/video_status

WORKER_DIR := services/worker
WORKER_BIN := $(WORKER_DIR)/worker

.PHONY: api
api:
	go build -o $(API_BIN) ./$(API_DIR)

.PHONY: auth
auth:
	go build -o $(AUTH_BIN) ./$(AUTH_DIR)

.PHONY: upload
upload:
	go build -o $(UPLOAD_BIN) ./$(UPLOAD_DIR)

.PHONY: video_status
video_status:
	go build -o $(VIDEO_STATUS_BIN) ./$(VIDEO_STATUS_DIR)

.PHONY: worker
worker:
	go build -o $(WORKER_BIN) ./$(WORKER_DIR)

auth-run: auth
	doppler run -p auth_service -c dev --command="./$(AUTH_BIN)"

api-run: api
	doppler run -p api_service -c dev --command="./$(API_BIN)"

upload-run: upload
	doppler run -p upload_service -c dev --command="./$(UPLOAD_BIN)"

worker-run: worker
	doppler run -p worker_service -c dev --command="./$(WORKER_BIN)"

video_status-run: video_status
	doppler run -p upload_service -c dev --command="./$(VIDEO_STATUS_BIN)"

auth-gen-proto:
	protoc --go_out=proto/generated --go_opt=paths=source_relative --go-grpc_out=proto/generated --go-grpc_opt=paths=source_relative proto/auth/v1/auth.proto

upload-gen-proto:
	protoc --go_out=proto/generated --go_opt=paths=source_relative --go-grpc_out=proto/generated --go-grpc_opt=paths=source_relative proto/upload/v1/upload.proto

API_DIR := services/api
API_BIN := $(API_DIR)/api
AUTH_DIR := services/auth
AUTH_BIN := $(AUTH_DIR)/auth

.PHONY: api
api:
	go build -o $(API_BIN) ./$(API_DIR)

.PHONY: auth
auth:
	go build -o $(AUTH_BIN) ./$(AUTH_DIR)

auth-run: auth
	./$(AUTH_BIN)

api-run: api
	./$(API_BIN)

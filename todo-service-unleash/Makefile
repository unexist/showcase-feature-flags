.DEFAULT_GOAL := build-gin
.ONESHELL:

EXE_GIN := todo-service.bin
API_URL := http://localhost:4242/api

--guard-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set: $*=abc123 make $(MAKECMDGOALS)"; \
		exit 1; \
	fi

swagger:
	@$(SHELL) -c "swag init"

# Build
build:
	@$(SHELL) -c "GO111MODULE=on; go mod download; go build -o $(EXE_GIN)"

# Analysis
vet:
	@$(SHELL) -c "go vet"

# Run
run: --guard-API_TOKEN
	@$(SHELL) -c "API_URL=$(API_URL) API_TOKEN=$(API_TOKEN) ./$(EXE_GIN)"

# Tests
test: --guard-API_TOKEN
	@$(SHELL) -c "API_URL=$(API_URL) API_TOKEN=$(API_TOKEN) go test -count=1 -v ./test"

# Helper
clear:
	rm -rf ./$(EXE_GIN)

install:
	go get github.com/Unleash/unleash-client-go/v3
	go install github.com/swaggo/swag/cmd/swag@latest
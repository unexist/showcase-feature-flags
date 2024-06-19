.DEFAULT_GOAL := build
.ONESHELL:

EXE_GIN := todo-service.bin

define JSON_TODO
curl -X 'POST' \
  'http://localhost:8080/todo' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "string",
  "done": true,
  "title": "string"
}'
endef
export JSON_TODO
# Tools
todo:
	@echo $$JSON_TODO | bash

list:
	@curl -X 'GET' 'http://localhost:8080/todo' -H 'accept: */*' | jq .

swagger:
	@$(SHELL) -c "cd todo-service-gin; swag init"

open-swagger:
	open http://localhost:8080/swagger/index.html

c4:
	@$(SHELL) -c "cd tools; ./$(EXE_STRUCT) && plantuml *.plantuml"

# Build
build-gin:
	@$(SHELL) -c "GO111MODULE=on; go mod download; go build -o $(EXE_GIN)"

# Analysis
vet-gin:
	@$(SHELL) -c "go vet"

# Run
run-gin:
	@$(SHELL) -c "./$(EXE_GIN)"

# Tests
test-gin:
	@$(SHELL) -c "go test -v ./test"

# Helper
clear:
	rm -rf ./$(EXE_GIN)

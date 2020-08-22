.PHONY: build

build:
	@echo "Building API..."
	@go build -o ./pos-api ./api

run:
	@echo "Running API..."
	@./pos-api

start: build run
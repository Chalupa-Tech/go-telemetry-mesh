.PHONY: build test lint clean run

APP_NAME := $(shell basename $(shell pwd))
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")

build:
	@echo "Building $(APP_NAME)..."
	go build -v ./...

test:
	@echo "Running tests..."
	go test -v -race ./...

lint:
	@echo "Linting..."
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Skipping."; \
	fi

clean:
	@echo "Cleaning..."
	go clean
	rm -f $(APP_NAME)

run: build
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME)

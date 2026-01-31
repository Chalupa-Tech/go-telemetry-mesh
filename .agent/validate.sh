#!/bin/bash
set -e

echo "Running Validation..."

echo "1. Formatting..."
go fmt ./...

echo "2. Linting..."
go vet ./...

echo "3. Testing..."
go test ./...

if [ -d "vendor" ]; then
    echo "4. Verifying Vendor..."
    go mod verify
fi

echo "Validation Passed!"

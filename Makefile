# Makefile for building the usqlmcp project

BINARY_NAME := bin/usqlmcp

# Default target
all: test

build.all:
	./build.sh -b -t all

build.most:
	@echo "Building $(BINARY_NAME)..."
	./build.sh -b -t most

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

# Run the application
run: build.most
	./$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test -v ./... | tee test.log | less -R

.PHONY: all build clean run test

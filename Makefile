# Makefile for Transfa Backend Microservices
#
# This Makefile provides a set of commands to manage the development lifecycle of all
# Transfa backend microservices located in the `services/` directory. It automates
# common tasks such as linting, building, and cleaning up artifacts.
#
# Key Targets:
# - `lint`: Formats all Go code with `gofmt` and runs static analysis with `go vet`.
# - `build`: Compiles each microservice into a production-ready binary in its respective `bin/` directory.
# - `clean`: Removes all build artifacts.
#
# This centralized script ensures that all services are treated uniformly, enforcing
# code quality and build consistency across the entire project.

# Discover all services by looking for directories inside the 'services/' folder.
SERVICES := $(wildcard services/*)

# Default target runs linting and then building.
.PHONY: all
all: lint build

# Lint all Go services. It iterates through each service directory,
# formats the code, and runs go vet to catch suspicious constructs.
# go vet is the modern standard for static analysis in Go.
.PHONY: lint
lint:
	@echo "Linting all services..."
	@for service in $(SERVICES); do \
		echo "--> Linting $$service"; \
		(cd $$service && go fmt ./... && go vet ./...); \
	done
	@echo "Linting complete."

# Build all Go services. It creates a `bin` directory within each service
# and compiles the main.go file into a binary named after the service.
.PHONY: build
build:
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "--> Building $$service"; \
		(cd $$service && mkdir -p bin && go build -v -o ./bin/$$(basename $$service) ./cmd/main.go); \
	done
	@echo "Build complete."

# Clean up all build artifacts by removing the `bin` directory from each service.
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@for service in $(SERVICES); do \
		echo "--> Cleaning $$service"; \
		(rm -rf $$service/bin); \
	done
	@echo "Cleaning complete."
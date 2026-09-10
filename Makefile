BINARY_NAME=dora
API_BINARY_NAME=dora-api
VERSION=v0.1.0
BUILD_DIR=bin

.PHONY: all build test vet lint run clean install uninstall validate lab-up lab-down release

all: build

build:
	@echo "Building DORA binaries..."
	mkdir -p $(BUILD_DIR)
	go build -ldflags="-X 'github.com/dora/dora/internal/configuration.Version=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/dora
	go build -ldflags="-X 'github.com/dora/dora/internal/configuration.Version=$(VERSION)'" -o $(BUILD_DIR)/$(API_BINARY_NAME) ./cmd/dora-api

test:
	@echo "Running tests..."
	go test -v -race ./...

vet:
	@echo "Running go vet..."
	go vet ./...

lint:
	@echo "Running static checks..."
	go vet ./...

run: build
	./$(BUILD_DIR)/$(BINARY_NAME) -help

clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME) $(API_BINARY_NAME)

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	install -m 0755 $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin..."
	rm -f /usr/local/bin/$(BINARY_NAME)

validate: vet test

lab-up:
	@echo "Starting local laboratory environment..."
	docker compose -f docker-compose.yml up -d --build

lab-down:
	@echo "Stopping local laboratory environment..."
	docker compose -f docker-compose.yml down

release:
	@echo "Building release binaries for all target architectures..."
	mkdir -p $(BUILD_DIR)/release
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 ./cmd/dora
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-arm64 ./cmd/dora
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 ./cmd/dora
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-arm64 ./cmd/dora
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-windows-amd64.exe ./cmd/dora

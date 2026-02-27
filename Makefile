BINARY_NAME=egg

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

BUILD_DIR=.

MAIN_PATH=.

VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Config values are baked into the binary at build time.
# Set these env vars before running make build for a release binary:
# API_ENDPOINT, COGNITO_USER_POOL_ID, COGNITO_CLIENT_ID, COGNITO_DOMAIN, COGNITO_REGION
LDFLAGS=-ldflags "\
  -X main.Version=$(VERSION) \
  -X main.BuildTime=$(BUILD_TIME) \
  -X main.GitCommit=$(GIT_COMMIT) \
  -X github.com/owenHochwald/egg-carton/cli/config.CompiledAPIEndpoint=$(API_ENDPOINT) \
  -X github.com/owenHochwald/egg-carton/cli/config.CompiledUserPoolID=$(COGNITO_USER_POOL_ID) \
  -X github.com/owenHochwald/egg-carton/cli/config.CompiledClientID=$(COGNITO_CLIENT_ID) \
  -X github.com/owenHochwald/egg-carton/cli/config.CompiledCognitoDomain=$(COGNITO_DOMAIN) \
  -X github.com/owenHochwald/egg-carton/cli/config.CompiledRegion=$(COGNITO_REGION)"

.PHONY: all build clean test install uninstall fmt vet help

all: clean build ## Build the binary

build: ## Build the egg binary
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

build-linux: ## Build for Linux
	@echo "Building $(BINARY_NAME) for Linux..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)

build-windows: ## Build for Windows
	@echo "Building $(BINARY_NAME) for Windows..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_PATH)

build-mac: ## Build for macOS
	@echo "Building $(BINARY_NAME) for macOS..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-mac-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-mac-arm64 $(MAIN_PATH)

build-all: build-linux build-windows build-mac ## Build for all platforms

clean: ## Remove build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f $(BUILD_DIR)/$(BINARY_NAME)-*
	rm -f $(BUILD_DIR)/$(BINARY_NAME).exe
	@echo "✅ Clean complete"
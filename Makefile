# =============================================================================
# MAKEFILE FOR GO APPLICATION - BOOKMARK SERVICE
# =============================================================================
# This Makefile provides automation for common development tasks:
#   - Local development (start, test, generate)
#   - Docker operations (build, test, release)
#   - Code quality (coverage thresholds)
#
# QUICK START:
#   make start         - Run the application locally
#   make test          - Run tests with coverage locally
#   make docker-build  - Build production Docker image
#   make docker-test   - Run tests inside Docker container
#
# CI/CD WORKFLOW:
#   1. make docker-login  - Authenticate with Docker Hub
#   2. make docker-test   - Run tests and check coverage threshold
#   3. make docker-build  - Build the production image
#   4. make docker-release - Push image to Docker Hub
# =============================================================================

# =============================================================================
# DOCKER IMAGE CONFIGURATION
# =============================================================================
# IMG_NAME: Docker Hub repository name (format: username/repository)
# This is where the built images will be pushed.
IMG_NAME=johnnyho3820/ebvn-bookmark-repo

# -----------------------------------------------------------------------------
# AUTOMATIC IMAGE TAGGING STRATEGY
# -----------------------------------------------------------------------------
# The IMG_TAG is automatically determined based on Git state:
#
# Priority order (highest to lowest):
#   1. Git tag (e.g., v1.0.0) → Uses the exact tag as image tag
#   2. master branch           → Tags image as "dev"
#   3. Other branches        → IMG_TAG may be empty (handle in CI)
#
# Examples:
#   - On tag v1.2.3:        IMG_TAG=v1.2.3
#   - On master branch:       IMG_TAG=dev
#   - On feature/xyz:       IMG_TAG=(empty, may need handling)
# -----------------------------------------------------------------------------

# Attempt to get the exact Git tag at HEAD (empty if not on a tag)
# The 2>/dev/null suppresses error messages when not on a tag
GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null)

# Get the current Git branch name
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

# Default: If on master branch, use "dev" as the image tag
# This enables continuous deployment to a dev/staging environment
ifeq ($(BRANCH),master)
	IMG_TAG := dev
endif

# Override: If we're on an exact Git tag, use that tag for the image
# Git tags take precedence over branch-based tagging
# This enables versioned releases (e.g., v1.0.0)
ifneq ($(GIT_TAG),)
	IMG_TAG := $(GIT_TAG)
endif

# Export IMG_TAG so it's available to child processes and scripts
export IMG_TAG

# =============================================================================
# PHONY TARGETS DECLARATION
# =============================================================================
# Declare targets that don't correspond to actual files.
# This prevents conflicts with files of the same name and improves performance.
.PHONY: start swagger-gen test gen-all docker-test docker-build docker-release docker-login clean

# =============================================================================
# LOCAL DEVELOPMENT TARGETS
# =============================================================================

# -----------------------------------------------------------------------------
# start: Run the application locally
# -----------------------------------------------------------------------------
# Starts the Go application in development mode.
# The server will be accessible at http://localhost:<port>
#
# Usage: make start
# -----------------------------------------------------------------------------
start:
	go run cmd/api/main.go

# -----------------------------------------------------------------------------
# swagger-gen: Generate Swagger/OpenAPI documentation
# -----------------------------------------------------------------------------
# Scans the codebase for Swagger annotations and generates:
#   - docs/swagger.json
#   - docs/swagger.yaml
#   - docs/docs.go (Go bindings)
#
# Flags:
#   -g cmd/api/main.go   Entry point for scanning
#   --parseDependency    Include types from imported packages
#   --parseInternal      Include types from internal packages
#
# Usage: make swagger-gen
# Note: Run this after modifying API annotations in handlers
# -----------------------------------------------------------------------------
swagger-gen:
	swag init -g cmd/api/main.go --parseDependency --parseInternal

# =============================================================================
# TEST CONFIGURATION
# =============================================================================

# Files/patterns to EXCLUDE from coverage calculations.
# These are typically auto-generated or test-related files.
# Pattern format: Extended regex for grep -E
#   - mocks:    Generated mock files
#   - main.go:  Entry point (minimal logic)
#   - docs.go:  Generated Swagger bindings
#   - test:     Test files themselves
COVERAGE_EXCLUDE=mocks|main.go|docs.go|test

# Minimum acceptable code coverage percentage.
# Build will FAIL if coverage drops below this threshold.
# Industry standard: 70-80% for business logic
COVERAGE_THRESHOLD=80

# -----------------------------------------------------------------------------
# test: Run tests locally with coverage analysis
# -----------------------------------------------------------------------------
# Executes all tests and enforces the coverage threshold.
#
# Steps:
#   1. Run all tests with coverage profiling
#   2. Filter out excluded files (mocks, generated code)
#   3. Generate HTML coverage report
#   4. Check if coverage meets threshold (fail if below)
#
# Outputs:
#   - coverage.out:  Machine-readable coverage data
#   - coverage.html: Visual HTML report (open in browser)
#
# Usage: make test
# Exit codes:
#   0 - All tests passed AND coverage >= threshold
#   1 - Tests failed OR coverage < threshold
# -----------------------------------------------------------------------------
test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

# -----------------------------------------------------------------------------
# gen-all: Run all Go generators
# -----------------------------------------------------------------------------
# Executes go generate for the entire project.
# This typically includes:
#   - Mock generation (mockgen, mockery)
#   - Code generators (stringer, enumer)
#
# Usage: make gen-all
# Note: Run this after adding new interfaces that need mocks
# -----------------------------------------------------------------------------
gen-all:
	go generate ./...

# -----------------------------------------------------------------------------
# clean: Remove build artifacts
# -----------------------------------------------------------------------------
# Cleans up generated files and build cache.
# Use this to start fresh or troubleshoot build issues.
#
# Usage: make clean
# -----------------------------------------------------------------------------
clean:
	rm -f myapp
	go clean

# =============================================================================
# DOCKER TARGETS
# =============================================================================

# Directory where coverage artifacts from Docker builds are extracted
COVERAGE_FOLDER=./coverage

# -----------------------------------------------------------------------------
# docker-test: Run tests inside Docker container
# -----------------------------------------------------------------------------
# Executes the test suite inside a Docker container, ensuring:
#   - Tests run in the same environment as production
#   - No local dependency differences
#   - Reproducible CI/CD builds
#
# Process:
#   1. Create local coverage output directory
#   2. Build the "test" stage from Dockerfile (runs tests)
#   3. Extract coverage artifacts to local filesystem
#   4. Validate coverage meets threshold
#
# Build args:
#   COVERAGE_EXCLUDE - Patterns to exclude from coverage
#
# Outputs:
#   ./coverage/coverage.out  - Machine-readable coverage
#   ./coverage/coverage.html - Visual HTML report
#
# Usage: make docker-test
# Usage with custom exclusions:
#   make docker-test COVERAGE_EXCLUDE="mocks|generated"
# -----------------------------------------------------------------------------
docker-test:
	mkdir -p ${COVERAGE_FOLDER}
	docker buildx build --build-arg COVERAGE_EXCLUDE="${COVERAGE_EXCLUDE}" --target test -t bookmark_service:dev --output ${COVERAGE_FOLDER} .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

# -----------------------------------------------------------------------------
# docker-build: Build production Docker image
# -----------------------------------------------------------------------------
# Builds the production-ready Docker image with automatic tagging.
#
# Image naming: $(IMG_NAME):$(IMG_TAG)
# Examples:
#   - On master branch:  johnnyho3820/ebvn-bookmark-repo:dev
#   - On tag v1.0.0:   johnnyho3820/ebvn-bookmark-repo:v1.0.0
#
# Usage: make docker-build
# Usage with custom tag:
#   make docker-build IMG_TAG=custom-tag
# -----------------------------------------------------------------------------
docker-build:
	docker build -t $(IMG_NAME):$(IMG_TAG) .

# -----------------------------------------------------------------------------
# docker-release: Push image to Docker Hub
# -----------------------------------------------------------------------------
# Pushes the built image to Docker Hub registry.
#
# Prerequisites:
#   1. Must run `make docker-login` first (or be already authenticated)
#   2. Must run `make docker-build` first
#
# Usage: make docker-release
# Note: Ensure you have push access to the repository
# -----------------------------------------------------------------------------
docker-release:
	docker push $(IMG_NAME):$(IMG_TAG)

# =============================================================================
# DOCKER HUB AUTHENTICATION
# =============================================================================
# Credentials for Docker Hub login.
# These use the ?= operator which means:
#   - If the variable is already set (e.g., from environment), keep it
#   - If not set, use empty string (will need to be provided)
#
# SECURITY NOTE:
#   Never commit actual credentials! Always use:
#   - CI/CD secrets
#   - Environment variables
#   - Docker credentials store
# -----------------------------------------------------------------------------
DOCKER_HUB_USERNAME ?=
DOCKER_HUB_PASSWORD ?=

# -----------------------------------------------------------------------------
# docker-login: Authenticate with Docker Hub
# -----------------------------------------------------------------------------
# Logs into Docker Hub using provided credentials.
# Uses --password-stdin for security (avoids password in process list).
#
# Usage (environment variables - recommended for CI):
#   DOCKER_HUB_USERNAME=myuser DOCKER_HUB_PASSWORD=mypass make docker-login
#
# Usage (command line - NOT recommended):
#   make docker-login DOCKER_HUB_USERNAME=myuser DOCKER_HUB_PASSWORD=mypass
#
# Note: For local development, prefer `docker login` interactive mode
# -----------------------------------------------------------------------------
docker-login:
	echo "$(DOCKER_HUB_PASSWORD)" | docker login -u "$(DOCKER_HUB_USERNAME)" --password-stdin
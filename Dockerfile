# =============================================================================
# MULTI-STAGE DOCKERFILE FOR GO APPLICATION
# =============================================================================
# This Dockerfile uses a multi-stage build pattern to create optimized images
# for different purposes: building, testing, and production runtime.
#
# STAGES OVERVIEW:
# ┌─────────────────────────────────────────────────────────────────────────────┐
# │  base       → Common foundation with Go, dependencies, and source code     │
# │     ├── build      → Compiles the production binary                        │
# │     └── test-exec  → Runs tests and generates coverage reports             │
# │            └── test → Minimal image containing only coverage artifacts     │
# │  final      → Production runtime image (~20MB)                             │
# └─────────────────────────────────────────────────────────────────────────────┘
#
# LAYER CACHING STRATEGY:
# The order of COPY and RUN commands is optimized for Docker layer caching:
# 1. Install build-base (rarely changes → long-lived cache)
# 2. Copy go.mod/go.sum only (changes when dependencies change)
# 3. Run go mod download (cached unless dependencies change)
# 4. Copy source code last (changes frequently)
#
# This means: changing source code does NOT invalidate the dependency cache,
# significantly speeding up rebuilds (from ~5min to ~30s).
# =============================================================================

# =============================================================================
# STAGE: BASE
# =============================================================================
# Purpose: Create a common foundation with all dependencies downloaded.
#          Both 'build' and 'test-exec' stages inherit from this base,
#          avoiding duplicate work (DRY principle).
#
# Base image: golang:1.25.5-alpine
# - Alpine Linux is chosen for its minimal size (~5MB base)
# - Go 1.25.5 provides the compiler and standard library
# =============================================================================
FROM golang:1.25.5-alpine AS base

# Create the application directory structure.
# Using /opt/app follows Linux filesystem hierarchy convention
# for optional application software.
RUN mkdir -p /opt/app

# Set working directory for all subsequent commands.
# This affects COPY destinations and RUN command execution paths.
WORKDIR /opt/app

# -----------------------------------------------------------------------------
# LAYER CACHING: STEP 1 - Install build tools (rarely changes)
# -----------------------------------------------------------------------------
# Install build-base package which includes essential build tools:
# - gcc: GNU C Compiler (required for CGO - Go's C interop)
# - make: Build automation tool
# - libc-dev: C library development headers
# - g++: GNU C++ Compiler
#
# This is needed because some Go packages (e.g., SQLite drivers) contain
# C code that must be compiled. The -tags musl flag in the build step
# ensures compatibility with Alpine's musl libc implementation.
RUN apk add build-base

# -----------------------------------------------------------------------------
# LAYER CACHING: STEP 2 - Copy dependency manifests (changes occasionally)
# -----------------------------------------------------------------------------
# Copy ONLY the dependency files first, before the source code.
# This is a critical optimization:
# - go.mod: Defines the module path and dependency requirements
# - go.sum: Contains cryptographic checksums for dependency verification
#
# By copying these separately, Docker can cache the dependency download
# step and reuse it even when source code changes.
COPY go.mod ./go.mod
COPY go.sum ./go.sum

# -----------------------------------------------------------------------------
# LAYER CACHING: STEP 3 - Download dependencies (expensive, but cached!)
# -----------------------------------------------------------------------------
# Download all Go module dependencies defined in go.mod.
# This step is expensive (can take 1-2 minutes) but is cached by Docker.
# The cache is only invalidated when go.mod or go.sum changes.
#
# This means adding/updating dependencies triggers a re-download,
# but plain code changes do NOT.
RUN go mod download

# -----------------------------------------------------------------------------
# LAYER CACHING: STEP 4 - Copy source code (changes frequently)
# -----------------------------------------------------------------------------
# Copy all remaining source files from the host to the container.
# This includes all .go files, configuration files, docs/, etc.
#
# This layer is invalidated on every code change, but crucially,
# it does NOT invalidate the cached dependency download above.
COPY . .

# =============================================================================
# STAGE: BUILD
# =============================================================================
# Purpose: Compile the Go application into a production-ready binary.
#
# Inherits: FROM base (includes source code and downloaded dependencies)
# Produces: /opt/app/bookmark_service (compiled binary)
#
# This stage is used by the 'final' stage to extract the compiled binary.
# =============================================================================
FROM base AS build

# Compile the Go application with production optimizations.
#
# Build flags explained:
# - GOOS=linux: Cross-compile for Linux (ensures binary runs on Alpine)
# - -tags musl: Use musl-compatible build tags (Alpine uses musl, not glibc)
# - -ldflags "-w -s": Linker flags for size optimization:
#     -w: Omit DWARF debug information (reduces binary size)
#     -s: Omit symbol table (further reduces binary size)
#
# Input: cmd/api/main.go (application entry point)
# Output: bookmark_service (compiled binary, ~10-20MB smaller than default)
RUN GOOS=linux go build -tags musl -ldflags "-w -s" -o bookmark_service cmd/api/main.go

# =============================================================================
# STAGE: TEST-EXEC
# =============================================================================
# Purpose: Run all unit tests and generate code coverage reports.
#
# Inherits: FROM base (includes source code and downloaded dependencies)
# Produces: Coverage reports at ${_outputdir}/coverage.{out,html}
#
# This stage is typically run in CI pipelines to validate code quality
# and measure test coverage before deploying.
# =============================================================================
FROM base AS test-exec

# Build argument: Directory to store coverage output files.
# Can be overridden at build time: --build-arg _outputdir=/custom/path
ARG _outputdir="/tmp/coverage"

# Build argument: Regex pattern for files/packages to EXCLUDE from coverage.
# Example: --build-arg COVERAGE_EXCLUDE="mock|generated"
# This filters out auto-generated code from coverage calculations.
ARG COVERAGE_EXCLUDE

# Run tests and generate coverage reports in a single layer.
# This combines multiple commands to reduce layer count.
#
# Command breakdown:
# 1. mkdir -p ${_outputdir}
#    Create the output directory (with parents if needed)
#
# 2. go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
#    Run all tests with coverage:
#    - ./...: Test all packages recursively
#    - -coverprofile: Output raw coverage data to coverage.tmp
#    - -covermode=atomic: Thread-safe coverage counting (required for -race)
#    - -coverpkg=./...: Include all packages in coverage calculation
#    - -p 1: Run tests serially (prevents race conditions in integration tests)
#
# 3. grep -v -E "${COVERAGE_EXCLUDE}" coverage.tmp > ${_outputdir}/coverage.out
#    Filter out excluded patterns (e.g., mocks, generated code) from coverage
#
# 4. go tool cover -html=... -o ...
#    Generate an interactive HTML coverage report for human review
RUN mkdir -p ${_outputdir} && \
    go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
    grep -v -E "${COVERAGE_EXCLUDE}}" coverage.tmp > ${_outputdir}/coverage.out && \
    go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html

# =============================================================================
# STAGE: TEST
# =============================================================================
# Purpose: Create a minimal image containing ONLY the coverage artifacts.
#
# Base image: scratch (completely empty, 0 bytes)
#
# This stage is used to extract coverage reports from Docker builds:
#   docker build --target test --output type=local,dest=./coverage .
#
# The scratch base image ensures the output contains nothing except
# the coverage files, making it easy to extract build artifacts.
# =============================================================================
FROM scratch AS test

# Re-declare ARG to make it available in this stage.
# ARG values do not persist across FROM statements.
ARG _outputdir="/tmp/coverage"

# Copy only the coverage artifacts from the test-exec stage.
# These files are placed at the root (/) of the scratch image.
# - coverage.out: Machine-readable coverage data (for CI tools)
# - coverage.html: Human-readable HTML report
COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /

# =============================================================================
# STAGE: FINAL (Production Runtime)
# =============================================================================
# Purpose: Create a minimal production image with only the compiled binary.
#
# Base image: alpine (minimal Linux, ~5MB)
# Final size: ~20MB (vs ~1GB if we included Go compiler)
#
# This stage does NOT include:
# - Go compiler/runtime source
# - Build tools (gcc, make)
# - Source code
# - Test files
#
# This dramatically reduces:
# - Image size (faster pulls, less storage)
# - Attack surface (fewer potential vulnerabilities)
# - Memory footprint
# =============================================================================
FROM alpine AS final

# -----------------------------------------------------------------------------
# SECURITY: Run as Non-Root User
# -----------------------------------------------------------------------------
# Create a dedicated user and group for running the application.
# This is a security best practice that follows the principle of least privilege.
#
# Why run as non-root?
#   - Limits damage if the application is compromised
#   - Prevents accidental modification of system files
#   - Required by many container orchestrators (OpenShift, some K8s policies)
#   - Reduces attack surface significantly
#
# Flags explained:
#   -S: Create a system user/group (no home dir, no password)
#   -G appgroup: Add user to the specified group
#
# After this, all subsequent commands (COPY, RUN, CMD) will run as 'appuser'
# -----------------------------------------------------------------------------
RUN addgroup -S appgroup \
    && adduser -S appuser -G appgroup

USER appuser

# Build argument for application name (unused in current config, but available
# for future customization like naming log files or process names).
ARG app_name=app

# Set timezone environment variable.
# This affects time formatting in application logs and any time-sensitive logic.
# Asia/Ho_Chi_Minh = UTC+7 (Vietnam timezone)
ENV TZ=Asia/Ho_Chi_Minh

# Set working directory for the production container.
# The application will run from this directory.
WORKDIR /app

# Copy the compiled binary from the build stage.
# --from=build references the stage named "build" above.
# This copies ONLY the binary, not the source code or build tools.
COPY --from=build /opt/app/bookmark_service /app/bookmark_service

# Copy the Swagger/OpenAPI documentation files.
# These are required for the /swagger endpoint to serve API documentation.
COPY --from=build /opt/app/docs /app/docs

# Configure the container's timezone.
# This creates a symlink from /etc/localtime to the appropriate timezone file
# and writes the timezone name to /etc/timezone.
# This ensures consistent timestamp formatting across the application.
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Define the default command to run when the container starts.
# This executes the compiled Go binary.
# Using exec form (JSON array) instead of shell form for proper signal handling.
CMD ["/app/bookmark_service"]

# =============================================================================
# BUILD COMMANDS
# =============================================================================
# Production image:
#   docker build --target final -t bookmark_service:latest .
#
# Run tests with coverage:
#   docker build --target test-exec -t bookmark_service:test .
#
# Extract coverage reports to local filesystem:
#   docker build --target test --output type=local,dest=./coverage .
#
# Build all stages (useful for CI):
#   docker build -t bookmark_service:dev .
# =============================================================================
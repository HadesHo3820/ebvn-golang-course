# =============================================================================
# MULTI-STAGE DOCKERFILE FOR GO APPLICATION
# =============================================================================
# This Dockerfile uses a multi-stage build to create a minimal production image.
# Stage 1 (build): Compiles the Go application with all build dependencies.
# Stage 2 (run): Creates a lightweight Alpine image with only the compiled binary.
# This approach reduces the final image size from ~1GB to ~20MB.
# =============================================================================

# -----------------------------------------------------------------------------
# STAGE 1: BUILD
# -----------------------------------------------------------------------------
# Use Go 1.25 on Alpine Linux as the build environment.
# Alpine is chosen for its small size (~5MB base).
FROM golang:1.25.5-alpine AS build

# Create the application directory inside the container.
RUN mkdir -p /opt/app

# Set the working directory for all subsequent commands.
WORKDIR /opt/app

# Copy all source files from the host to the container.
# This includes go.mod, go.sum, and all .go files.
COPY . .

# Install build-base package which includes:
# - gcc: GNU C Compiler (required for CGO)
# - make: Build automation tool
# - libc-dev: C library development files
# Some Go packages with C bindings (like SQLite) require these tools.
RUN apk add build-base

# Download Go module dependencies and compile the application.
# - go mod download: Fetches all dependencies defined in go.mod
# - go build -o bookmark_service: Compiles and outputs the binary
# The resulting binary is created at: /opt/app/bookmark_service
RUN go mod download && \
    go build -o bookmark_service cmd/api/main.go

# -----------------------------------------------------------------------------
# STAGE 2: RUN (Production)
# -----------------------------------------------------------------------------
# Use a minimal Alpine image for the production container.
# This image does NOT include Go or build tools, only the compiled binary.
FROM alpine AS run

# Set the working directory for the runtime container.
WORKDIR /app

# Copy the compiled binary from the build stage.
# --from=build references the first stage named "build".
COPY --from=build /opt/app/bookmark_service .

# Copy the Swagger documentation files from the build stage.
# These are needed for the /swagger endpoint to work.
COPY --from=build /opt/app/docs ./docs

# Define the default command to run when the container starts.
# This executes the compiled Go binary.
CMD ["/app/bookmark_service"]

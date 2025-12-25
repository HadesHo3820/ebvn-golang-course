# Golang Backend Development Course - Bookmark & Shorten Link System

A comprehensive Golang backend development course project focusing on professional backend architecture, clean code practices, and production-ready API development.

## 📚 Course Introduction

This course provides in-depth training on:

- **Writing Clean Go Code**: Following industry best practices and standards
- **Professional Backend Architecture Design**: Building robust and scalable systems
- **Understanding Patterns and Principles**: Deep dive into clean code and clean architecture
- **Building Production-Ready APIs**: From error handling to CI/CD

## 🎯 Course Objectives

The goal is to help students learn to think and develop backend applications like a **software engineer**, not just a "coder". By the end of this course, you'll understand how to:

- Design scalable and maintainable backend systems
- Implement clean architecture patterns
- Build production-grade APIs with proper error handling and testing
- Deploy and manage microservices in cloud environments

---

## 🏗️ Project Structure

```
ebvn-golang-course/
├── cmd/
│   └── api/                    # Application entrypoint
│       └── main.go
├── internal/                   # Private application code
│   ├── api/                    # API server setup & configuration
│   │   ├── api.go              # HTTP server & routing
│   │   └── config.go           # Environment configuration
│   ├── handler/                # HTTP handlers (adapters)
│   │   ├── pass_handler.go
│   │   ├── pass_handler_test.go
│   │   ├── health_check_handler.go
│   │   └── health_check_handler_test.go
│   ├── service/                # Business logic (core)
│   │   ├── pass_service.go
│   │   ├── pass_service_test.go
│   │   ├── health_check_service.go
│   │   ├── health_check_service_test.go
│   │   └── mocks/              # Generated mocks for testing
│   ├── model/                  # Domain models
│   ├── repository/             # Data access layer
│   └── test/
│       └── endpoint/           # Integration/E2E tests
├── docs/                       # Swagger documentation
├── Makefile                    # Build & run commands
├── go.mod
└── go.sum
```

## 🛠️ Technology Stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.25+ |
| **Web Framework** | [Gin](https://github.com/gin-gonic/gin) |
| **Configuration** | [envconfig](https://github.com/kelseyhightower/envconfig) |
| **API Documentation** | [Swagger/OpenAPI](https://github.com/swaggo/swag) |
| **Testing** | [testify](https://github.com/stretchr/testify), [mockery](https://github.com/vektra/mockery) |
| **Architecture** | Hexagonal Architecture (Ports & Adapters) |

## 🚀 Getting Started

### Prerequisites

- Go 1.25 or later
- Make (optional, for using Makefile commands)

### Installation

```bash
# Clone the repository
git clone https://github.com/HadesHo3820/ebvn-golang-course.git
cd ebvn-golang-course

# Install dependencies
go mod download
```

### Running the Application

```bash
# Using Makefile
make start

# Or directly with Go
go run cmd/api/main.go
```

The server starts on port `8080` by default.

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_PORT` | `8080` | Server port |
| `SERVICE_NAME` | `bookmark-api` | Service name for health check |
| `INSTANCE_ID` | Auto-generated UUID | Unique instance identifier |

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health-check` | Health check endpoint |
| `GET` | `/gen-pass` | Generate a random password |
| `GET` | `/swagger/*` | Swagger UI documentation |

### Health Check Response

```json
{
  "message": "OK",
  "service_name": "bookmark-api",
  "instance_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## 🧪 Testing

The project includes comprehensive tests at multiple levels:

```bash
# Run all tests
make test

# Or directly with Go
go test ./... -v
```

### Test Coverage

- **Unit Tests**: Service and handler layer tests with mocks
- **Integration Tests**: End-to-end endpoint tests

### Mock Generation

Mocks are generated using `mockery`:

```bash
go generate ./...
```

## 📖 API Documentation

Swagger documentation is available at:

```
http://localhost:8080/swagger/index.html
```

To regenerate Swagger docs after API changes:

```bash
make swagger-gen
```

## 🏛️ Architecture

This project follows **Hexagonal Architecture** (Ports & Adapters):

```
┌─────────────────────────────────────────────────┐
│                   Handlers                      │  ← HTTP Adapters
│         (pass_handler, health_handler)          │
├─────────────────────────────────────────────────┤
│                   Services                      │  ← Core Business Logic
│         (pass_service, health_service)          │
├─────────────────────────────────────────────────┤
│                 Repositories                    │  ← Data Adapters
│              (Future: Database)                 │
└─────────────────────────────────────────────────┘
```

**Key Principles:**
- **Dependency Injection**: Services are injected into handlers
- **Interface-based Design**: All layers communicate via interfaces
- **Testability**: Mocks enable isolated unit testing

---

## 📚 Makefile Commands

| Command | Description |
|---------|-------------|
| `make start` | Run the application |
| `make test` | Run all tests |
| `make swagger-gen` | Regenerate Swagger documentation |

---

## 🎓 Target Audience

This course is designed for developers who want to:
- Level up from writing basic code to building professional backend systems
- Understand the thinking process of professional software engineers
- Build production-ready applications with industry best practices
- Gain hands-on experience with modern backend technologies

---

**Note**: This project serves as a practical, hands-on learning experience for backend development with Go, emphasizing real-world application and professional engineering practices.

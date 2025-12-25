# Swagger API Documentation Tutorial

This guide explains how to integrate [swaggo/swag](https://github.com/swaggo/swag) into a Go Gin project to automatically generate OpenAPI (Swagger) documentation.

## Overview

```
┌─────────────────┐     ┌──────────────┐     ┌─────────────────┐
│  Go Comments    │ ──▶ │  swag init   │ ──▶ │  docs/ folder   │
│  (annotations)  │     │  (CLI tool)  │     │  swagger.json   │
└─────────────────┘     └──────────────┘     └─────────────────┘
                                                     │
                                                     ▼
                                            ┌─────────────────┐
                                            │  Swagger UI     │
                                            │  /swagger/      │
                                            └─────────────────┘
```

---

## Step 1: Install the Swag CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Add to your PATH (in `~/.zshrc`):
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## Step 2: Install Dependencies

```bash
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

---

## Step 3: Add General API Info (main.go)

Add these annotations **before** the `main()` function:

```go
package main

import (
    _ "your-module/docs"  // Import generated docs
    "your-module/internal/api"
)

// @title Your API Title
// @version 1.0
// @description Your API description here.

// @host localhost:8080
// @BasePath /
func main() {
    // ...
}
```

| Annotation | Description |
|------------|-------------|
| `@title` | API title shown in Swagger UI |
| `@version` | API version |
| `@description` | Brief description of your API |
| `@host` | Host URL (without protocol) |
| `@BasePath` | Base path for all endpoints |

---

## Step 4: Add Endpoint Annotations (handlers)

Add annotations **before** each handler function:

```go
// @Summary Short description of endpoint
// @Description Detailed description of what this endpoint does
// @Tags tag-name
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param body body RequestStruct true "Request body"
// @Success 200 {object} ResponseStruct "Success response"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal error"
// @Router /your-endpoint [get]
func (h *handler) YourHandler(c *gin.Context) {
    // ...
}
```

### Common Annotations

| Annotation | Example | Description |
|------------|---------|-------------|
| `@Summary` | `Generate password` | Short one-line description |
| `@Description` | `Creates a secure password` | Detailed description |
| `@Tags` | `password` | Groups endpoints in UI |
| `@Accept` | `json` | Request content type |
| `@Produce` | `json` or `plain` | Response content type |
| `@Param` | `id path int true "ID"` | Parameter definition |
| `@Success` | `200 {object} User` | Success response |
| `@Failure` | `400 {string} string` | Error response |
| `@Router` | `/users/{id} [get]` | Route path and method |

### @Param Format

```
@Param name location type required "description"
```

- **location**: `path`, `query`, `header`, `body`, `formData`
- **type**: `int`, `string`, `bool`, or struct name
- **required**: `true` or `false`

---

## Step 5: Register Swagger Route (api.go)

```go
import (
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func (a *api) RegisterEP() {
    // ... your other routes
    
    // Swagger documentation endpoint
    a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

---

## Step 6: Generate Documentation

Run from your project root:

```bash
swag init -g cmd/api/main.go --parseDependency --parseInternal
```

This creates:
```
docs/
├── docs.go       # Go code to embed docs
├── swagger.json  # OpenAPI spec (JSON)
└── swagger.yaml  # OpenAPI spec (YAML)
```

### Flags

| Flag | Description |
|------|-------------|
| `-g` | Path to main.go with general API info |
| `--parseDependency` | Parse dependencies for struct definitions |
| `--parseInternal` | Parse internal packages |

---

## Step 7: Access Swagger UI

1. Run your server:
   ```bash
   go run cmd/api/main.go
   ```

2. Open in browser:
   ```
   http://localhost:8080/swagger/index.html
   ```

---

## Regenerating Docs

**Important**: Run `swag init` every time you change annotations!

Add to your Makefile:
```makefile
.PHONY: docs
docs:
	swag init -g cmd/api/main.go --parseDependency --parseInternal
```

Or add a `//go:generate` directive:
```go
//go:generate swag init -g cmd/api/main.go --parseDependency --parseInternal
```

---

## Project Structure Example

```
your-project/
├── cmd/api/
│   └── main.go           # General API info annotations
├── internal/
│   ├── api/
│   │   └── api.go        # Swagger route registration
│   └── handler/
│       └── handler.go    # Endpoint annotations
├── docs/                  # Generated (git-ignore or commit)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
└── go.mod
```

---

## Tips

1. **Don't forget the import**: `_ "your-module/docs"` in main.go
2. **Regenerate after changes**: Always run `swag init` after modifying annotations
3. **Struct documentation**: Add `// @Description` above struct fields for better docs
4. **Example values**: Use `example:"value"` struct tags for sample data

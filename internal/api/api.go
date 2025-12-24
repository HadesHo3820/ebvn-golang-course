// Package api provides the HTTP server setup and routing configuration.
// It acts as the entry point for the application, wiring together
// all handlers and services following the Hexagonal Architecture pattern.
package api

import (
	"fmt"
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/handler"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Engine defines the interface for the API server.
// It abstracts the server implementation, allowing for easier testing
// and potential swapping of the underlying HTTP framework.
type Engine interface {
	// Start runs the HTTP server on the default port (8080).
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// api is the concrete implementation of the Engine interface.
// It wraps a Gin engine and manages the application's HTTP routing.
type api struct {
	app *gin.Engine
	cfg *Config
}

// New creates and initializes a new API server.
// It sets up the Gin engine and registers all endpoints.
// Returns an Engine interface to hide the implementation details.
func New(cfg *Config) Engine {
	a := &api{
		app: gin.New(),
		cfg: cfg,
	}
	a.RegisterEP()
	return a
}

// Start begins listening for HTTP requests.
// By default, Gin listens on port 8080.
// Returns an error if the server fails to start.
func (a *api) Start() error {
	return a.app.Run(fmt.Sprintf(":%s", a.cfg.AppPort))
}

// ServeHTTP serves an HTTP request with the given response writer and request.
// It uses the underlying Gin engine to handle the request and write the response.
func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

// RegisterEP sets up all API endpoints and their handlers.
// It performs dependency injection by:
//  1. Creating service instances (business logic layer)
//  2. Injecting services into handlers (HTTP adapter layer)
//  3. Registering handlers with their respective routes
//
// Endpoints:
//   - GET /gen-pass: Generates a random password
//   - GET /swagger/*any: Swagger UI documentation
func (a *api) RegisterEP() {
	// Initialize the password service (core business logic)
	passSvc := service.NewPassword()

	// Create the password handler with injected service dependency
	passHandler := handler.NewPassword(passSvc)

	// Register the password generation endpoint
	a.app.GET("/gen-pass", passHandler.GenPass)

	// Register Swagger documentation endpoint
	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

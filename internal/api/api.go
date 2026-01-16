// Package api provides the HTTP server setup and routing configuration.
// It acts as the entry point for the application, wiring together
// all handlers and services following the Hexagonal Architecture pattern.
package api

import (
	"fmt"
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/docs"
	"github.com/HadesHo3820/ebvn-golang-course/internal/handler"
	"github.com/HadesHo3820/ebvn-golang-course/internal/repository"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
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
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
	keyGen      stringutils.KeyGenerator
	db          *gorm.DB
}

// New creates and initializes a new API server.
// It sets up the Gin engine and registers all endpoints.
// Returns an Engine interface to hide the implementation details.
func New(app *gin.Engine, cfg *Config, redisClient *redis.Client, keyGen stringutils.KeyGenerator, db *gorm.DB) Engine {
	a := &api{
		app:         app,
		cfg:         cfg,
		redisClient: redisClient,
		keyGen:      keyGen,
		db:          db,
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
//   - GET /health-check: Health check endpoint
//   - POST /links/shorten: Shorten a URL
//   - GET /swagger/*any: Swagger UI documentation
func (a *api) RegisterEP() {
	// Initialize the password service (core business logic)
	passSvc := service.NewPassword()

	// Initialize Redis health checker for dependency verification
	// Use the injected Redis client from the constructor
	redisHealthChecker := repository.NewRedisHealthChecker(a.redisClient)

	healthSvc := service.NewHealthCheck(a.cfg.ServiceName, a.cfg.InstanceID, redisHealthChecker)

	// Initialize URL storage repository and service
	urlRepo := repository.NewUrlStorage(a.redisClient)
	urlSvc := service.NewShortenUrl(urlRepo, a.keyGen)

	// Create the password handler with injected service dependency
	passHandler := handler.NewPassword(passSvc)

	// Create the health handler with injected service dependency
	healthHandler := handler.NewHealthCheck(healthSvc)

	// Create the URL shorten handler with injected service dependency
	urlShortenHandler := handler.NewUrlShorten(urlSvc)

	// create user handler
	userRepo := repository.NewUser(a.db)
	userSvc := service.NewUser(userRepo)
	userHandler := handler.NewUser(userSvc)

	// v1Routes creates a route group with "/v1" prefix for API versioning.
	// All routes registered under this group will be prefixed with "/v1",
	// allowing for future API versions (e.g., "/v2") without breaking existing clients.
	// The curly braces are purely for visual grouping and have no effect on scope.
	v1Routes := a.app.Group("/v1")
	{
		// GET /v1/gen-pass - Generates a random password
		v1Routes.GET("/gen-pass", passHandler.GenPass)

		// GET /v1/health-check - Returns service health status including Redis connectivity
		v1Routes.GET("/health-check", healthHandler.Check)

		// POST /v1/links/shorten - Creates a shortened URL code for the provided URL
		v1Routes.POST("/links/shorten", urlShortenHandler.ShortenUrl)

		// GET /v1/links/redirect/{code} - Redirects to the original URL for the provided short code
		v1Routes.GET("/links/redirect/:code", urlShortenHandler.GetUrl)

		// POST /v1/users/register - Registers a new user
		v1Routes.POST("/users/register", userHandler.RegisterUser)
	}

	// Register Swagger documentation endpoint
	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Configure Swagger host dynamically at runtime.
	// This overrides the @host annotation defined in the main.go swagger comments.
	// Why this is needed:
	//   - When running behind a reverse proxy (e.g., NGINX), the API is accessed
	//     via a different hostname/path (e.g., "localhost/api/bookmark_service")
	//   - The default @host value may not match the actual deployment URL
	//   - Setting this dynamically allows the Swagger UI "Try it out" feature
	//     to send requests to the correct endpoint based on environment config
	// Example: APP_HOSTNAME=localhost/api/bookmark_service makes Swagger send
	//          requests to http://localhost/api/bookmark_service/v1/health-check
	docs.SwaggerInfo.Host = a.cfg.AppHostName
}

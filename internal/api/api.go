// Package api provides the HTTP server setup and routing configuration.
// It acts as the entry point for the application, wiring together
// all handlers and services following the Hexagonal Architecture pattern.
package api

import (
	"fmt"
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/docs"
	"github.com/HadesHo3820/ebvn-golang-course/internal/api/middleware"
	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/healthcheck"
	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/password"
	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/url"
	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/user"
	"github.com/HadesHo3820/ebvn-golang-course/internal/repository"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/jwtutils"
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
	app          *gin.Engine
	cfg          *Config
	redisClient  *redis.Client
	db           *gorm.DB
	keyGen       stringutils.KeyGenerator
	jwtGen       jwtutils.JWTGenerator
	jwtValidator jwtutils.JWTValidator
}

type EngineOpts struct {
	Engine       *gin.Engine
	Cfg          *Config
	RedisClient  *redis.Client
	SqlDB        *gorm.DB
	KeyGen       stringutils.KeyGenerator
	JwtGen       jwtutils.JWTGenerator
	JwtValidator jwtutils.JWTValidator
}

// New creates and initializes a new API server.
// It sets up the Gin engine and registers all endpoints.
// Returns an Engine interface to hide the implementation details.
func New(opts *EngineOpts) Engine {
	a := &api{
		app:          opts.Engine,
		cfg:          opts.Cfg,
		redisClient:  opts.RedisClient,
		keyGen:       opts.KeyGen,
		db:           opts.SqlDB,
		jwtGen:       opts.JwtGen,
		jwtValidator: opts.JwtValidator,
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

// handlers aggregates all HTTP handlers required by the API.
// It serves as a container for dependency-injected handler instances,
// making it easier to manage and pass handlers to route registration.
type handlers struct {
	healthCheckHandler healthcheck.HealthCheckHandler // Handles health check endpoints
	passwordHandler    password.PasswordHandler       // Handles password generation endpoints
	urlShortenHandler  url.UrlHandler                 // Handles URL shortening endpoints
	userHandler        user.UserHandler               // Handles user management endpoints
}

// initHandlers initializes all handlers with their required dependencies.
// It follows the Hexagonal Architecture pattern by:
//  1. Creating repository instances (infrastructure layer)
//  2. Creating service instances with injected repositories (domain layer)
//  3. Creating handler instances with injected services (adapter layer)
//
// This method centralizes dependency injection, making it easier to:
//   - Understand the dependency graph of the application
//   - Test handlers with mock dependencies
//   - Add new handlers consistently
//
// Returns a handlers struct containing all initialized handler instances.
func (a *api) initHandlers() *handlers {
	// Create health check service with Redis connectivity checker
	healthCheckRepo := repository.NewRedisHealthChecker(a.redisClient)
	healthSvc := service.NewHealthCheck(a.cfg.ServiceName, a.cfg.InstanceID, healthCheckRepo)

	// Create URL shortening service with Redis storage
	urlRepo := repository.NewUrlStorage(a.redisClient)
	urlSvc := service.NewShortenUrl(urlRepo, a.keyGen)

	// Create password service (stateless, no repository needed)
	passSvc := service.NewPassword()

	// Create user service with PostgreSQL repository
	userRepo := repository.NewUser(a.db)
	userSvc := service.NewUser(userRepo, a.jwtGen)

	return &handlers{
		healthCheckHandler: healthcheck.NewHealthCheckHandler(healthSvc),
		passwordHandler:    password.NewPasswordHandler(passSvc),
		urlShortenHandler:  url.NewUrlHandler(urlSvc),
		userHandler:        user.NewUserHandler(userSvc),
	}
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
	allHandlers := a.initHandlers()

	// GET /health-check - Returns service health status including Redis connectivity
	a.app.GET("/health-check", allHandlers.healthCheckHandler.Ping)

	// v1PublicRoutes creates a route group with "/v1" prefix for API versioning.
	// All routes registered under this group will be prefixed with "/v1",
	// allowing for future API versions (e.g., "/v2") without breaking existing clients.
	// The curly braces are purely for visual grouping and have no effect on scope.
	v1PublicRoutes := a.app.Group("/v1")
	{
		// GET /v1/gen-pass - Generates a random password
		v1PublicRoutes.GET("/gen-pass", allHandlers.passwordHandler.GenPass)

		// POST /v1/links/shorten - Creates a shortened URL code for the provided URL
		v1PublicRoutes.POST("/links/shorten", allHandlers.urlShortenHandler.ShortenUrl)

		// GET /v1/links/redirect/{code} - Redirects to the original URL for the provided short code
		v1PublicRoutes.GET("/links/redirect/:code", allHandlers.urlShortenHandler.GetUrl)

		// POST /v1/users/register - Registers a new user
		v1PublicRoutes.POST("/users/register", allHandlers.userHandler.Register)

		// POST /v1/users/login - Logs in a user and returns a JWT token
		v1PublicRoutes.POST("/users/login", allHandlers.userHandler.Login)
	}

	jwtMiddleware := middleware.NewJWTAuth(a.jwtValidator)

	v1PrivateRoutes := a.app.Group("/v1")
	v1PrivateRoutes.Use(jwtMiddleware.JWTAuth())
	{
		// GET /v1/self/info - Gets the authenticated user's profile information
		v1PrivateRoutes.GET("/self/info", allHandlers.userHandler.GetSelfInfo)

		// PUT /v1/self/info - Updates the authenticated user's profile information
		v1PrivateRoutes.PUT("/self/info", allHandlers.userHandler.UpdateSelfInfo)
	}

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

	// Register Swagger documentation endpoint
	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

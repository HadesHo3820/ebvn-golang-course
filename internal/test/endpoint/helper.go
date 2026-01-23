package endpoint

import (
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/api"
	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/jwtutils"
	jwtMocks "github.com/HadesHo3820/ebvn-golang-course/pkg/jwtutils/mocks"
	redisPkg "github.com/HadesHo3820/ebvn-golang-course/pkg/redis"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// defaultTestConfig returns the default API config for testing.
func defaultTestConfig() *api.Config {
	return &api.Config{
		ServiceName: "test-service",
		InstanceID:  "1234",
	}
}

// TestEngineOpts holds options for creating a test API engine.
type TestEngineOpts struct {
	// T is the testing instance (required)
	T *testing.T

	// Cfg is the API config (optional, uses default if nil)
	Cfg *api.Config

	// Fixture is the database fixture to use (optional, uses UserCommonTestDB if nil)
	Fixture fixture.Fixture

	// JwtGen is the JWT generator mock (optional, creates new mock if nil)
	JwtGen jwtutils.JWTGenerator

	// JwtValidator is the JWT validator mock (optional, creates new mock if nil)
	JwtValidator jwtutils.JWTValidator
}

// TestEngine wraps the API engine with test dependencies for easy access.
type TestEngine struct {
	Engine       api.Engine
	DB           *gorm.DB
	JwtGen       *jwtMocks.JWTGenerator
	JwtValidator *jwtMocks.JWTValidator
}

// NewTestEngine creates an API engine configured for testing with fixture database
// and mock dependencies. This reduces boilerplate in endpoint tests.
//
// Parameters:
//   - opts: Configuration options for the test engine
//
// Returns a TestEngine with the API engine and mock references for assertions.
func NewTestEngine(opts *TestEngineOpts) *TestEngine {
	if opts.T == nil {
		panic("TestEngineOpts.T is required")
	}

	// Use default config if not provided
	cfg := opts.Cfg
	if cfg == nil {
		cfg = defaultTestConfig()
	}

	// Use default fixture if not provided
	fix := opts.Fixture
	if fix == nil {
		fix = &fixture.UserCommonTestDB{}
	}

	// Setup test database with fixture
	db := fixture.NewFixture(opts.T, fix)

	// Create or use provided JWT mocks
	var jwtGen *jwtMocks.JWTGenerator
	var jwtValidator *jwtMocks.JWTValidator

	if opts.JwtGen != nil {
		// Type assert if it's already a mock
		if mock, ok := opts.JwtGen.(*jwtMocks.JWTGenerator); ok {
			jwtGen = mock
		}
	}
	if jwtGen == nil {
		jwtGen = jwtMocks.NewJWTGenerator(opts.T)
	}

	if opts.JwtValidator != nil {
		if mock, ok := opts.JwtValidator.(*jwtMocks.JWTValidator); ok {
			jwtValidator = mock
		}
	}
	if jwtValidator == nil {
		jwtValidator = jwtMocks.NewJWTValidator(opts.T)
	}

	// Create API engine with dependencies
	engine := api.New(&api.EngineOpts{
		Engine:          gin.New(),
		Cfg:             cfg,
		RedisClient:     redisPkg.InitMockRedis(opts.T),
		SqlDB:           db,
		KeyGen:          stringutils.NewKeyGenerator(),
		PasswordHashing: utils.NewPasswordHashing(),
		JwtGen:          jwtGen,
		JwtValidator:    jwtValidator,
	})

	return &TestEngine{
		Engine:       engine,
		DB:           db,
		JwtGen:       jwtGen,
		JwtValidator: jwtValidator,
	}
}

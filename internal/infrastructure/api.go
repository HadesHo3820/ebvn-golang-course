// Package infrastructure provides factory functions for initializing
// and wiring up application dependencies during startup.
package infrastructure

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/api"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/common"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/logger"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateAPIConfig loads configuration from environment variables and
// generates a unique instance ID if not provided.
func CreateAPIConfig() *api.Config {
	cfg, err := api.NewConfig()

	common.HandleError(err)

	if cfg.InstanceID == "" {
		// The most commonly used UUID library for Go is github.com/google/uuid
		cfg.InstanceID = uuid.New().String()
	}

	return cfg
}

// CreateAPI initializes and returns a fully configured API engine.
func CreateAPI() api.Engine {
	logger.SetLogLevel()

	// Init config
	cfg := CreateAPIConfig()
    
	// Init redis - 
	redisClient := CreateRedisGeneralConn()
	cacheRedisClient := CreateRedisCacheConn()

	// Init sql DB
	sqlDB := CreateSQLDBWithMigration()

	// Init key gen
	keyGen := stringutils.NewKeyGenerator()

	// Init jwt gen and validator
	jwtGen, jwtValidator := CreateJWTProvider()

	// Init password hashing
	passwordHashing := utils.NewPasswordHashing()

	app := gin.New()

	return api.New(&api.EngineOpts{
		Engine:           app,
		Cfg:              cfg,
		RedisClient:      redisClient,
		CacheRedisClient: cacheRedisClient,
		SqlDB:            sqlDB,
		KeyGen:           keyGen,
		PasswordHashing:  passwordHashing,
		JwtGen:           jwtGen,
		JwtValidator:     jwtValidator,
	})
}

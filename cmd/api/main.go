package main

import (
	_ "github.com/HadesHo3820/ebvn-golang-course/docs"
	"github.com/HadesHo3820/ebvn-golang-course/internal/api"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/logger"
	redisPkg "github.com/HadesHo3820/ebvn-golang-course/pkg/redis"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/sqldb"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// @title EBVN Bookmark API
// @version 1.2
// @description A simple password generator API built with Go and Gin.
// @host localhost:8080
// @BasePath /
func main() {
	logger.SetLogLevel()

	cfg := createAPIConfig()

	redisClient := createRedisClient()

	app := createAPIApp(cfg, redisClient)

	app.Start()
}

// createRedisClient creates Redis client
func createRedisClient() *redis.Client {
	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}
	return redisClient
}

// createAPIConfig creates API configuration based on environment variables
func createAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

// createAPIApp creates API application based on API configuration
func createAPIApp(cfg *api.Config, redis *redis.Client) api.Engine {
	app := gin.New()

	// create db
	db, err := sqldb.NewClient("")
	if err != nil {
		panic(err)
	}

	keyGen := stringutils.NewKeyGenerator()

	a := api.New(app, cfg, redis, keyGen, db)
	return a
}

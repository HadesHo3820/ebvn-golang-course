package main

import (
	_ "github.com/HadesHo3820/ebvn-golang-course/docs"
	"github.com/HadesHo3820/ebvn-golang-course/internal/api"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/logger"
	redisPkg "github.com/HadesHo3820/ebvn-golang-course/pkg/redis"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils"
)

// @title EBVN Bookmark API
// @version 1.0
// @description A simple password generator API built with Go and Gin.
// @host localhost:8080
// @BasePath /
func main() {
	logger.SetLogLevel()

	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}

	keyGen := stringutils.NewKeyGenerator()

	app := api.New(cfg, redisClient, keyGen)
	app.Start()
}

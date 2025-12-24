package main

import (
	_ "github.com/HadesHo3820/ebvn-golang-course/docs"
	"github.com/HadesHo3820/ebvn-golang-course/internal/api"
)

// @title EBVN Bookmark API
// @version 1.0
// @description A simple password generator API built with Go and Gin.
// @BasePath /
func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	app := api.New(cfg)
	app.Start()
}

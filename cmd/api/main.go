package main

import (
	_ "github.com/HadesHo3820/ebvn-golang-course/docs"
	"github.com/HadesHo3820/ebvn-golang-course/internal/infrastructure"
)

// @title EBVN Bookmark API
// @version 1.3
// @description A simple password generator API built with Go and Gin.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your Bearer token in the format: Bearer {token}
func main() {
	// Init api
	a := infrastructure.CreateAPI()

	// Start api app
	a.Start()
}

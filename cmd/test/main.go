package main

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/repository"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/redis"
)

func main() {
	ctx := context.Background()
	urlStorage, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}

	urlRepo := repository.NewUrlStorage(urlStorage)

	urlService := service.NewShortenUrl(urlRepo)

	key, _ := urlService.ShortenUrl(ctx, "https://google.com")

	println(key)
}

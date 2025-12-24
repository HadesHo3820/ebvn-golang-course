package main

import "github.com/HadesHo3820/ebvn-golang-course/internal/api"

func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	app := api.New(cfg)
	app.Start()
}

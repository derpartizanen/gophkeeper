package main

import (
	"log"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/app"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}

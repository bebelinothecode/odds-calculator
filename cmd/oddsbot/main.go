package main

import (
	"context"
	"log"

	"oddsbot/internal/app"
	"oddsbot/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	a, err := app.New(cfg)
	if err != nil {
		log.Fatalf("app init: %v", err)
	}

	if err := a.Run(context.Background()); err != nil {
		log.Fatalf("server: %v", err)
	}
}

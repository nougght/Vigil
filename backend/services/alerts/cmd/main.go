package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/nougght/monitoring-system/backend/services/alerts/internal/app"
	"github.com/nougght/monitoring-system/backend/services/alerts/internal/config"
)

func main() {
	cfg := config.MustLoadConfig("config.yaml")
	rootCtx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	app := app.New(rootCtx, cfg)

	if err := app.Run(rootCtx); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}

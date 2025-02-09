package main

import (
	"github.com/ArtemSarafannikov/OzonTestTask/internal/app"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/config"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	application := app.New(log, cfg)
	application.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

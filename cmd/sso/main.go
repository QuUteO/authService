package main

import (
	"Auth/internal/app"
	"Auth/internal/config"
	"Auth/internal/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// start:go run cmd/sso/main.go --config=./config/config.yaml

func main() {
	// инициализация конфига
	cfg := config.NewConfig()

	// инициализация логгера
	log := logger.SetupLogger(cfg.Env)

	application := app.NewApp(log, cfg.Grpc.Port, cfg.StoragePath, cfg.TokenTTL)

	// запускаем сервер с слушателем
	go application.GRPCServer.MustRun()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	result := <-sigChan
	application.GRPCServer.Stop()
	log.Info("shutting down", slog.String("reason", result.String()))

}

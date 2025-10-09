package app

import (
	grpcapp "Auth/internal/app/grpc"
	"Auth/internal/service/auth"
	"Auth/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

// NewApp инициализирует grpc + service слой + storage
func NewApp(log *slog.Logger, grpcPort int, storagePath string, TokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, TokenTTL)

	// grpc—сервер инициализация
	grpcApp := grpcapp.New(log, grpcPort, authService)
	return &App{
		GRPCServer: grpcApp,
	}
}

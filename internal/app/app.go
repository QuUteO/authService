package app

import (
	grpcapp "Auth/internal/app/grpc"
	authgrpc "Auth/internal/grpc/Auth"
	"Auth/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

// NewApp инициализирует grpc + service слой + storage
func NewApp(log *slog.Logger, grpcPort int, storagePath string, TokenTTL time.Duration, auth authgrpc.Auth) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New()
	// grpc—сервер инициализация
	grpcApp := grpcapp.New(log, grpcPort, auth)
	return &App{
		GRPCServer: grpcApp,
	}
}

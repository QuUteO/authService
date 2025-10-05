package app

import (
	grpcapp "Auth/internal/app/grpc"
	authgrpc "Auth/internal/grpc/Auth"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

// NewApp инициализирует grpc + service слой + storage
func NewApp(log *slog.Logger, grpcPort int, storagePath string, TokenTTL time.Duration, auth authgrpc.Auth) *App {
	// TODO: иницилазция Storage

	// TODO: инициализцаия auth service

	// grpc—сервер инициализация
	grpcApp := grpcapp.New(log, grpcPort, auth)
	return &App{
		GRPCServer: grpcApp,
	}
}

package grpcapp

import (
	authgrpc "Auth/internal/grpc/Auth"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New инициализирует grpc—сервера
func New(log *slog.Logger, port int, auth authgrpc.Auth) *App {
	grpcServer := grpc.NewServer()

	// обработчик
	authgrpc.RegisterAuthServer(grpcServer, auth)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
	}
}

func (a *App) Start() error {
	const op = "grpcapp.Start"

	log := a.log.With(
		slog.String("place", op),
		slog.Int("port", a.port),
	)

	// создаем слушателя
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}

	// запускаем grpc—сервер
	if err := a.gRPCServer.Serve(l); err != nil {
		return err
	}

	log.Info("gRPC server started")

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log := a.log.With(
		slog.String("place", op),
		slog.Int("port", a.port),
	)

	// останавливаем grpc—сервер
	a.gRPCServer.GracefulStop()

	log.Info("gRPC server stopped")

}

// MustRun обертка над Run()
func (a *App) MustRun() {
	if err := a.Start(); err != nil {
		panic(err)
	}
}

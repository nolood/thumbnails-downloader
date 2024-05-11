package grpcapp

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	downloadergrpc "thumbs/internal/grpc/thumbs"
	"thumbs/internal/services/youtube"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *zap.Logger, port int, youtubeService youtube.Service) *App {
	gRPCServer := grpc.NewServer()

	downloadergrpc.Register(gRPCServer, youtubeService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(zap.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running on port: ", zap.Int("port", a.port))

	if err = a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(zap.String("op", op))

	a.gRPCServer.GracefulStop()
}

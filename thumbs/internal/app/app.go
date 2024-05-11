package app

import (
	"go.uber.org/zap"
	grpcapp "thumbs/internal/app/grpc"
	"thumbs/internal/services/youtube"
	"thumbs/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *zap.Logger, grpcPort int, storagePath string, youtubeKey string) *App {

	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Panic("failed to create storage", zap.Error(err))
	}

	youtubeService := youtube.New(log, storage, youtubeKey)

	gprcApp := grpcapp.New(log, grpcPort, youtubeService)

	return &App{
		GRPCSrv: gprcApp,
	}
}

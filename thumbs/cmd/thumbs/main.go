package main

import (
	"os"
	"os/signal"
	"syscall"
	"thumbs/internal/app"
	"thumbs/internal/config"
	"thumbs/internal/lib/logger"
)

func main() {

	cfg := config.MustLoad()

	log := logger.Setup(cfg.Env)

	application := app.New(log, cfg.GPRC.Port, cfg.StoragePath, cfg.YoutubeKey)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()

	log.Info("Application stopped")

}

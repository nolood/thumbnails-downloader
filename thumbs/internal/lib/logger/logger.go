package logger

import "go.uber.org/zap"

const envLocal = "local"
const envProd = "prod"

// В идеале написать обёртку для логгера, на случай если понадобится переезд на другую реализацию

func Setup(env string) *zap.Logger {
	var log *zap.Logger

	switch env {
	case envLocal:
		log, _ = zap.NewDevelopment()
	case envProd:
		log, _ = zap.NewProduction()
	default:
		log, _ = zap.NewProduction()
	}

	return log
}

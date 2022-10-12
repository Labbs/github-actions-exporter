package logging

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	if logger != nil {
		return logger
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize logger: %v", err)
	}

	return logger.Sugar()
}

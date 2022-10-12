package logging

import (
	"log"

	"github.com/spendesk/github-actions-exporter/pkg/config"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	if logger != nil {
		return logger
	}

	var (
		freshLogger *zap.Logger
		err         error
	)
	if config.LogStructured {
		freshLogger, err = zap.NewProduction()
	} else {
		freshLogger, err = zap.NewDevelopment()
	}

	defer freshLogger.Sync()

	if err != nil {
		log.Fatalf("Can't initialize logger: %v", err)
	}

	logger = freshLogger.Sugar()
	return logger
}

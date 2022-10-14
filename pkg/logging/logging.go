package logging

import (
	"log"

	"github.com/spendesk/github-actions-exporter/pkg/config"
	"go.uber.org/zap"
)

func InitLogger() *zap.SugaredLogger {
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

	return freshLogger.Sugar()
}

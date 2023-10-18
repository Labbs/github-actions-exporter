package logging

import (
	"log"

	"github.com/spendesk/github-actions-exporter/pkg/config"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

var logger *zap.SugaredLogger

func validateFormat(text string) {
	if !slices.Contains([]string{"json", "plain"}, text) {
		logger.Fatalf("Invalid log_format '%v'", text)
	}
}

func InitLogger() *zap.SugaredLogger {
	var (
		freshLogger *zap.Logger
		err         error
	)
	if config.LogFormat == "plain" {
		freshLogger, err = zap.NewDevelopment()
	} else {
		freshLogger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalf("Can't initialize logger: %v", err)
	}

	defer freshLogger.Sync()
	logger = freshLogger.Sugar()

	validateFormat(config.LogFormat)

	return logger
}

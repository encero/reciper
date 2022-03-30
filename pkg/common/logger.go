package common

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

func LoggerFromEnv() (*zap.Logger, error) {
	var loggerConfig zap.Config
	if strings.ToLower(os.Getenv("LOGGER")) == "dev" {
		loggerConfig = zap.NewDevelopmentConfig()
	} else {
		loggerConfig = zap.NewProductionConfig()
	}

	return loggerConfig.Build()
}

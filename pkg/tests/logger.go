package tests

import (
	"testing"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) *zap.Logger {
	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		t.Fatalf("test zap logger %v", err)
	}

	logger = logger.With(zap.String("system", "tests"))

	return logger
}

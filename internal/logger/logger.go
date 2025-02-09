package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func NewLogger(env string) *zap.Logger {
	var logger *zap.Logger
	var config zap.Config
	var err error

	switch env {
	case envDev:
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = config.Build()
	case envProd:
		// included atomic level: INFO
		config = zap.NewProductionConfig()
		logger, err = config.Build()
	default:
		log.Fatal("Logger is uninitialized: wrong env value")
	}

	if err != nil {
		log.Fatal("error configuring logger: " + err.Error())
	}

	return logger
}

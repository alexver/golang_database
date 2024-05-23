package logger

import (
	"fmt"
	"log"

	"github.com/alexver/golang_database/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

var supportedLoggingLevels = map[string]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
}

func CreateLogger(config *config.LoggerConfig) *zap.Logger {
	if config == nil {
		log.Fatal("there is no logger config provided")
	}

	logLevel, ok := supportedLoggingLevels[config.Level]
	if !ok {
		logLevel = zapcore.DebugLevel
	}
	loggerConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{config.Output},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	logger.Info(fmt.Sprintf("Logger is configured to apply '%s' level and to write into '%s' file", config.Level, config.Output))

	return logger
}

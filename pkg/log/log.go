package log

import (
	"strings"

	"go.uber.org/zap"
)

const (
	PROD        = "prod"
	VERBOSE     = "verbose"
	VERYVERBOSE = "veryverbose"
)

func InitLog(appEnv string) ILogger {
	appEnv = strings.ToLower(appEnv)
	if appEnv == PROD {
		return NewZapLogger(produceProdConfig())
	}
	return NewZapLogger(produceDevConfig(appEnv))
}

func produceProdConfig() zap.Config {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	config.Encoding = "console"
	return config
}

func produceDevConfig(appEnv string) zap.Config {
	config := zap.NewDevelopmentConfig()
	if appEnv == VERBOSE {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return config
}

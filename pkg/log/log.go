package log

import (
	"strings"

	"go.uber.org/zap"
)

const (
	Prod        = "prod"
	Verbose     = "verbose"
	VeryVerbose = "veryverbose"
)

func InitLog(appEnv string) ILogger {
	appEnv = strings.ToLower(appEnv)
	if appEnv == Prod {
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
	if appEnv == Verbose {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return config
}

package log

import (
	"strings"

	"go.uber.org/zap"
)

func InitLog(appEnv string) ZapLogger {
	if strings.ToLower(appEnv) == "prod" {
		return NewZapLogger(produceProdConfig())
	}
	return NewZapLogger(produceDevConfig())
}

func produceProdConfig() zap.Config {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	return config
}

func produceDevConfig() zap.Config {
	return zap.NewDevelopmentConfig()
}

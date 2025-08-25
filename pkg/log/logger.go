package log

import (
	"go.uber.org/zap"
)

// Reference: https://dave.cheney.net/2015/11/05/lets-talk-about-logging
// for philosophy on log levels

type ZapLogger struct {
	Log *zap.SugaredLogger
}

func NewZapLogger(conf zap.Config) ILogger {
	logger, err := conf.Build()

	if err != nil {
		panic("Could not setup logger.\n" + err.Error())
	}

	defer logger.Sync()

	sugar := logger.Sugar()

	return &ZapLogger{Log: sugar}
}

func (zl *ZapLogger) Debug(msg string, keysAndValues ...any) {
	zl.Log.Debugw(msg, keysAndValues...)
}

func (zl *ZapLogger) Info(msg string, keysAndValues ...any) {
	zl.Log.Infow(msg, keysAndValues...)
}

func (zl *ZapLogger) Panic(msg string, keysAndValues ...any) {
	zl.Log.Panicw(msg, keysAndValues...)
}

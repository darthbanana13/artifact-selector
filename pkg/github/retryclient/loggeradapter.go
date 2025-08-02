package retryclient

import (
	"github.com/darthbanana13/artifact-selector/pkg/log"
)

type LeveledLoggerAdapter struct {
	log log.ILogger
}

func NewLeveledLoggerAdapter(logger log.ILogger) *LeveledLoggerAdapter {
	return &LeveledLoggerAdapter{
		log: logger,
	}
}

func (l *LeveledLoggerAdapter) Debug(msg string, keysAndValues ...any) {
	l.log.Debug(msg, keysAndValues...)
}

func (l *LeveledLoggerAdapter) Info(msg string, keysAndValues ...any) {
	l.log.Info(msg, keysAndValues...)
}

func (l *LeveledLoggerAdapter) Warn(msg string, keysAndValues ...any) {
	l.log.Info(msg, keysAndValues...)
}

func (l *LeveledLoggerAdapter) Error(msg string, keysAndValues ...any) {
	l.log.Info(msg, keysAndValues...)
}

func (l *LeveledLoggerAdapter) Panic(msg string, keysAndValues ...any) {
	l.log.Panic(msg, keysAndValues...)
}

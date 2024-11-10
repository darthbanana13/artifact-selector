package log

import (
  "go.uber.org/zap"
)

type ZapLogger struct {
  log *zap.Logger
}

func NewZapLogger(conf zap.Config) ZapLogger {
  logger, err := conf.Build()

  if err != nil {
    panic("Could not setup logger.\n" + err.Error())
  }

  defer logger.Sync()

  return ZapLogger{log: logger}
}

func (zl *ZapLogger) Debug(msg string) {
  zl.log.Debug(msg)
}

func (zl *ZapLogger) Info(msg string) {
    zl.log.Info(msg)
}

func (zl *ZapLogger) Warn(msg string) {
    zl.log.Warn(msg)
}

func (zl *ZapLogger) Error(msg string) {
    zl.log.Error(msg)
}

func (zl *ZapLogger) Fatal(msg string) {
    zl.log.Fatal(msg)
}

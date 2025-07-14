package log

type ILogger interface {
  Debug(msg string, keysAndValues ...any)
  Info(msg string, keysAndValues ...any)
  Panic(msg string, keysAndValues ...any)
}

package log

type ILogger interface {
  Debug(msg string)
  Info(msg string)
  Warn(msg string)
  Error(msg string)
  Fatal(msg string)
}

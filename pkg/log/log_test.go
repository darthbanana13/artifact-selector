package log

import (
  "testing"

  "github.com/stretchr/testify/assert"

  "go.uber.org/zap"
)

func TestInitLogProd(t *testing.T) {
  logger := InitLog("ProD")
  actualLevel := logger.log.Level()
  assert.Equal(t, zap.WarnLevel, actualLevel)
}

func TestInitLogDev(t *testing.T) {
  logger := InitLog("dEv")
  actualLevel := logger.log.Level()
  assert.Equal(t, zap.DebugLevel, actualLevel)
}

// This test does not have any value
func TestNewZapLoggerPanic(t *testing.T) {
  var emptyConfig zap.Config
  assert.Panics(t, func() { NewZapLogger(emptyConfig) })
}

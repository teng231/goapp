package logger

import (
	"errors"
	"testing"

	"go.uber.org/zap"
)

func TestTestOutlog(t *testing.T) {
	Init(false)
	L().Info("test log")
	L().Warn("err", zap.Error(errors.New("looix")))
}

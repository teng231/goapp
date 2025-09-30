package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

func Init(isDev bool) {
	once.Do(func() {
		var err error
		if isDev {
			log, err = zap.NewDevelopment() // có màu, dễ debug
		} else {
			cfg := zap.NewProductionConfig()
			cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // format thời gian dễ đọc
			cfg.EncoderConfig.CallerKey = "caller"                    // giữ key caller
			log, err = cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
			// log, err = zap.NewProduction() // output JSON chuẩn
			// log = _log
		}
		if err != nil {
			panic(err)
		}
	})
}

func L() *zap.Logger {
	if log == nil {
		panic("logger not initialized, call logger.Init() first")
	}
	return log
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

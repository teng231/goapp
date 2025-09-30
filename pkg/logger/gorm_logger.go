package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

type ZapGormLogger struct {
	log        *zap.Logger
	logLevel   gormlogger.LogLevel
	slowThresh time.Duration
}

func NewZapGormLogger(log *zap.Logger) *ZapGormLogger {
	return &ZapGormLogger{
		log:        log,
		logLevel:   gormlogger.Info,
		slowThresh: 200 * time.Millisecond,
	}
}

func (l *ZapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.logLevel = level
	return &newlogger
}

func (l *ZapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Info {
		l.logWithTrace(ctx).Sugar().Infof(msg, data...)
	}
}
func (l *ZapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Warn {
		l.logWithTrace(ctx).Sugar().Warnf(msg, data...)
	}
}
func (l *ZapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Error {
		l.logWithTrace(ctx).Sugar().Errorf(msg, data...)
	}
}
func (l *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel == gormlogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	}
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	if err != nil && l.logLevel >= gormlogger.Error {
		l.log.Error("gorm query error", append(fields, zap.Error(err))...)
	} else if elapsed > l.slowThresh && l.logLevel >= gormlogger.Warn {
		l.log.Warn("slow query", fields...)
	} else if l.logLevel >= gormlogger.Info {
		l.log.Info("gorm query", fields...)
	}
}

func (l *ZapGormLogger) logWithTrace(ctx context.Context) *zap.Logger {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return l.log.With(zap.String("trace_id", traceID))
	}
	return l.log
}

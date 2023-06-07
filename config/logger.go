package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func newCustomLogger() (*zap.Logger, error) {
	os.Mkdir("./logs", os.ModePerm)
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller", // 不记录日志调用位置
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", "./logs/test.log"},
		ErrorOutputPaths: []string{"stdout", "./logs/error.log"},
	}
	return cfg.Build()
}

func Info(msg string, keysAndValues ...interface{}) {
	logger, _ := newCustomLogger()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow(msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	logger, _ := newCustomLogger()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Debugw(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	logger, _ := newCustomLogger()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Warnw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	logger, _ := newCustomLogger()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Errorw(msg, keysAndValues...)
}

func Panic(msg string, keysAndValues ...interface{}) {
	logger, _ := newCustomLogger()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Panicw(msg, keysAndValues...)
}

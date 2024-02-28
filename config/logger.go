package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Log *zap.SugaredLogger

func init() {
	logger, _ := newCustomLogger()
	defer logger.Sync()
	Log = logger.Sugar()
}

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
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", "./logs/test.log"},
		ErrorOutputPaths: []string{"stdout", "./logs/error.log"},
	}
	return cfg.Build()
}

package logger

import (
	"os"
	"strings"

	"github.com/Vikktttoriya/flight-tracker/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(cfg config.LoggerConfig) (*zap.Logger, error) {
	level := parseLevel(cfg.Level)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		MessageKey:    "msg",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	var jsonEncoder = zapcore.NewJSONEncoder(encoderConfig)

	consoleEncodeConfig := encoderConfig
	consoleEncodeConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	var consoleEncoder = zapcore.NewConsoleEncoder(consoleEncodeConfig)

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename: cfg.FilePath,
		MaxAge:   int(cfg.MaxAge.Hours() / 24),
		Compress: true,
	})

	consoleWriter := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, fileWriter, level),
		zapcore.NewCore(consoleEncoder, consoleWriter, level),
	)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCallerSkip(1),
	)

	return logger, nil
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

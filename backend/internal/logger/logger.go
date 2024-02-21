package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	loggerOutputDir      = "../logs"
	loggerOutputFormat   = "stdout"
	loggerOutputFilePath = "../logs/todo_list.log"
)

type logger struct {
	*zap.SugaredLogger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	With(args ...interface{}) *zap.SugaredLogger
	Sync() error
}

func New() (Logger, error) {
	if err := createLogsDirectory(); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Создание конфигурации логгера
	config := zap.Config{
		Encoding:    "console",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{loggerOutputFormat, loggerOutputFilePath},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "message",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
	}

	zapLogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании логгера - %s", err)
	}

	sugar := zapLogger.Sugar()

	return &logger{
		SugaredLogger: sugar,
	}, nil
}

func createLogsDirectory() error {
	if _, err := os.Stat(loggerOutputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(loggerOutputDir, os.ModePerm); err != nil {
			return fmt.Errorf("ошибка при создании директории - %s", err)
		}
	}

	return nil
}

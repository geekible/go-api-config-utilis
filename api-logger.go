package goapiconfigutilis

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ApiLogger struct {
	logFilePath string
	logger      *zap.SugaredLogger
}

func InitApiLogger(filePath string) *ApiLogger {
	return &ApiLogger{
		logFilePath: filePath,
	}
}

func (l *ApiLogger) BuilderLogger() error {
	logFile, err := os.Create(l.logFilePath)
	if err != nil {
		return fmt.Errorf("unable to create %s log file with error %v", l.logFilePath, err)
	}

	pe := zap.NewProductionEncoderConfig()
	fileEncoder := zapcore.NewJSONEncoder(pe)
	pe.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	level := zap.InfoLevel

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	logger := zap.New(core)
	l.logger = logger.Sugar()

	return nil
}

func (l *ApiLogger) LogInfo(source, message string) {
	go func() {
		l.logger.Infof("Source: %s\nMessage: %s", source, message)
	}()
}

func (l *ApiLogger) LogWarning(source, message string) {
	go func() {
		l.logger.Warnf("Source: %s\nMessage: %s", source, message)
	}()
}

func (l *ApiLogger) LogError(source string, err error) {
	go func() {
		l.logger.Errorf("Source: %s\nError: %v", source, err)
	}()
}

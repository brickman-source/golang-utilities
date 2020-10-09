/*
 * Copyright (c) 2019. Brickman Source.
 */

package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogFormat string

const (
	LogFormatJSON LogFormat = "json"
	LogFormatText LogFormat = "text"
)

type FileLogger struct {
	fields    Fields
	zapLogger *zap.Logger
}

func NewFileLogger(filePath string, format LogFormat, maxSize int, maxBackups int, maxAge int) *FileLogger {
	hook := lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize, // megabytes
		MaxBackups: maxBackups,
		MaxAge:     maxAge, //days
		Compress:   false,  // disabled by default
	}

	levelFunc := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zap.InfoLevel
	})

	fileWriter := zapcore.AddSync(&hook)

	var fileEncoder zapcore.Encoder
	if format == LogFormatJSON {
		fileEncoder = zapcore.NewJSONEncoder(newProductionEncoderConfig())
	} else if format == LogFormatText {
		fileEncoder = zapcore.NewConsoleEncoder(newConsoleEncoderConfig())
	}
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, levelFunc)

	core := zapcore.NewTee(
		fileCore,
	)

	return &FileLogger{
		zapLogger: zap.New(core),
	}
}

func (fileLogger *FileLogger) AppName() string {
	return ""
}
func (fileLogger *FileLogger) ZapLogger() *zap.Logger {
	return fileLogger.zapLogger
}
func (fileLogger *FileLogger) With(fields Fields) *FileLogger {
	logger := fileLogger.clone()
	logger.fields = fields
	return logger
}
func (fileLogger *FileLogger) Panicf(fmt string, args ...interface{}) {
	write(fileLogger, LevelPanic, fileLogger.fields, fmt, args...)
}
func (fileLogger *FileLogger) Fatalf(fmt string, args ...interface{}) {
	write(fileLogger, LevelFatal, fileLogger.fields, fmt, args...)
}
func (fileLogger *FileLogger) Errorf(fmt string, args ...interface{}) {
	write(fileLogger, LevelError, fileLogger.fields, fmt, args...)
}
func (fileLogger *FileLogger) Infof(fmt string, args ...interface{}) {
	write(fileLogger, LevelInfo, fileLogger.fields, fmt, args...)
}
func (fileLogger *FileLogger) Debugf(fmt string, args ...interface{}) {
	write(fileLogger, LevelDebug, fileLogger.fields, fmt, args...)
}
func (fileLogger *FileLogger) Warnf(fmt string, args ...interface{}) {
	write(fileLogger, LevelWarn, fileLogger.fields, fmt, args...)
}
func (fileLogger *FileLogger) Messagef(level Level, fmt string, args ...interface{}) {
	write(fileLogger, level, fileLogger.fields, fmt, args...)
}
func (fileLogger *FileLogger) Flush() {
	if fileLogger.zapLogger != nil {
		fileLogger.zapLogger.Sync()
	}
}

func (fileLogger *FileLogger) clone() *FileLogger {
	c := *fileLogger
	return &c
}

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
		zap.New(core),
	}
}

func (fileLogger *FileLogger) Panicf(fields Fields, fmt string, args ...interface{}) {
	write(fileLogger.zapLogger, LevelPanic, fields, fmt, args...)
}
func (fileLogger *FileLogger) Fatalf(fields Fields, fmt string, args ...interface{}) {
	write(fileLogger.zapLogger, LevelFatal, fields, fmt, args...)
}
func (fileLogger *FileLogger) Errorf(fields Fields, fmt string, args ...interface{}) {
	write(fileLogger.zapLogger, LevelError, fields, fmt, args...)
}
func (fileLogger *FileLogger) Infof(fields Fields, fmt string, args ...interface{}) {
	write(fileLogger.zapLogger, LevelInfo, fields, fmt, args...)
}
func (fileLogger *FileLogger) Debugf(fields Fields, fmt string, args ...interface{}) {
	write(fileLogger.zapLogger, LevelDebug, fields, fmt, args...)
}
func (fileLogger *FileLogger) Warnf(fields Fields, fmt string, args ...interface{}) {
	write(fileLogger.zapLogger, LevelWarn, fields, fmt, args...)
}
func (fileLogger *FileLogger) Messagef(fields Fields, level Level, fmt string, args ...interface{}) {
	write(fileLogger.zapLogger, level, fields, fmt, args...)
}
func (fileLogger *FileLogger) Flush() {
	if fileLogger.zapLogger != nil {
		fileLogger.zapLogger.Sync()
	}
}

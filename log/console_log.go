package log

import (
	"fmt"
	"github.com/brickman-source/golang-utilities/converter"
	"github.com/brickman-source/golang-utilities/json"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Level string
type Format string

type Logger struct {
	fields Fields
	zapLogger *zap.Logger
	appName   string
}

var (
	defaultLogger *Logger
)

const (
	FormatJSON Format = "json"
	FormatText Format = "text"

	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelFatal Level = "fatal"
	LevelPanic Level = "panic"
)

func init() {
	Config("golang-util",
		LevelDebug,
		true,
		"",
		LogFormatJSON,
		0,
		0,
		0,
	)
}

func Config(appName string,
	level Level,
	console bool,
	filename string,
	format LogFormat,
	maxSize int,
	maxBackups int,
	maxAge int) {

	if defaultLogger == nil {
		defaultLogger = &Logger{}
	}

	defaultLogger.appName = appName

	levels := map[Level]zapcore.Level{
		LevelDebug: zap.DebugLevel,
		LevelInfo:  zap.InfoLevel,
		LevelWarn:  zap.WarnLevel,
		LevelError: zap.ErrorLevel,
		LevelFatal: zap.PanicLevel,
		LevelPanic: zap.FatalLevel,
	}

	zaplvl := levels[level]

	levelFunc := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zaplvl
	})

	cores := make([]zapcore.Core, 0)

	if console {
		consoleDebugging := zapcore.Lock(os.Stdout)
		consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:       "T",
			LevelKey:      "L",
			NameKey:       "N",
			CallerKey:     "C",
			MessageKey:    "M",
			StacktraceKey: "S",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Local().Format("2006-01-02 15:04:05"))
			},
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
		consoleCore := zapcore.NewCore(consoleEncoder, consoleDebugging, levelFunc)
		cores = append(cores, consoleCore)
	}
	if filename != "" {
		hook := lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize, // megabytes
			MaxBackups: maxBackups,
			MaxAge:     maxAge, //days
			Compress:   false,  // disabled by default
		}

		fileWriter := zapcore.AddSync(&hook)

		var fileEncoder zapcore.Encoder
		if format == LogFormatJSON {
			fileEncoder = zapcore.NewJSONEncoder(newProductionEncoderConfig())
		} else if format == LogFormatText {
			fileEncoder = zapcore.NewConsoleEncoder(newConsoleEncoderConfig())
		}
		fileCore := zapcore.NewCore(fileEncoder, fileWriter, levelFunc)

		cores = append(cores, fileCore)
	}

	core := zapcore.NewTee(
		cores...,
	)

	logger := zap.New(core, zap.AddStacktrace(zap.ErrorLevel))

	if defaultLogger.zapLogger != nil {
		defaultLogger.zapLogger.Sync()
	}
	defaultLogger.zapLogger = logger
}
func With(fields Fields) *Logger {
	logger := defaultLogger.clone()
	logger.fields = fields
	return logger
}
func Panicf(fmt string, args ...interface{}) {
	write(defaultLogger, LevelPanic, nil, fmt, args...)
}
func Fatalf(fmt string, args ...interface{}) {
	write(defaultLogger, LevelFatal, nil, fmt, args...)
}
func Errorf(fmt string, args ...interface{}) {
	write(defaultLogger, LevelError, nil, fmt, args...)
}
func Infof(fmt string, args ...interface{}) {
	write(defaultLogger, LevelInfo, nil, fmt, args...)
}
func Debugf(fmt string, args ...interface{}) {
	write(defaultLogger, LevelDebug, nil, fmt, args...)
}
func Warnf(fmt string, args ...interface{}) {
	write(defaultLogger, LevelWarn, nil, fmt, args...)
}
func Messagef(level Level, fmt string, args ...interface{}) {
	write(defaultLogger, level, nil, fmt, args...)
}
func Json(j interface{}) {
	jd := json.ShouldMarshal(j)
	if jd == nil {
		return
	}
	write(defaultLogger, LevelInfo, nil, "%s", string(jd))
}
func Flush() {
	if defaultLogger.zapLogger != nil {
		defaultLogger.zapLogger.Sync()
	}
}
func (log *Logger) AppName() string {
	return log.appName
}
func (log *Logger) ZapLogger() *zap.Logger {
	return log.zapLogger
}
func (log *Logger) Panicf(fmt string, args ...interface{}) {
	write(log, LevelPanic, log.fields, fmt, args...)
}
func (log *Logger) Fatalf(fmt string, args ...interface{}) {
	write(log, LevelFatal, log.fields, fmt, args...)
}
func (log *Logger) Errorf(fmt string, args ...interface{}) {
	write(log, LevelError, log.fields, fmt, args...)
}
func (log *Logger) Infof(fmt string, args ...interface{}) {
	write(log, LevelInfo, log.fields, fmt, args...)
}
func (log *Logger) Debugf(fmt string, args ...interface{}) {
	write(log, LevelDebug, log.fields, fmt, args...)
}
func (log *Logger) Warnf(fmt string, args ...interface{}) {
	write(log, LevelWarn, log.fields, fmt, args...)
}
func (log *Logger) Messagef(level Level, fmt string, args ...interface{}) {
	write(log, level, log.fields, fmt, args...)
}
func (log *Logger) Json(j interface{}) {
	jd := json.ShouldMarshal(j)
	if jd == nil {
		return
	}
	write(log, LevelInfo, log.fields, "%s", string(jd))
}
func (log *Logger) Flush() {
	if log.zapLogger != nil {
		_ = log.zapLogger.Sync()
	}
}

func DumpFormattedJson(j interface{}) {
	jd, err := json.MarshalIndent(j, "", " ")
	if err != nil {
		return
	}
	println(string(jd))
}
func DumpKeyValue(j interface{}) {
	m := converter.ConvertToMap(j)
	for key, value := range m {
		fmt.Printf("%s:%v\n", key, value)
	}
}

func newProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
func newConsoleEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Local().Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func (log *Logger) clone() *Logger {
	c := *log
	return &c
}

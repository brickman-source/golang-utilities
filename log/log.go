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

var (
	globalZapLogger *zap.Logger
	globalAppName string
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
	globalAppName = appName

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

	if globalZapLogger != nil {
		globalZapLogger.Sync()
	}
	globalZapLogger = logger
}

func Panicf(fields Fields, fmt string, args ...interface{}) {
	write(globalZapLogger, LevelPanic, fields, fmt, args...)
}
func Fatalf(fields Fields, fmt string, args ...interface{}) {
	write(globalZapLogger, LevelFatal, fields, fmt, args...)
}
func Errorf(fields Fields, fmt string, args ...interface{}) {
	write(globalZapLogger, LevelError, fields, fmt, args...)
}
func Infof(fields Fields, fmt string, args ...interface{}) {
	write(globalZapLogger, LevelInfo, fields, fmt, args...)
}
func Debugf(fields Fields, fmt string, args ...interface{}) {
	write(globalZapLogger, LevelDebug, fields, fmt, args...)
}
func Warnf(fields Fields, fmt string, args ...interface{}) {
	write(globalZapLogger, LevelWarn, fields, fmt, args...)
}

func Messagef(fields Fields, level Level, fmt string, args ...interface{}) {
	write(globalZapLogger, level, fields, fmt, args...)
}
func Flush() {
	if globalZapLogger != nil {
		globalZapLogger.Sync()
	}
}
func Json(j interface{}) {
	jd := json.ShouldMarshal(j)
	if jd == nil {
		return
	}
	write(globalZapLogger, LevelInfo, nil,"%s", string(jd))
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
func write(zapLogger *zap.Logger, level Level, fields Fields, format string, args ...interface{}) {
	if zapLogger == nil {
		zapLogger = globalZapLogger
	}
	if zapLogger == nil {
		return
	}

	zapFields := make([]zap.Field, 0)
	if fields != nil {
		for k, v := range fields {
			switch v.(type) {
			case map[string]interface{}, []interface{}, []map[string]interface{}:
				zapFields = append(zapFields, zap.String(k, string(json.ShouldMarshal(v))))
			default:
				zapFields = append(zapFields, zap.Any(k, v))
			}

		}
	}
	if globalAppName != "" {
		zapFields = append(zapFields, zap.String("app", globalAppName))
	}
	if level == LevelPanic {
		zapLogger.Panic(fmt.Sprintf(format, args...), zapFields...)
		os.Exit(1)
	} else if level == LevelFatal {
		zapLogger.Fatal(fmt.Sprintf(format, args...), zapFields...)
		os.Exit(1)
	} else if level == LevelError {
		zapLogger.Error(fmt.Sprintf(format, args...), zapFields...)
	} else if level == LevelWarn {
		zapLogger.Warn(fmt.Sprintf(format, args...), zapFields...)
	} else if level == LevelInfo {
		zapLogger.Info(fmt.Sprintf(format, args...), zapFields...)
	} else if level == LevelDebug {
		zapLogger.Debug(fmt.Sprintf(format, args...), zapFields...)
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

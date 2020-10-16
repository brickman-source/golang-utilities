/*
 * Copyright (c) 2020. Brickman Source.
 */

package log

import (
	"fmt"
	"github.com/brickman-source/golang-utilities/json"
	"go.uber.org/zap"
	"os"
)

type logger interface {
	AppName() string
	ZapLogger() *zap.Logger
}

func write(log logger, level Level, fields Fields, format string, args ...interface{}) {
	zapLogger := log.ZapLogger()
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
	if log.AppName() != "" {
		zapFields = append(zapFields, zap.String("app", log.AppName()))
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

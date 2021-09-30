package map_util

import (
	"reflect"
	"strings"

	"github.com/brickman-source/golang-utilities/converter"
)

func SetValue(m map[string]interface{}, path string, val interface{}) {
	if path == "" {
		return
	}

	index := strings.Index(path, ".")
	if index > 0 {
		pp := path[:index]
		child := m[pp]
		if child == nil || reflect.ValueOf(child).IsNil() {
			child = make(map[string]interface{})
			m[pp] = child
		}
		if childM, ok := child.(map[string]interface{}); ok {
			SetValue(childM, path[index+1:], val)
		}
		return
	} else {
		m[path] = val
	}
}
func GetValue(m map[string]interface{}, path string) interface{} {
	if m == nil || reflect.ValueOf(m).IsNil() {
		return nil
	}
	if path == "" {
		return m
	}

	index := strings.Index(path, ".")
	if index > 0 {
		child := converter.ConvertToMap(m[path[:index]])
		if child == nil {
			return nil
		}
		return GetValue(child, path[index+1:])
	} else {
		return m[path]
	}
}

func Exists(m map[string]interface{}, path string) bool {
	if m == nil || reflect.ValueOf(m).IsNil() {
		return false
	}
	if path == "" {
		return false
	}
	index := strings.Index(path, ".")
	if index > 0 {
		child := converter.ConvertToMap(m[path[:index]])
		if child == nil {
			return false
		}
		return Exists(child, path[index+1:])
	} else {
		if _, ok := m[path]; ok {
			return true
		}
		return false
	}
}

func GetValueAsInt(m map[string]interface{}, path string, defaultValue int) int {
	if m == nil || reflect.ValueOf(m).IsNil() {
		return defaultValue
	}
	if path == "" {
		return converter.AsInt(m, defaultValue)
	}
	index := strings.Index(path, ".")
	if index > 0 {
		child := converter.ConvertToMap(m[path[:index]])
		if child == nil {
			return defaultValue
		}
		return GetValueAsInt(child, path[index+1:], defaultValue)
	} else {
		return converter.AsInt(m[path], defaultValue)
	}
}
func GetValueAsInt32(m map[string]interface{}, path string, defaultValue int32) int32 {
	if m == nil || reflect.ValueOf(m).IsNil() {
		return defaultValue
	}
	if path == "" {
		return converter.AsInt32(m, defaultValue)
	}

	index := strings.Index(path, ".")
	if index > 0 {
		child := converter.ConvertToMap(m[path[:index]])
		if child == nil {
			return defaultValue
		}
		return GetValueAsInt32(child, path[index+1:], defaultValue)
	} else {
		return converter.AsInt32(m[path], defaultValue)
	}
}
func GetValueAsInt64(m map[string]interface{}, path string, defaultValue int64) int64 {
	if m == nil || reflect.ValueOf(m).IsNil() {
		return defaultValue
	}
	if path == "" {
		return converter.AsInt64(m, defaultValue)
	}

	index := strings.Index(path, ".")
	if index > 0 {
		child := converter.ConvertToMap(m[path[:index]])
		if child == nil {
			return defaultValue
		}
		return GetValueAsInt64(child, path[index+1:], defaultValue)
	} else {
		return converter.AsInt64(m[path], defaultValue)
	}
}
func GetValueAsString(m map[string]interface{}, path string, defaultValue string) string {
	if m == nil || reflect.ValueOf(m).IsNil() {
		return defaultValue
	}
	if path == "" {
		return converter.AsString(m, defaultValue)
	}

	index := strings.Index(path, ".")
	if index > 0 {
		child := converter.ConvertToMap(m[path[:index]])
		if child == nil {
			return defaultValue
		}
		return GetValueAsString(child, path[index+1:], defaultValue)
	} else {
		return converter.AsString(m[path], defaultValue)
	}
}
func GetValueAsStringSlice(m map[string]interface{}, path string, defaultValue []string) []string {
	if m == nil || reflect.ValueOf(m).IsNil() {
		return defaultValue
	}
	if path == "" {
		return converter.AsStringSlice(m, defaultValue)
	}

	index := strings.Index(path, ".")
	if index > 0 {
		child := converter.ConvertToMap(m[path[:index]])
		if child == nil {
			return defaultValue
		}
		return GetValueAsStringSlice(child, path[index+1:], defaultValue)
	} else {
		return converter.AsStringSlice(m[path], defaultValue)
	}
}
func GetValueAsBool(m map[string]interface{}, path string, defaultValue bool) bool {
	if m == nil || reflect.ValueOf(m).IsNil() {
		return defaultValue
	}
	if path == "" {
		return converter.AsBool(m, defaultValue)
	}

	index := strings.Index(path, ".")
	if index > 0 {
		child := converter.ConvertToMap(m[path[:index]])
		if child == nil {
			return defaultValue
		}
		return GetValueAsBool(child, path[index+1:], defaultValue)
	} else {
		return converter.AsBool(m[path], defaultValue)
	}
}
func GetValueAsFloat64(m map[string]interface{}, path string, defaultValue float64) float64 {
	if m == nil || reflect.ValueOf(m).IsNil() {
		return defaultValue
	}
	if path == "" {
		return converter.AsFloat64(m, defaultValue)
	}

	index := strings.Index(path, ".")
	if index > 0 {
		child := converter.ConvertToMap(m[path[:index]])
		if child == nil {
			return defaultValue
		}
		return GetValueAsFloat64(child, path[index+1:], defaultValue)
	} else {
		return converter.AsFloat64(m[path], defaultValue)
	}
}
func GetValueAsMap(m map[string]interface{}, path string, defaultValue map[string]interface{}) map[string]interface{} {
	if m == nil || reflect.ValueOf(m).IsNil() {
		return defaultValue
	}
	if path == "" {
		return m
	}

	index := strings.Index(path, ".")
	if index > 0 {
		child := converter.ConvertToMap(m[path[:index]])
		if child == nil {
			return defaultValue
		}
		return GetValueAsMap(child, path[index+1:], defaultValue)
	} else {
		child := m[path]
		switch child.(type) {
		case map[string]interface{}:
			return child.(map[string]interface{})
		default:
			m := converter.ConvertToMap(child)
			if m != nil {
				return m
			}
			return defaultValue
		}
	}
}

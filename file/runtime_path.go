package file

import (
	"path"
	"runtime"
)

func GetFilePathFromRuntime(part string) string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), part)
}

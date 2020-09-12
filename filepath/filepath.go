package filepath

import (
	"os"
	goFilePath "path/filepath"
	"runtime"
)

func ExpandUser(s string) string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	home := os.Getenv(env)
	if home == "" {
		return s
	}

	if len(s) >= 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		if runtime.GOOS == "windows" {
			s = goFilePath.ToSlash(goFilePath.Join(home, s[2:]))
		} else {
			s = goFilePath.Join(home, s[2:])
		}
	}
	return os.Expand(s, func(env string) string {
		if env == "HOME" {
			return home
		}
		return os.Getenv(env)
	})
}

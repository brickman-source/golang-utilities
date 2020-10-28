package text

import "strings"

func IsEmpty(str string) bool {
	if strings.TrimSpace(str) == "" {
		return true
	}
	return false
}

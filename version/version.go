package version

import (
	"fmt"
)

var (
	buildTime string
	gitHash   string
)

func String(app string, version string) string {
	return fmt.Sprintf("%s v%s (hash %s build %s)", app, version, gitHash, buildTime)
}

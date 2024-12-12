package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "0.1"
	Revision  string
	BuildDate string
	GoVersion = runtime.Version()
	GoOS      = runtime.GOOS
	GoArch    = runtime.GOARCH
)

func Info() string {
	return fmt.Sprintf("(version=%s, revision=%s)", Version, Revision)
}

func BuildContext() string {
	return fmt.Sprintf("(go=%s, platform=%s, date=%s)", GoVersion, GoOS+"/"+GoArch, BuildDate)
}

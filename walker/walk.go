package walker

import (
	"os"
	"strings"

	"github.com/sf9133/fanal/analyzer"

	"github.com/sf9133/fanal/utils"
)

var (
	ignoreDirs       = []string{".git", "node_modules", "vendor"}
	ignoreSystemDirs = []string{"proc", "sys"}
)

type WalkFunc func(filePath string, info os.FileInfo, opener analyzer.Opener) error

func isIgnored(filePath string) bool {
	filePath = strings.TrimLeft(filePath, "/")
	for _, path := range strings.Split(filePath, utils.PathSeparator) {
		if utils.StringInSlice(path, ignoreDirs) {
			return true
		}
	}

	// skip system directories such as /sys and /proc
	for _, ignore := range ignoreSystemDirs {
		if strings.HasPrefix(filePath, ignore) {
			return true
		}
	}

	return false
}

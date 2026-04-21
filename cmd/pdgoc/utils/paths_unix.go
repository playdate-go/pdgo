//go:build !windows

package utils

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/playdate-go/pdgo/cmd/pdgoc/scripts"
)

func GetBuildScriptFilename() string {
	return "device-build-*.sh"
}

func GetBuildScript() []byte {
	return scripts.DeviceBuildScriptUnix
}

func GetExecutable(path string) string {
	return path
}

func GetLs(path string) (string, []string) {
	return "ls", []string{path}
}

func GetShellExecutableName() string {
	return "bash"
}

func GetTinyGoPath() string {
	return path.Join(GetTinyGoDir(), "build/tinygo")
}

func FindPlaydatePort() (string, error) {
	var patterns = PlaydatePortPatterns()

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		if len(matches) > 0 {
			// Prefer ones with "PD" in name
			for _, m := range matches {
				if strings.Contains(m, "PD") {
					return m, nil
				}
			}
			return matches[0], nil
		}
	}

	return "", fmt.Errorf("no Playdate device found (is it connected and unlocked?)")
}

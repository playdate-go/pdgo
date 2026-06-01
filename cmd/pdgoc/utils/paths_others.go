//go:build !windows && !darwin

package utils

import (
	"fmt"
	"path"
)

func GetPlaydateSDKFallbackPath() (string, error) {
	return "", ErrNoFallbackValue
}

func GetSimulatorPath() (string, error) {
	sdkPath, err := GetPlaydateSDKPath()
	if err != nil {
		return "", err
	}

	return path.Join(sdkPath, "bin/PlaydateSimulator"), nil
}

func GetLibrary(path string) string {
	return fmt.Sprintf("%s.so", path)
}

func PlaydatePortPatterns() []string {
	return []string{
		"/dev/ttyACM*",
		"/dev/ttyUSB*",
	}
}

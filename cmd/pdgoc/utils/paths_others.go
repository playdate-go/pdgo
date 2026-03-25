//go:build !windows && !darwin

package utils

import (
	"fmt"
	"path"
)

func GetPlayDateSDKFallbackPath() (string, error) {
	return "", ErrNoFallbackValue
}

func GetSimulatorPath() (string, error) {
	sdkPath, err := GetPlayDateSDKPath()
	if err != nil {
		return "", err
	}

	return path.Join(sdkPath, "bin/PlaydateSimulator"), nil
}

func GetExecutable(path string) string {
	return path
}

func GetLibrary(path string) string {
	return fmt.Sprintf("%s.so", path)
}

func GetLs() string {
	return "ls"
}

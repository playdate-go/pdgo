package utils

import (
	"fmt"
	"os"
	"path"
)

func GetPlayDateSDKFallbackPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine fallback PlayDateSDK path: %w", err)
	}

	return path.Join(homeDir, "Developer/PlaydateSDK"), nil
}

func GetSimulatorPath() (string, error) {
	sdkPath, err := GetPlayDateSDKPath()
	if err != nil {
		return "", err
	}

	return path.Join(sdkPath, "bin/Playdate Simulator.app/Contents/MacOS/Playdate Simulator"), nil
}

func GetExecutable(path string) string {
	return path
}

func GetLibrary(path string) string {
	return fmt.Sprintf("%s.dylib", path)
}

func GetLs() string {
	return "ls"
}

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

	return path.Join(homeDir, "Documents/PlaydateSDK"), nil
}

func GetSimulatorPath() (string, error) {
	sdkPath, err := GetPlayDateSDKPath()
	if err != nil {
		return "", err
	}

	return path.Join(sdkPath, "bin/PlaydateSimulator.exe"), nil
}

func GetExecutable(path string) string {
	return fmt.Sprintf("%s.exe", path)
}

func GetLibrary(path string) string {
	return fmt.Sprintf("%s.dll", path)
}

func GetLs() string {
	return ""
}

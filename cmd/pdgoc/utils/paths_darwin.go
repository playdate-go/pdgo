package utils

import (
	"fmt"
	"os"
	"path"
)

func GetPlaydateSDKFallbackPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine fallback PlaydateSDK path: %w", err)
	}

	return path.Join(homeDir, "Developer/PlaydateSDK"), nil
}

func GetSimulatorPath() (string, error) {
	sdkPath, err := GetPlaydateSDKPath()
	if err != nil {
		return "", err
	}

	return path.Join(sdkPath, "bin/Playdate Simulator.app/Contents/MacOS/Playdate Simulator"), nil
}

func GetLibrary(path string) string {
	return fmt.Sprintf("%s.dylib", path)
}

func PlaydatePortPatterns() []string {
	return []string{
		"/dev/cu.usbmodemPD*",
		"/dev/cu.usbmodem*",
	}
}

package utils

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

const PlayDateSDKPathEnvVar = "PLAYDATE_SDK_PATH"

var ErrPlayDateSDKPathNotSet = fmt.Errorf("%s isn't set and no fallback value is available, please set env var to SDK path", PlayDateSDKPathEnvVar)

func GetPlayDateSDKPath() (string, error) {
	envVar := os.Getenv(PlayDateSDKPathEnvVar)

	if envVar != "" {
		return envVar, nil
	}

	switch runtime.GOOS {
	case "windows":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return path.Join(homeDir, "Documents/PlaydateSDK"), nil
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return path.Join(homeDir, "Developer/PlaydateSDK"), nil
	}

	return "", ErrPlayDateSDKPathNotSet
}

func GetPdcPath() (string, error) {
	sdkPath, err := GetPlayDateSDKPath()
	if err != nil {
		return "", err
	}

	pdcPath := "bin/pdc"
	if runtime.GOOS == "windows" {
		pdcPath = "bin/pdc.exe"
	}

	return path.Join(sdkPath, pdcPath), nil
}

func GetLsExec() string {
	if runtime.GOOS == "windows" {
		return ""
	}

	return "ls"
}

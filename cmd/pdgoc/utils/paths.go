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

func GetSimulatorPath() (string, error) {
	sdkPath, err := GetPlayDateSDKPath()
	if err != nil {
		return "", err
	}

	simulatorPath := "bin/PlaydateSimulator"
	switch runtime.GOOS {
	case "darwin":
		simulatorPath = "bin/Playdate Simulator.app/Contents/MacOS/Playdate Simulator"
	case "windows":
		simulatorPath = "bin/PlaydateSimulator.exe"
	}

	return path.Join(sdkPath, simulatorPath), nil
}

func GetLsExec() string {
	if runtime.GOOS == "windows" {
		return ""
	}

	return "ls"
}

func GetTinyGoPath() string {
	if runtime.GOOS == "windows" {
		return path.Join(GetTinyGoDir(), "bin/tinygo")
	}
	return path.Join(GetTinyGoDir(), "build/tinygo")
}

func GetTinyGoDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return path.Join(home, "tinygo-playdate")
}

func GetShellExecutableName() string {
	if runtime.GOOS == "windows" {
		return "powershell.exe"
	}
	return "bash"
}

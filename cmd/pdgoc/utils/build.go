package utils

import (
	"runtime"

	"github.com/playdate-go/pdgo/cmd/pdgoc/scripts"
)

func GetBuildScriptFilename() string {
	if runtime.GOOS == "windows" {
		return "device-build-*.ps1"
	}
	return "device-build-*.sh"
}

func GetBuildScript() []byte {
	if runtime.GOOS == "windows" {
		return scripts.DeviceBuildScriptWindows
	}
	return scripts.DeviceBuildScriptUnix
}

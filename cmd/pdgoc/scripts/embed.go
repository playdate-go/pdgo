package scripts

import (
	_ "embed"
)

//go:embed DeviceBuildScriptUnix.sh
var DeviceBuildScriptUnix []byte

//go:embed DeviceBuildScriptWindows.ps1
var DeviceBuildScriptWindows []byte

package utils

import (
	"errors"
	"runtime"
)

var ErrUnsupportedPlatform = errors.New("unsupported platform")

func HostLibExt() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return "so", nil
	case "darwin":
		return "dylib", nil
	case "windows":
		return "dll", nil
	default:
		return "", ErrUnsupportedPlatform
	}
}

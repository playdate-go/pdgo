package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
)

func GetPlayDateSDKPath() (string, error) {
	envVar := os.Getenv(PlayDateSDKPathEnvVar)

	if envVar != "" {
		return envVar, nil
	}

	fallback, err := GetPlayDateSDKFallbackPath()
	if err != nil {
		return "", ErrPlayDateSDKPathNotSet
	}

	return fallback, nil
}

func CheckPlayDateSDKVersion(sdkPath string) error {
	versionFilePath := filepath.Join(sdkPath, "VERSION.txt")
	if stat, err := os.Stat(sdkPath); err == nil && stat.IsDir() {
		log.Printf("auto-found SDK path: %s", sdkPath)

		b, err := os.ReadFile(versionFilePath)
		if err != nil {
			return err
		}
		verStr := string(b)
		verStr = verStr[:len(verStr)-1]

		current, err := semver.NewVersion(verStr)
		if err != nil {
			return fmt.Errorf("invalid semver in VERSION.txt '%s': %v", verStr, err)
		}

		minReq, _ := semver.NewVersion("3.0.2")

		if current.LessThan(minReq) {
			return fmt.Errorf("current SDK version %s is less than required 3.0.2", current.String())
		}

		return nil
	}

	return fmt.Errorf("failed to get sdk version information")
}

func GetPdcPath() (string, error) {
	sdkPath, err := GetPlayDateSDKPath()
	if err != nil {
		return "", err
	}

	return path.Join(sdkPath, GetExecutable("bin/pdc")), nil
}

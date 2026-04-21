package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
)

func GetPlaydateSDKPath() (string, error) {
	envVar := os.Getenv(PlaydateSDKPathEnvVar)

	if envVar != "" {
		return envVar, nil
	}

	fallback, err := GetPlaydateSDKFallbackPath()
	if err != nil {
		return "", ErrPlaydateSDKPathNotSet
	}

	return fallback, nil
}

func CheckPlaydateSDKVersion(sdkPath string) error {
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
	sdkPath, err := GetPlaydateSDKPath()
	if err != nil {
		return "", err
	}

	return path.Join(sdkPath, GetExecutable("bin/pdc")), nil
}

func GetPdutilPath() (string, error) {
	sdkPath, err := GetPlaydateSDKPath()
	if err != nil {
		return "", err
	}

	return path.Join(sdkPath, GetExecutable("bin/pdutil")), nil
}

func GetTinyGoDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return path.Join(home, "tinygo-playdate")
}

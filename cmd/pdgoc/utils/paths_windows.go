package utils

import (
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/playdate-go/pdgo/cmd/pdgoc/scripts"
	"golang.org/x/sys/windows/registry"
)

func GetPlaydateSDKFallbackPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine fallback PlaydateSDK path: %w", err)
	}

	return path.Join(homeDir, "Documents/PlaydateSDK"), nil
}

func GetSimulatorPath() (string, error) {
	sdkPath, err := GetPlaydateSDKPath()
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

func GetLs(path string) (string, []string) {
	return "cmd", []string{"/c", "dir", "/b", "/w", path}
}

func GetBuildScriptFilename() string {
	return "device-build-*.ps1"
}

func GetBuildScript() []byte {
	return scripts.DeviceBuildScriptWindows
}

func GetShellExecutableName() string {
	return "powershell.exe"
}

func GetTinyGoPath() string {
	return path.Join(GetTinyGoDir(), "bin/tinygo")
}

func FindPlaydatePort() (string, error) {
	helperMsg := "ensure your playdate is connected and unlocked"
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Enum\\USB\\VID_1331&PID_5740", registry.READ)
	if err != nil {
		return "", fmt.Errorf("failed to open registry, %s: %w", helperMsg, err)
	}
	defer key.Close()

	devices, err := key.ReadSubKeyNames(1)
	if err != nil && len(devices) == 1 {
		return "", fmt.Errorf("failed to find device info, %s: %w", helperMsg, err)
	}

	portKey, err := registry.OpenKey(key, fmt.Sprintf("%s\\Device Parameters", devices[0]), registry.READ)
	if err != nil {
		return "", fmt.Errorf("failed to read device port info, %s: %w", helperMsg, err)
	}
	defer portKey.Close()

	port, _, err := portKey.GetStringValue("PortName")
	if err = testCOMPort(port); err != nil {
		return "", fmt.Errorf("failed to connect to Playdate, %s: %w", helperMsg, err)
	}

	return port, nil
}

func testCOMPort(port string) error {
	portU16, _ := syscall.UTF16PtrFromString(port)

	handle, err := syscall.CreateFile(
		portU16,
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		0,
		nil,
		syscall.OPEN_EXISTING,
		0,
		0,
	)

	if err != nil {
		return fmt.Errorf("failed to open port %s: %w", port, err)
	}

	if err = syscall.CloseHandle(handle); err != nil {
		return fmt.Errorf("failed to close port %s: %w", port, err)
	}

	return nil
}

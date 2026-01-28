package proc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func (p *Processor) Deploy(pdxPath string) error {
	log.Println("deploying to Playdate device...")

	port, err := findPlaydatePort()
	if err != nil {
		return fmt.Errorf("failed to find Playdate device: %s", err)
	}
	log.Printf("found Playdate at: %s", port)

	sdkPath := os.Getenv("PLAYDATE_SDK_PATH")
	pdutil := filepath.Join(sdkPath, "bin", "pdutil")

	log.Println("switching Playdate to Data Disk mode...")
	if err = p.execCmd(pdutil, []string{port, "datadisk"}); err != nil {
		return fmt.Errorf("failed to switch to datadisk mode: %s", err)
	}

	log.Println("waiting for Playdate to mount...")
	mountPath, err := waitForMount(10 * time.Second)
	if err != nil {
		return fmt.Errorf("failed to detect mounted Playdate: %s", err)
	}
	log.Printf("Playdate mounted at: %s", mountPath)

	gamesPath := filepath.Join(mountPath, "Games")
	pdxName := filepath.Base(pdxPath)
	destPath := filepath.Join(gamesPath, pdxName)

	log.Printf("copying %s to %s...", pdxPath, destPath)

	// Remove existing if present
	if _, err = os.Stat(destPath); err == nil {
		if err = os.RemoveAll(destPath); err != nil {
			return fmt.Errorf("failed to remove existing game: %s", err)
		}
	}

	if err = p.execCmd("cp", []string{"-r", pdxPath, destPath}); err != nil {
		return fmt.Errorf("failed to copy game: %s", err)
	}
	log.Println("game has been copied successfully!")

	log.Println("ejecting Playdate...")
	if err = p.ejectPlaydate(mountPath); err != nil {
		return fmt.Errorf("failed to eject Playdate: %s", err)
	}

	log.Println("waiting for Playdate to reconnect...")
	time.Sleep(3 * time.Second)

	port, err = waitForPlaydatePort(30 * time.Second)
	if err != nil {
		return fmt.Errorf("failed to find Playdate after eject: %s", err)
	}
	log.Printf("Playdate reconnected at: %s", port)

	gamePath := fmt.Sprintf("Games/%s", pdxName)
	log.Printf("launching game: %s...", gamePath)

	if err = p.execCmd(pdutil, []string{port, "run", gamePath}); err != nil {
		return fmt.Errorf("failed to run game: %s", err)
	}

	log.Println("game has been successfully deployed and run")
	return nil
}

func findPlaydatePort() (string, error) {
	var patterns []string

	if runtime.GOOS == "darwin" {
		patterns = []string{
			"/dev/cu.usbmodemPD*",
			"/dev/cu.usbmodem*",
		}
	} else {
		// Linux
		patterns = []string{
			"/dev/ttyACM*",
			"/dev/ttyUSB*",
		}
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		if len(matches) > 0 {
			// Prefer ones with "PD" in name
			for _, m := range matches {
				if strings.Contains(m, "PD") {
					return m, nil
				}
			}
			return matches[0], nil
		}
	}

	return "", fmt.Errorf("no Playdate device found (is it connected and unlocked?)")
}

func waitForPlaydatePort(timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		port, err := findPlaydatePort()
		if err == nil {
			return port, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return "", fmt.Errorf("timeout waiting for Playdate device")
}

func waitForMount(timeout time.Duration) (string, error) {
	var possiblePaths []string

	if runtime.GOOS == "darwin" {
		possiblePaths = []string{
			"/Volumes/PLAYDATE",
		}
	} else {
		// Linux - check common mount points
		user := os.Getenv("USER")
		possiblePaths = []string{
			fmt.Sprintf("/media/%s/PLAYDATE", user),
			fmt.Sprintf("/run/media/%s/PLAYDATE", user),
			"/mnt/PLAYDATE",
		}
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		for _, path := range possiblePaths {
			if stat, err := os.Stat(path); err == nil && stat.IsDir() {
				gamesPath := filepath.Join(path, "Games")
				if _, err = os.Stat(gamesPath); err == nil {
					return path, nil
				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	return "", fmt.Errorf("timeout waiting for Playdate to mount")
}

func (p *Processor) ejectPlaydate(mountPath string) error {
	cmdName := ""
	var args []string

	if runtime.GOOS == "darwin" {
		cmdName = "diskutil"
		args = []string{"eject", mountPath}
	} else {
		cmdName = "umount"
		args = []string{mountPath}
	}

	return p.execCmd(cmdName, args)
}

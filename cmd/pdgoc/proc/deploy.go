package proc

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/playdate-go/pdgo/cmd/pdgoc/utils"
)

func (p *Processor) Deploy(pdxPath string) error {
	log.Println("deploying to Playdate device...")
	log.Println("ensure your device is unlocked and connected")

	if _, err := waitForPlaydatePort(10 * time.Second); err != nil {
		return fmt.Errorf("failed to find Playdate: %s", err)
	}

	pdutil, err := utils.GetPdutilPath()
	if err != nil {
		return fmt.Errorf("failed to find pdutil: %w", err)
	}

	log.Println("installing game to Playdate...")
	if err = p.execCmd(pdutil, []string{"install", pdxPath}); err != nil {
		return fmt.Errorf("failed to install game: %s", err)
	}

	pdxName := filepath.Base(pdxPath)
	installedPdxPath := fmt.Sprintf("Games/%s", pdxName)

	log.Println("waiting for Playdate...")
	port, err := waitForPlaydatePort(10 * time.Second)
	if err != nil {
		return fmt.Errorf("failed to find Playdate after eject: %s", err)
	}
	log.Printf("Playdate reconnected at: %s", port)

	log.Printf("launching game: %s...", installedPdxPath)
	if err = p.execCmd(pdutil, []string{"run", installedPdxPath}); err != nil {
		return fmt.Errorf("failed to run game: %s", err)
	}

	log.Println("game has been successfully deployed and run")
	return nil
}

func waitForPlaydatePort(timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		port, err := utils.FindPlaydatePort()
		if err == nil {
			return port, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return "", fmt.Errorf("timeout waiting for Playdate device")
}

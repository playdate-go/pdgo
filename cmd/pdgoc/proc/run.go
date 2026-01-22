package proc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (p *Processor) processRun() error {
	if err := p.processSim(); err != nil {
		return err
	}
	simulatorBinPath := filepath.Join(os.Getenv("PLAYDATE_SDK_PATH"), "bin", "Playdate Simulator.app")

	log.Printf("running '%s' in Playdate Simulator...", p.cfg.System.OutputPath)
	if err := p.execCmd("open", []string{simulatorBinPath, p.cfg.System.OutputPath}); err != nil {
		return fmt.Errorf("failed to run %s in Playdate Simulator: %s", p.cfg.System.OutputPath, err)
	}
	log.Printf("successfully ran '%s' in Playdate Simulator", p.cfg.System.OutputPath)
	return nil
}

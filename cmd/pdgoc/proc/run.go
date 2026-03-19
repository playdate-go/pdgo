package proc

import (
	"fmt"
	"log"

	"github.com/playdate-go/pdgo/cmd/pdgoc/utils"
)

func (p *Processor) processRun() error {
	simulatorPath, err := utils.GetSimulatorPath()
	if err != nil {
		return err
	}

	log.Printf("running '%s' in Playdate Simulator...", p.cfg.System.OutputPath)
	if err := p.execCmd(simulatorPath, []string{p.cfg.System.OutputPath}); err != nil {
		return fmt.Errorf("failed to run %s in Playdate Simulator: %s", p.cfg.System.OutputPath, err)
	}
	log.Printf("successfully ran '%s' in Playdate Simulator", p.cfg.System.OutputPath)
	return nil
}

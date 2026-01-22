package proc

import (
	"fmt"
	"github.com/playdate-go/pdgoc/config"
	"log"
	"os"
	"os/exec"
)

type Processor struct {
	cfg *config.Config
}

func NewProcessor(cfg *config.Config) *Processor {
	return &Processor{cfg: cfg}
}

func (p *Processor) Process() error {
	switch {
	case p.cfg.System.RunMode:
		log.Println("mode: run (build and run project in Playdate Simulator)")
		if err := p.processRun(); err != nil {
			return err
		}

	case p.cfg.System.SimMode && p.cfg.System.DeviceMode:
		log.Println("mode: simulator and device (build project for Playdate Simulator and real Playdate Console)")
		if err := p.processSim(); err != nil {
			return err
		}
		if err := p.processDevice(); err != nil {
			return err
		}

	case p.cfg.System.SimMode:
		log.Println("mode: simulator (build project for Playdate Simulator)")
		if err := p.processSim(); err != nil {
			return err
		}

	case p.cfg.System.DeviceMode:
		log.Println("mode: device (build project for real Playdate Console)")
		if err := p.processDevice(); err != nil {
			return err
		}

	default:
		return fmt.Errorf("no mode specified")
	}

	return nil
}

func (p *Processor) execCmd(name string, args []string) error {
	log.Printf("runnign cmd: %s %v ...", name, args)

	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	log.Printf("successfully ran cmd: %s", name)
	return nil
}

package config

import (
	"encoding/json"
	"flag"
	"fmt"
)

type Meta struct {
	Name            string
	Author          string
	Desc            string
	BundleID        string
	Version         string
	BuildNumber     string
	ImagePath       string
	LaunchSoundPath string
	ContentWarn     string
	ContentWarn2    string
}

func (meta *Meta) String() string {
	b, _ := json.Marshal(meta)
	return fmt.Sprintf(string(b))
}

func (meta *Meta) Validate() error {
	if meta == nil {
		return fmt.Errorf("passed nil config.Meta struct")
	}
	if meta.Name == "" {
		return fmt.Errorf("'name' pdxinfo property is mandatory, please set it using '-name' flag")
	}
	if meta.Author == "" {
		return fmt.Errorf("'author' pdxinfo property is mandatory, please set it using '-author' flag")
	}
	if meta.Desc == "" {
		return fmt.Errorf("'description' pdxinfo property is mandatory, please set it using '-desc' flag")
	}
	if meta.BundleID == "" {
		return fmt.Errorf("'bundleID' pdxinfo property is mandatory, please set it using '--bundle-id' flag")
	}
	if meta.Version == "" {
		return fmt.Errorf("'version' pdxinfo property is mandatory, please set it using '-version' flag")
	}
	if meta.BuildNumber == "" {
		return fmt.Errorf("'buildNumber' pdxinfo property is mandatory, please set it using '-build-number' flag")
	}
	// imagePath is optional
	return nil
}

type System struct {
	SimMode    bool
	DeviceMode bool
	RunMode    bool
	DeployMode bool
	InputPath  string
	OutputPath string
}

func (system *System) String() string {
	b, _ := json.Marshal(system)
	return fmt.Sprintf(string(b))
}

func (system *System) Validate() error {
	if !system.SimMode && !system.DeviceMode && !system.RunMode {
		return fmt.Errorf("at least '-sim' or '-device' or '-run' must be defined")
	}
	if system.DeployMode && !system.DeviceMode {
		return fmt.Errorf("'-deploy' requires '-device' flag")
	}

	return nil
}

type Config struct {
	System *System
	Meta   *Meta
}

func (cfg *Config) String() string {
	b, _ := json.Marshal(cfg)
	return fmt.Sprintf(string(b))
}

func Load() (*Config, error) {
	cfg := &Config{
		System: &System{},
		Meta:   &Meta{},
	}

	flag.BoolVar(&cfg.System.SimMode, "sim", false, "build project for Playdate Simulator")
	flag.BoolVar(&cfg.System.DeviceMode, "device", false, "build project for real Playdate console")
	flag.BoolVar(&cfg.System.RunMode, "run", false, "build and run project in Playdate Simulator")
	flag.BoolVar(&cfg.System.DeployMode, "deploy", false, "deploy and run on connected Playdate device (requires -device)")

	flag.StringVar(&cfg.Meta.Name, "name", "", "set pdxinfo 'name' property")
	flag.StringVar(&cfg.Meta.Author, "author", "", "set pdxinfo 'author' property")
	flag.StringVar(&cfg.Meta.Desc, "desc", "", "set pdxinfo 'description' property")
	flag.StringVar(&cfg.Meta.BundleID, "bundle-id", "", "set pdxinfo 'bundleID' property")
	flag.StringVar(&cfg.Meta.Version, "version", "", "set pdxinfo 'version' property")
	flag.StringVar(&cfg.Meta.BuildNumber, "build-number", "", "set pdxinfo 'buildNumber' property")
	flag.StringVar(&cfg.Meta.ImagePath, "image-path", "", "set pdxinfo 'imagePath' property")
	flag.StringVar(&cfg.Meta.LaunchSoundPath, "launch-sound-path", "", "set pdxinfo 'launchSoundPath'")
	flag.StringVar(&cfg.Meta.ContentWarn, "content-warn", "", "set pdxinfo 'contentWarning'")
	flag.StringVar(&cfg.Meta.ContentWarn2, "content-warn2", "", "set pdxinfo 'contentWarning2'")

	flag.Parse()

	if err := cfg.Meta.Validate(); err != nil {
		return nil, err
	}
	if err := cfg.System.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

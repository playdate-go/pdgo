package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/playdate-go/pdgo/cmd/pdgoc/config"
	"github.com/playdate-go/pdgo/cmd/pdgoc/pdxinfo"
	"github.com/playdate-go/pdgo/cmd/pdgoc/proc"
	"github.com/playdate-go/pdgo/cmd/pdgoc/utils"
)

var Version string
var Commit string
var Date string

func main() {
	log.Println("PdGo (Playdate Go) build tool: https://github.com/playdate-go")
	log.Printf("version: %s", Version)
	log.Printf("commit: %s", Commit)
	log.Printf("date: %s", Date)

	if runtime.GOOS == "windows" {
		log.Print("windows support is experimental, features may not work as expected")
	}

	sdkPath, err := utils.GetPlayDateSDKPath()
	if err != nil {
		log.Fatal(err)
	}

	err = utils.CheckPlayDateSDKVersion(sdkPath)
	if err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %s", err)
	}

	sourcePath, err := checkSourcePath(cwd)
	if err != nil {
		log.Fatalf("failed to check 'Source' path: %s\n", err)
	}

	log.Println("loading config...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %s\n", err)
	}

	log.Printf("config has been successfully loaded! :\n{%s}", cfg)

	log.Printf("working directory: %s, project name: %s", cwd, cfg.Meta.Name)

	log.Println("creating pdxinfo file...")

	if err = pdxinfo.CreateFile(cfg, sourcePath); err != nil {
		log.Fatalf("failed to create pdxinfo file: %s\n", err)
	}

	log.Printf("pdxinfo file has been successfully created! : %s", filepath.Join(sourcePath, "pdxinfo"))

	startTime := time.Now()

	p := proc.NewProcessor(cfg)

	if err = p.Process(); err != nil {
		log.Fatalf("failed to process: %s\n", err)
	}

	pdxInfoPath := filepath.Join(cwd, "Source", "pdxinfo")

	if err = os.Remove(pdxInfoPath); err != nil {
		log.Printf("failed to remove tmp file: %s\n", err)
	}
	log.Printf("tmp file has been successfully removed! : %s", pdxInfoPath)

	log.Printf("project has been sucessfully executed in %s", time.Since(startTime))
}

func checkSourcePath(cwd string) (string, error) {
	sourcePath := filepath.Join(cwd, "Source")
	if stat, err := os.Stat(sourcePath); err == nil && stat.IsDir() {
		log.Printf("auto-found 'Source' dir: %s", sourcePath)
		return sourcePath, nil
	}
	return "", fmt.Errorf("'Source' dir not found: %s", sourcePath)
}

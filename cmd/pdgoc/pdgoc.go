package main

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/playdate-go/pdgoc/config"
	"github.com/playdate-go/pdgoc/pdxinfo"
	"github.com/playdate-go/pdgoc/proc"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
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
		log.Fatal("currently this OS is unsupported: windows, supported: Linux, MacOS")
	}

	if err := checkSDKPathAndAPI(); err != nil {
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

func checkSDKPathAndAPI() error {
	if os.Getenv("PLAYDATE_SDK_PATH") == "" {
		return errors.New("PLAYDATE_SDK_PATH environment variable not set")
	}

	sdkPath := os.Getenv("PLAYDATE_SDK_PATH")
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

	return fmt.Errorf("failed to check SDK path")
}

func checkSourcePath(cwd string) (string, error) {
	sourcePath := filepath.Join(cwd, "Source")
	if stat, err := os.Stat(sourcePath); err == nil && stat.IsDir() {
		log.Printf("auto-found 'Source' dir: %s", sourcePath)
		return sourcePath, nil
	}
	return "", fmt.Errorf("'Source' dir not found: %s", sourcePath)
}

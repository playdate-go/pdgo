package proc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func (p *Processor) processSim() error {
	log.Printf("building host dynamic shared library from Go source: Go = %s, Arch = %s, OS = %s ...", runtime.Version(), runtime.GOARCH, runtime.GOOS)

	libPath, err := p.execBuildGoHostLib()
	if err != nil {
		return fmt.Errorf("failed to build host dynamic shared library from Go source: %s", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %s\n", err)
	}

	log.Println("successfully built host dynamic shared library from Go source")

	sourcePath := filepath.Join(cwd, "Source")
	p.cfg.System.OutputPath = fmt.Sprintf("%s_sim.pdx", p.cfg.Meta.Name)

	log.Println("building Playdate Application...")

	if err = p.execCmd("pdc", []string{"-k", "-s", "-v", sourcePath, p.cfg.System.OutputPath}); err != nil {
		return err
	}

	log.Printf("'%s' content:", p.cfg.System.OutputPath)
	if err = p.execCmd("ls", []string{"-l", p.cfg.System.OutputPath}); err != nil {
		return err
	}

	log.Println("Playdate Application has been successfully built!")

	if err = os.Remove(libPath); err != nil {
		return fmt.Errorf("failed to remove tmp file: %s", libPath)
	}
	log.Printf("tmp file has been successfully removed! : %s", libPath)

	return nil
}

func (p *Processor) execBuildGoHostLib() (string, error) {
	libExt, err := p.resolveHostLibExt()
	if err != nil {
		return "", err
	}
	log.Printf("host OS: %s, host dynamic shared library extension: %s", runtime.GOOS, libExt)

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %s\n", err)
	}

	goSourcePath := filepath.Join(cwd, "Source")
	libPath := filepath.Join(goSourcePath, fmt.Sprintf("pdex.%s", libExt))

	// change to 'Source' directory to build from there
	if err = os.Chdir(goSourcePath); err != nil {
		return "", fmt.Errorf("failed to change directory to %s: %s", goSourcePath, err)
	}
	// change back after all
	defer func() {
		if err = os.Chdir(cwd); err != nil {
			log.Printf("failed to change directory to %s: %s", cwd, err)
		}
	}()

	mainCgoGoFilePath := filepath.Join(goSourcePath, "main_cgo.go")

	log.Printf("writing %s...", mainCgoGoFilePath)
	if err = os.WriteFile(mainCgoGoFilePath, []byte(rawMainCgoGo), 0644); err != nil {
		return "", fmt.Errorf("failed to write main_cgo.go: %s", err)
	}
	log.Printf("file has been seccessfully created: %s", mainCgoGoFilePath)

	// Build using "." since we're now in the Source directory
	if err = p.execCmd("go", []string{"build", "-ldflags", "-w -s", "-gcflags", "all=-l", "-trimpath", "-buildvcs=false", "-race=false", "-o", libPath, "-buildmode=c-shared", "."}); err != nil {
		return "", err
	}

	if err = p.execCmd("file", []string{libPath}); err != nil {
		return "", err
	}

	pdexHPath := filepath.Join(goSourcePath, "pdex.h")
	if err = os.Remove(filepath.Join(goSourcePath, "pdex.h")); err == nil {
		log.Printf("tmp file has been seccessfully removed: %s", pdexHPath)
	}

	if err = os.Remove(mainCgoGoFilePath); err == nil {
		log.Printf("tmp file has been seccessfully removed: %s", mainCgoGoFilePath)
	}

	return libPath, nil
}

func (p *Processor) resolveHostLibExt() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return "so", nil
	case "darwin":
		return "dylib", nil
	default:
		return "", fmt.Errorf("unsupported platform")
	}
}

package proc

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (p *Processor) processDevice() error {
	log.Println("building Playdate executable and linkable format (ELF) file...")
	if err := p.runBuildScript(); err != nil {
		return fmt.Errorf("failed to build Playdate executable and linkable format (ELF) file...: %s", err)
	}

	log.Println("Playdate executable and linkable format (ELF) file has been successfully built!")
	return nil
}

// findPdgoPath finds the path to the pdgo module using go list
func findPdgoPath(goSrcDir string) (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/playdate-go/pdgo")
	cmd.Dir = goSrcDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to find pdgo module: %s", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (p *Processor) runBuildScript() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %s", err)
	}

	gameDir := cwd
	buildDir := filepath.Join(cwd, "build")

	var goSrcDir string
	if _, err = os.Stat(filepath.Join(cwd, "Source", "go.mod")); err == nil {
		goSrcDir = filepath.Join(cwd, "Source")
	} else {
		return fmt.Errorf("go.mod and .go sources not found in %s/Source", cwd)
	}

	if err = os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %s", err)
	}

	// Find pdgo module path and copy pd_cgo.c
	pdgoPath, err := findPdgoPath(goSrcDir)
	if err != nil {
		return fmt.Errorf("failed to find pdgo module: %s", err)
	}
	log.Printf("found pdgo module at: %s", pdgoPath)

	// Read pd_cgo.c from pdgo module
	pdCgoSrc := filepath.Join(pdgoPath, "pd_cgo.c")
	pdCgoContent, err := os.ReadFile(pdCgoSrc)
	if err != nil {
		return fmt.Errorf("failed to read pd_cgo.c from pdgo: %s", err)
	}

	// Write pd_cgo.c to build directory (will be compiled with -DTARGET_PLAYDATE=1)
	pdRuntimeFilePath := filepath.Join(buildDir, "pd_runtime.c")
	log.Printf("copying %s to %s...", pdCgoSrc, pdRuntimeFilePath)
	if err = os.WriteFile(pdRuntimeFilePath, pdCgoContent, 0644); err != nil {
		return fmt.Errorf("failed to write pd_runtime.c: %s", err)
	}
	log.Printf("file has been successfully created: %s", pdRuntimeFilePath)

	// Write main_tinygo.go to Source directory
	mainTinyGoFilePath := filepath.Join(goSrcDir, "main_tinygo.go")
	log.Printf("writing %s...", mainTinyGoFilePath)
	if err = os.WriteFile(mainTinyGoFilePath, []byte(rawMainTinyGo), 0644); err != nil {
		return fmt.Errorf("failed to write main_tinygo.go: %s", err)
	}
	log.Printf("file has been successfully created: %s", mainTinyGoFilePath)

	// Run go mod tidy
	log.Println("running cmd 'go mod tidy'...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = goSrcDir
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run 'go mod tidy': %s\n%s", err, string(output))
	}
	log.Printf("successfully ran 'go mod tidy'")

	// Create temporary build script
	buildScriptFile, err := os.CreateTemp("", "device-build-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create temp build script file: %s", err)
	}
	defer func() {
		if err = os.Remove(buildScriptFile.Name()); err != nil {
			log.Printf("warning: failed to remove temp file: %s", err)
		}
	}()

	if _, err = buildScriptFile.WriteString(rawBuildScript); err != nil {
		return fmt.Errorf("failed to write temp file '%s': %s", buildScriptFile.Name(), err)
	}
	if err = buildScriptFile.Close(); err != nil {
		log.Printf("warning: failed to close temp file '%s': %s", buildScriptFile.Name(), err)
	}

	if err = os.Chmod(buildScriptFile.Name(), 0755); err != nil {
		return fmt.Errorf("warning: failed to chmod build script: %s", err)
	}

	// Run the build script
	cmd := exec.Command("bash", buildScriptFile.Name())
	cmd.Dir = goSrcDir
	cmd.Env = append(os.Environ(),
		"GAME_NAME="+p.cfg.Meta.Name,
		"GAME_DIR="+gameDir,
		"GO_SRC_DIR="+goSrcDir,
		"BUILD_DIR="+buildDir,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	buildErr := cmd.Run()

	// Cleanup main_tinygo.go
	if os.Remove(mainTinyGoFilePath) == nil {
		log.Printf("tmp file has been successfully removed! : %s", mainTinyGoFilePath)
	}

	// Optionally cleanup build directory on success
	if buildErr == nil {
		os.RemoveAll(buildDir)
		log.Printf("build directory has been cleaned up: %s", buildDir)
	}

	return buildErr
}

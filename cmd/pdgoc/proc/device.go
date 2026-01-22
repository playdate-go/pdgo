package proc

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func (p *Processor) processDevice() error {
	log.Println("building Playdate executable and linkable format (ELF) file...")
	if err := p.runBuildScript(); err != nil {
		return fmt.Errorf("failed to build Playdate executable and linkable format (ELF) file...: %s", err)
	}

	log.Println("Playdate executable and linkable format (ELF) file has been successfully built!")
	return nil
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

	pdRuntimeFilePath := filepath.Join(buildDir, "pd_runtime.c")

	log.Printf("writing %s...", pdRuntimeFilePath)

	if err = os.WriteFile(pdRuntimeFilePath, []byte(rawRuntimeC), 0644); err != nil {
		return fmt.Errorf("failed to write pd_runtime.c: %s", err)
	}
	log.Printf("file has been successfuly created: %s", pdRuntimeFilePath)

	bridgeTemplateFilePath := filepath.Join(goSrcDir, "bridge_template.go")

	log.Printf("writing %s...", bridgeTemplateFilePath)
	if err := os.WriteFile(bridgeTemplateFilePath, []byte(rawBridgeTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write bridge_template.go: %s", err)
	}
	log.Printf("file has been seccessfully created: %s", bridgeTemplateFilePath)

	log.Println("running cmd 'go mod tidy'...")

	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = goSrcDir
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run 'go mod tidy': %s\n%s", err, string(output))
	}
	log.Printf("successfully ran 'go mod tidy'")

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

	mainTinyGoFilePath := filepath.Join(goSrcDir, "main_tinygo.go")

	log.Printf("writing %s...", mainTinyGoFilePath)
	if err = os.WriteFile(mainTinyGoFilePath, []byte(rawMainTinyGo), 0644); err != nil {
		return fmt.Errorf("failed to write main_tinygo.go: %s", err)
	}
	log.Printf("file has been seccessfully created: %s", mainTinyGoFilePath)

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

	if err = os.Remove(bridgeTemplateFilePath); err == nil {
		log.Printf("tmp file has been successfully removed! : %s", bridgeTemplateFilePath)
	}

	if os.Remove(mainTinyGoFilePath) == nil {
		log.Printf("tmp file has been successfully removed! : %s", mainTinyGoFilePath)
	}

	return buildErr
}

package pdxinfo

import (
	"fmt"
	"github.com/playdate-go/pdgoc/config"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateFile(cfg *config.Config, path string) error {
	if cfg == nil {
		panic("passed nil config struct")
	}

	pdxInfoPath := filepath.Join(path, "pdxinfo")
	f, err := os.Create(pdxInfoPath)
	if err != nil {
		return fmt.Errorf("failed to create pdxinfo file %s: %w", pdxInfoPath, err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Printf("failed to close pdxinfo file %s: %s", pdxInfoPath, err)
			return
		}
	}()

	// Mandatory fields
	fields := []struct {
		key   string
		value string
	}{
		{"name", cfg.Meta.Name},
		{"author", cfg.Meta.Author},
		{"description", cfg.Meta.Desc},
		{"bundleID", cfg.Meta.BundleID},
		{"version", cfg.Meta.Version},
		{"buildNumber", cfg.Meta.BuildNumber},
	}

	for _, field := range fields {
		if _, err = fmt.Fprintf(f, "%s=%s\n", field.key, field.value); err != nil {
			return fmt.Errorf("failed to write %s to pdxinfo: %w", field.key, err)
		}
		log.Printf("pdxinfo: added property, name='%s', value=%s", field.key, field.value)
	}

	// Optional fields
	optional := []struct {
		key   string
		value string
	}{
		{"imagePath", cfg.Meta.ImagePath},
		{"launchSoundPath", cfg.Meta.LaunchSoundPath},
		{"contentWarning", cfg.Meta.ContentWarn},
		{"contentWarning2", cfg.Meta.ContentWarn2},
	}

	for _, field := range optional {
		if field.value != "" {
			if _, err = fmt.Fprintf(f, "%s=%s\n", field.key, field.value); err != nil {
				return fmt.Errorf("failed to write %s to pdxinfo: %w", field.key, err)
			}
			log.Printf("pdxinfo: added property, name='%s', value=%s", field.key, field.value)
		}
	}

	if _, err = f.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek to beginning of pdxinfo file: %w", err)
	}
	content := make([]byte, 1024*10)
	n, err := f.Read(content)
	if err != nil {
		return fmt.Errorf("failed to read content of pdxinfo file: %w", err)
	}
	log.Printf("pdxinfo content (%d bytes):\n%s", n, strings.TrimSpace(string(content)))

	return nil
}

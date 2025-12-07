package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureCommonIgnore ensures common.gitignore exists, creates it if not
func EnsureCommonIgnore(configDir string) (string, error) {
	path := filepath.Join(configDir, "common.gitignore")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Creating default common.gitignore: %s\n", path)
		if err := createDefaultIgnore(path); err != nil {
			return "", fmt.Errorf("failed to create common.gitignore: %w", err)
		}
	}
	return path, nil
}

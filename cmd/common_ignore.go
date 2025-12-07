package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

// デフォルトの無視ルール
//
//go:embed default.txt
var defaultContent []byte

func DefaultContent() []byte {
	return defaultContent
}

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

// createDefaultIgnore creates a common.gitignore file with default ignore rules
func createDefaultIgnore(path string) error {
	// ディレクトリを作成
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(defaultContent), 0644)
}

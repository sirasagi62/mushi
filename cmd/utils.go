package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	interactive bool
	noUpdate    bool
	print       bool
)

// getCacheDir returns the path to the cache directory
func getCacheDir() (string, error) {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome != "" {
		return filepath.Join(cacheHome, "mushi", "github-gitignore"), nil
	}

	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("HOME environment variable is not set")
	}

	return filepath.Join(home, ".cache", "mushi", "github-gitignore"), nil
}

// ResolveImports は、content 内の "#Import:template" 行を展開して、
// 対応するテンプレートの内容に置き換えます。
func ResolveImports(content []byte, cacheDir string) ([]byte, error) {
	var result []byte
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#Import:") {
			templateName := strings.TrimSpace(trimmed[8:])
			if templateName == "" {
				continue
			}

			templatePath := filepath.Join(cacheDir, templateName+".gitignore")
			imported, err := os.ReadFile(templatePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to import %s: %v\n", templateName, err)
				continue
			}

			result = append(result, imported...)
			result = append(result, '\n')
		} else {
			result = append(result, line...)
			result = append(result, '\n')
		}
	}

	return result, nil
}

// getConfigDir returns the path to the config directory
func getConfigDir() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome != "" {
		return filepath.Join(configHome, "mushi"), nil
	}

	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("HOME environment variable is not set")
	}

	return filepath.Join(home, ".config", "mushi"), nil
}

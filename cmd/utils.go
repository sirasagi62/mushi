package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	interactive bool
	noUpdate    bool
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

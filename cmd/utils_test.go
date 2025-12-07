package cmd

import (
	"path/filepath"
	"testing"
)

func TestGetCacheDir(t *testing.T) {
	home := "/home/user"
	cacheHome := "/cache/custom"

	tests := []struct {
		name           string
		home           string
		cacheHome      string
		expectedSuffix string
	}{
		{
			name:           "XDG_CACHE_HOME set",
			home:           home,
			cacheHome:      cacheHome,
			expectedSuffix: "custom/mushi/github-gitignore",
		},
		{
			name:           "XDG_CACHE_HOME not set",
			home:           home,
			cacheHome:      "",
			expectedSuffix: ".cache/mushi/github-gitignore",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 環境変数を設定
			if tt.cacheHome != "" {
				t.Setenv("XDG_CACHE_HOME", tt.cacheHome)
				t.Setenv("HOME", tt.home)
			} else {
				t.Setenv("XDG_CACHE_HOME", "")
				t.Setenv("HOME", tt.home)
			}

			got, err := getCacheDir()
			if err != nil {
				t.Fatalf("getCacheDir() error = %v", err)
			}

			if !filepath.IsAbs(got) {
				t.Errorf("getCacheDir() returned non-absolute path: %s", got)
			}

			if tt.cacheHome != "" {
				expected := filepath.Join(tt.cacheHome, "mushi", "github-gitignore")
				if got != expected {
					t.Errorf("getCacheDir() = %v, expected %v", got, expected)
				}
			} else {
				expected := filepath.Join(tt.home, ".cache", "mushi", "github-gitignore")
				if got != expected {
					t.Errorf("getCacheDir() = %v, expected %v", got, expected)
				}
			}
		})
	}

	t.Run("HOME not set", func(t *testing.T) {
		t.Setenv("HOME", "")
		t.Setenv("XDG_CACHE_HOME", "")

		_, err := getCacheDir()
		if err == nil {
			t.Error("getCacheDir() should return error when HOME is not set")
		}
	})
}

func TestGetConfigDir(t *testing.T) {
	home := "/home/user"
	configHome := "/config/custom"

	tests := []struct {
		name           string
		home           string
		configHome     string
		expectedSuffix string
	}{
		{
			name:           "XDG_CONFIG_HOME set",
			home:           home,
			configHome:     configHome,
			expectedSuffix: "custom/mushi",
		},
		{
			name:           "XDG_CONFIG_HOME not set",
			home:           home,
			configHome:     "",
			expectedSuffix: ".config/mushi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 環境変数を設定
			if tt.configHome != "" {
				t.Setenv("XDG_CONFIG_HOME", tt.configHome)
				t.Setenv("HOME", tt.home)
			} else {
				t.Setenv("XDG_CONFIG_HOME", "")
				t.Setenv("HOME", tt.home)
			}

			got, err := getConfigDir()
			if err != nil {
				t.Fatalf("getConfigDir() error = %v", err)
			}

			if !filepath.IsAbs(got) {
				t.Errorf("getConfigDir() returned non-absolute path: %s", got)
			}

			if tt.configHome != "" {
				expected := filepath.Join(tt.configHome, "mushi")
				if got != expected {
					t.Errorf("getConfigDir() = %v, expected %v", got, expected)
				}
			} else {
				expected := filepath.Join(tt.home, ".config", "mushi")
				if got != expected {
					t.Errorf("getConfigDir() = %v, expected %v", got, expected)
				}
			}
		})
	}

	t.Run("HOME not set", func(t *testing.T) {
		t.Setenv("HOME", "")
		t.Setenv("XDG_CONFIG_HOME", "")

		_, err := getConfigDir()
		if err == nil {
			t.Error("getConfigDir() should return error when HOME is not set")
		}
	})
}

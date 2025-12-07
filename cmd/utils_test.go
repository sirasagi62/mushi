package cmd

import (
	"io"
	"os"
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

// TestResolveImports は ResolveImports 関数の動作をテストします
// ResolveImportsはそれぞれのImportの末尾とファイル全体の末尾に改行を挿入する
func TestResolveImports(t *testing.T) {
	tmpDir := t.TempDir()
	cacheDir := filepath.Join(tmpDir, "cache")
	if err := os.MkdirAll(filepath.Join(cacheDir, "Global"), 0755); err != nil {
		t.Fatal(err)
	}

	// ダミーテンプレートを作成
	write := func(path, content string) {
		if err := os.WriteFile(filepath.Join(cacheDir, path), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("Go.gitignore", "bin/\n*.exe\n")
	write("Python.gitignore", "*.pyc\n__pycache__/\n")
	write("Global/Windows.gitignore", "Thumbs.db\n")

	tests := []struct {
		name     string
		input    string
		expected string
		warning  bool
	}{
		{
			name:     "simple import",
			input:    "#Import:Go\n*.log",
			expected: "bin/\n*.exe\n\n*.log\n",
		},
		{
			name:     "subdirectory import",
			input:    "#Import:Global/Windows\n",
			expected: "Thumbs.db\n\n\n",
		},
		{
			name:     "multiple imports",
			input:    "#Import:Go\n#Import:Python\n",
			expected: "bin/\n*.exe\n\n*.pyc\n__pycache__/\n\n\n",
		},
		{
			name:     "nonexistent template",
			input:    "#Import:Rust\n*.tmp\n",
			expected: "*.tmp\n\n",
			warning:  true,
		},
		{
			name:     "empty import",
			input:    "#Import:\n*.tmp\n",
			expected: "*.tmp\n\n",
		},
		{
			name:     "import in comment",
			input:    "# #Import:Go\n*.tmp\n",
			expected: "# #Import:Go\n*.tmp\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderrBuf []byte
			// os.Stderr を一時的に差し替え
			origStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			got, err := ResolveImports([]byte(tt.input), cacheDir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// 標準エラー出力をキャプチャ
			w.Close()
			stderrBuf, _ = io.ReadAll(r)
			os.Stderr = origStderr

			if string(got) != tt.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tt.expected, got)
			}

			if tt.warning && len(stderrBuf) == 0 {
				t.Error("expected warning message, but none was printed")
			}
		})
	}
}

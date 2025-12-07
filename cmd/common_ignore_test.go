package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureCommonIgnore(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "mushi")

	t.Run("creates common.gitignore when not exists", func(t *testing.T) {
		path, err := EnsureCommonIgnore(configDir)
		if err != nil {
			t.Fatalf("EnsureCommonIgnore failed: %v", err)
		}

		expectedPath := filepath.Join(configDir, "common.gitignore")
		if path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, path)
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("common.gitignore should exist at %s", path)
		}
	})

	t.Run("does not recreate when already exists", func(t *testing.T) {
		// 既にファイルが存在する状態を作る
		path, err := EnsureCommonIgnore(configDir)
		if err != nil {
			t.Fatalf("First call failed: %v", err)
		}

		// 2回目
		path2, err := EnsureCommonIgnore(configDir)
		if err != nil {
			t.Fatalf("Second call failed: %v", err)
		}

		if path != path2 {
			t.Errorf("Path should be consistent: %s vs %s", path, path2)
		}
	})
}

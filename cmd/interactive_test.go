package cmd

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestFindTemplates(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()

	// テスト用の .gitignore ファイルを作成
	testFiles := []string{
		"Go.gitignore",
		"Python.gitignore",
		"Node.gitignore",
		"nested/Rust.gitignore",
	}

	for _, file := range testFiles {
		path := filepath.Join(tmpDir, file)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// 関数をテスト
	templates, err := findTemplates(tmpDir)
	if err != nil {
		t.Fatalf("findTemplates() error: %v", err)
	}

	// 期待される結果
	expected := []string{"Go", "Python", "Node", "nested/Rust"}

	if len(templates) != len(expected) {
		t.Errorf("Expected %d templates, got %d", len(expected), len(templates))
	}

	for _, exp := range expected {
		found := slices.Contains(templates, exp)
		if !found {
			t.Errorf("Expected template %s not found", exp)
		}
	}
}

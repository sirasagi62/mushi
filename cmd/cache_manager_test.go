package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureCache(t *testing.T) {
	tmpDir := t.TempDir()
	cacheDir := filepath.Join(tmpDir, "cache")

	// テスト1: キャッシュディレクトリが存在しない場合、作成されるべき
	t.Run("creates cache when not exists", func(t *testing.T) {
		err := EnsureCache(cacheDir, true)
		if err != nil {
			t.Fatalf("EnsureCache failed: %v", err)
		}
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			t.Error("Cache directory should be created")
		}
	})

	// テスト2: キャッシュが存在する場合、skipUpdate=trueなら更新しない
	t.Run("skips update when skipUpdate is true", func(t *testing.T) {
		// キャッシュディレクトリを作成
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			t.Fatalf("Failed to create cache dir: %v", err)
		}

		// updateCacheは失敗してもEnsureCacheはエラーを返さない
		err := EnsureCache(cacheDir, true)
		if err != nil {
			t.Errorf("EnsureCache should not error when skipUpdate=true, got: %v", err)
		}
	})

	// テスト3: キャッシュが存在する場合、skipUpdate=falseなら更新を試みる（失敗しても続行）
	t.Run("attempts update when skipUpdate is false", func(t *testing.T) {
		// キャッシュディレクトリを作成
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			t.Fatalf("Failed to create cache dir: %v", err)
		}

		// updateCacheは内部で失敗しても、EnsureCacheはエラーを返さずに続行する
		err := EnsureCache(cacheDir, false)
		if err != nil {
			t.Errorf("EnsureCache should not return error even if update fails, got: %v", err)
		}
	})
}

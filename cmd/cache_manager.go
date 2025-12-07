package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// cloneCache clones the github/gitignore repository to the cache directory
func cloneCache(cacheDir string) error {
	// 親ディレクトリを作成
	if err := os.MkdirAll(filepath.Dir(cacheDir), 0755); err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/github/gitignore", cacheDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// EnsureCache ensures the cache directory exists and updates it if needed
// skipUpdateがtrueの場合は更新をスキップ
func EnsureCache(cacheDir string, skipUpdate bool) error {
	// キャッシュディレクトリが存在しない場合はクローン
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		fmt.Println("Cache not found. Cloning github/gitignore repository...")
		if err := cloneCache(cacheDir); err != nil {
			return fmt.Errorf("failed to clone cache: %w", err)
		}
		return nil
	}

	// キャッシュが存在する場合は更新を確認
	if skipUpdate {
		fmt.Println("Skipping cache update...")
	} else {
		fmt.Println("Updating cache...")
		if err := updateCache(cacheDir); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update cache: %v\nSkipping cache update.\n", err)
			// 更新失敗はエラーとせず続行
		}
	}

	return nil
}

package cmd

import (
	"fmt"
	"os"
)

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

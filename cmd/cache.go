package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage local cache of gitignore templates",
}

var cacheUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the local cache",
	Run: func(cmd *cobra.Command, args []string) {
		// キャッシュディレクトリが存在しない場合は、自動的に取得
		if _, err := os.Stat(CacheDir); os.IsNotExist(err) {
			fmt.Println("Cache not found. Cloning github/gitignore repository...")
			if err := cloneCache(CacheDir); err != nil {
				fmt.Fprintf(os.Stderr, "Error cloning cache: %v\n", err)
				os.Exit(1)
			}
		} else {
			// キャッシュディレクトリが存在する場合は、更新を確認
			fmt.Println("Updating cache...")
			if err := updateCache(CacheDir); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating cache: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

var cacheCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the local cache",
	Run: func(cmd *cobra.Command, args []string) {
		// キャッシュディレクトリが存在するか確認
		if _, err := os.Stat(CacheDir); os.IsNotExist(err) {
			fmt.Println("Cache directory does not exist")
			return
		}

		// キャッシュディレクトリを削除
		fmt.Printf("Removing cache directory: %s\n", CacheDir)
		if err := os.RemoveAll(CacheDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing cache directory: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Cache cleaned successfully")
	},
}

func init() {
	cacheCmd.AddCommand(cacheUpdateCmd)
	cacheCmd.AddCommand(cacheCleanCmd)
	RootCmd.AddCommand(cacheCmd)
}

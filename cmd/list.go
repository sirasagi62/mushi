package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available gitignore templates",
	Run: func(cmd *cobra.Command, args []string) {
		// キャッシュディレクトリのパスを取得
		cacheDir, err := getCacheDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting cache directory: %v\n", err)
			os.Exit(1)
		}

		// キャッシュディレクトリが存在しない場合はクローン
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			fmt.Println("Cache not found. Cloning github/gitignore repository...")
			if err := cloneCache(cacheDir); err != nil {
				fmt.Fprintf(os.Stderr, "failed to clone cache: %v", err)
				os.Exit(1)
			}
		}

		// キャッシュディレクトリ内のすべての .gitignore ファイルを再帰的に検索
		var templates []string
		err = filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".gitignore") {
				// キャッシュディレクトリからの相対パスを取得
				relPath, err := filepath.Rel(cacheDir, path)
				if err != nil {
					return err
				}
				// .gitignore 拡張子を除去
				templateName := strings.TrimSuffix(relPath, ".gitignore")
				templates = append(templates, templateName)
			}
			return nil
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading cache directory: %v\n", err)
			os.Exit(1)
		}

		// テンプレートをソートして表示
		if len(templates) == 0 {
			fmt.Println("No templates found in cache")
			return
		}

		fmt.Printf("Available gitignore templates (%d):\n", len(templates))
		for _, template := range templates {
			fmt.Printf("  %s\n", template)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}

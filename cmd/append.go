package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var appendCmd = &cobra.Command{
	Use:   "append [template]",
	Short: "Append template to existing .gitignore",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// 既存の .gitignore が存在するか確認
		if _, err := os.Stat(".gitignore"); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: .gitignore does not exist in current directory\n")
			os.Exit(1)
		}

		// キャッシュディレクトリのパスを解決
		cacheDir, err := getCacheDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting cache directory: %v\n", err)
			os.Exit(1)
		}

		// 設定ディレクトリのパスを解決
		configDir, err := getConfigDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting config directory: %v\n", err)
			os.Exit(1)
		}

		var template string
		if interactive {
			// インタラクティブモード
			template, err = runInteractiveSelector(cacheDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error in interactive mode: %v\n", err)
				os.Exit(1)
			}
			if template == "" {
				fmt.Println("No template selected")
				return
			}
		} else {
			// 非インタラクティブモード
			if len(args) < 1 {
				fmt.Fprintln(os.Stderr, "Error: template name is required")
				os.Exit(1)
			}
			template = args[0]
		}

		// キャッシュディレクトリが存在しない場合は、自動的に取得
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			fmt.Println("Cache not found. Cloning github/gitignore repository...")
			if err := cloneCache(cacheDir); err != nil {
				fmt.Fprintf(os.Stderr, "Error cloning cache: %v\n", err)
				os.Exit(1)
			}
		} else {
			// キャッシュ更新の設定を決定
			// コマンドラインフラグが設定されている場合は、それ優先
			skipUpdate := noUpdate
			if !noUpdate {
				// コマンドラインフラグが未設定の場合は、設定ファイルの値を使用
				skipUpdate = config.NoUpdate
			}

			// キャッシュディレクトリが存在する場合は、更新を確認
			if skipUpdate {
				fmt.Println("Skipping cache update...")
			} else {
				fmt.Println("Updating cache...")
				if err := updateCache(cacheDir); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to update cache: %v\nSkipping cache update.", err)
				}
			}
		}

		// テンプレートファイルのパスを構築
		templateFile := template + ".gitignore"
		templatePath := filepath.Join(cacheDir, templateFile)

		// テンプレートファイルの内容を読み込む
		templateContent, err := os.ReadFile(templatePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading template file %s: %v\n", templateFile, err)
			os.Exit(1)
		}

		// 既存の .gitignore を読み込む
		existingContent, err := os.ReadFile(".gitignore")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading existing .gitignore: %v\n", err)
			os.Exit(1)
		}

		// common.gitignore を連結するかどうか
		var finalContent []byte
		if !noCommon {
			// 共通無視ファイルのパスを取得
			commonIgnorePath := filepath.Join(configDir, "common.gitignore")

			// 共通無視ファイルの内容を読み込む
			commonContent, err := os.ReadFile(commonIgnorePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading common.gitignore file: %v\n", err)
				os.Exit(1)
			}

			if len(commonContent) > 0 {
				finalContent = append(finalContent, commonContent...)
				finalContent = append(finalContent, '\n')
			}
		}

		// 既存の内容に追記
		finalContent = append(existingContent, '\n')
		finalContent = append(finalContent, templateContent...)

		// 結果を .gitignore に出力
		if err := os.WriteFile(".gitignore", finalContent, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to .gitignore: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✨️ Successfully appended %s to .gitignore\n", template)
	},
}

var (
	noCommon bool
)

func init() {
	appendCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactively select a template")
	appendCmd.Flags().BoolVar(&noUpdate, "no-update", false, "Skip updating the local cache")
	appendCmd.Flags().BoolVar(&noCommon, "no-common", false, "Do not include common.gitignore patterns")
	RootCmd.AddCommand(appendCmd)
}

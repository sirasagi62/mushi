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

		// キャッシュの存在確認と更新
		skipUpdate := noUpdate || config.NoUpdate
		if err := EnsureCache(cacheDir, skipUpdate); err != nil {
			fmt.Fprintf(os.Stderr, "Error managing cache: %v\n", err)
			os.Exit(1)
		}

		// テンプレートファイルのパスを構築
		templatePath := filepath.Join(cacheDir, template+".gitignore")

		// テンプレートファイルの内容を読み込む
		templateContent, err := os.ReadFile(templatePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading template file %s: %v\n", templatePath, err)
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
			// 共通無視ファイルの存在確認と作成
			commonIgnorePath, err := EnsureCommonIgnore(configDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error managing common.gitignore: %v\n", err)
				os.Exit(1)
			}

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

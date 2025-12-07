package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [template]",
	Short: "Generate .gitignore from a template",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
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

		// 共通無視ファイルの存在確認と作成
		commonIgnorePath, err := EnsureCommonIgnore(configDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error managing common.gitignore: %v\n", err)
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

		// 共通無視ファイルの内容を読み込む
		commonContent, err := os.ReadFile(commonIgnorePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading common.gitignore file: %v\n", err)
			os.Exit(1)
		}

		// インポートを解決
		var resolvedCommon []byte
		if len(commonContent) > 0 {
			resolvedCommon, err = ResolveImports(commonContent, cacheDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error resolving imports in common.gitignore: %v\n", err)
				os.Exit(1)
			}
		}

		// 両方の内容を結合
		var finalContent []byte
		if len(resolvedCommon) > 0 {
			finalContent = append(finalContent, resolvedCommon...)
			finalContent = append(finalContent, '\n')
		}
		finalContent = append(finalContent, templateContent...)

		// --print が指定されたら標準出力に表示
		if print {
			os.Stdout.Write(finalContent)
			return
		}
		// 出力ファイルのパスを設定
		outputPath := ".gitignore"

		// 既に .gitignore が存在するか確認
		if _, err := os.Stat(outputPath); err == nil {
			if !force {
				fmt.Fprintf(os.Stderr, "Error: %s already exists. Use -f or --force to overwrite.\n", outputPath)
				os.Exit(1)
			}
			fmt.Printf("Overwriting existing %s\n", outputPath)
		} else {
			fmt.Printf("Generating %s\n", outputPath)
		}

		// 結果を .gitignore に出力
		if err := os.WriteFile(outputPath, finalContent, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to %s: %v\n", outputPath, err)
			os.Exit(1)
		}

		fmt.Printf("✨️ Successfully generated %s\n", outputPath)
	},
}

var (
	force bool
)

func init() {
	createCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactively select a template")
	createCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing .gitignore file")
	createCmd.Flags().BoolVar(&noUpdate, "no-update", false, "Skip updating the local cache")
	createCmd.Flags().BoolVar(&print, "print", false, "Print the result to stdout instead of writing to .gitignore")
	RootCmd.AddCommand(createCmd)
}

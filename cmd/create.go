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

		// キャッシュディレクトリが存在しない場合は、自動的に取得
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			fmt.Println("Cache not found. Cloning github/gitignore repository...")
			if err := cloneCache(cacheDir); err != nil {
				fmt.Fprintf(os.Stderr, "Error cloning cache: %v\n", err)
				os.Exit(1)
			}
		} else {
			// キャッシュディレクトリが存在する場合は、更新を確認
			fmt.Println("Updating cache...")
			if err := updateCache(cacheDir); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating cache: %v\n", err)
				os.Exit(1)
			}
		}

		// 共通無視ファイルのパスを取得
		commonIgnorePath := filepath.Join(configDir, "common.ignore")

		// 共通無視ファイルが存在しない場合は、デフォルトの無視ルールで作成
		if _, err := os.Stat(commonIgnorePath); os.IsNotExist(err) {
			fmt.Printf("Creating common ignore file: %s\n", commonIgnorePath)
			if err := createDefaultIgnore(commonIgnorePath); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating common ignore file: %v\n", err)
				os.Exit(1)
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

		// 共通無視ファイルの内容を読み込む
		commonContent, err := os.ReadFile(commonIgnorePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading common ignore file: %v\n", err)
			os.Exit(1)
		}

		// 両方の内容を結合
		var finalContent []byte
		if len(commonContent) > 0 {
			finalContent = append(finalContent, commonContent...)
			finalContent = append(finalContent, '\n')
		}
		finalContent = append(finalContent, templateContent...)

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

		fmt.Printf("Successfully generated %s\n", outputPath)
	},
}

var (
	interactive bool
	force       bool
)

func init() {
	createCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactively select a template")
	createCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing .gitignore file")
	RootCmd.AddCommand(createCmd)
}

// getCacheDir returns the path to the cache directory
func getCacheDir() (string, error) {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome != "" {
		return filepath.Join(cacheHome, "mushi", "github-gitignore"), nil
	}

	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("HOME environment variable is not set")
	}

	return filepath.Join(home, ".cache", "mushi", "github-gitignore"), nil
}

// getConfigDir returns the path to the config directory
func getConfigDir() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome != "" {
		return filepath.Join(configHome, "mushi"), nil
	}

	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("HOME environment variable is not set")
	}

	return filepath.Join(home, ".config", "mushi"), nil
}

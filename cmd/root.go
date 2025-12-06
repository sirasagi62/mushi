package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	// ConfigDir は ~/.config/mushi を指します
	ConfigDir string
	// CacheDir は ~/.cache/mushi/github-gitignore を指します
	CacheDir string
	// CommonIgnorePath は共通無視ファイルのパスです
	CommonIgnorePath string
)

var RootCmd = &cobra.Command{
	Use:   "mushi",
	Short: "mushi is a gitignore template generator",
}

// デフォルトの無視ルール
//
//go:embed default.txt
var defaultContent []byte

func init() {
	// 共通フラグやサブコマンドの初期化
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// 設定ディレクトリのパスを解決
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home := os.Getenv("HOME")
		if home == "" {
			fmt.Fprintf(os.Stderr, "Error: HOME environment variable is not set\n")
			os.Exit(1)
		}
		configHome = filepath.Join(home, ".config")
	}
	ConfigDir = filepath.Join(configHome, "mushi")

	// キャッシュディレクトリのパスを解決
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home := os.Getenv("HOME")
		if home == "" {
			fmt.Fprintf(os.Stderr, "Error: HOME environment variable is not set\n")
			os.Exit(1)
		}
		cacheHome = filepath.Join(home, ".cache")
	}
	CacheDir = filepath.Join(cacheHome, "mushi", "github-gitignore")

	// 必要なディレクトリの作成
	dirs := []string{ConfigDir, filepath.Dir(CacheDir)}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	// 共通無視ファイルのパスを設定
	CommonIgnorePath = filepath.Join(ConfigDir, "common.ignore")

	// 共通無視ファイルが存在しない場合はデフォルトの無視ルールで作成
	if _, err := os.Stat(CommonIgnorePath); os.IsNotExist(err) {
		if err := createDefaultIgnore(CommonIgnorePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating common ignore file %s: %v\n", CommonIgnorePath, err)
			os.Exit(1)
		}
	}
}

// createDefaultIgnore creates a common.ignore file with default ignore rules
func createDefaultIgnore(path string) error {
	// ディレクトリを作成
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(defaultContent), 0644)
}

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

// updateCache updates the local cache with git pull
func updateCache(cacheDir string) error {
	cmd := exec.Command("git", "-C", cacheDir, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// ConfigDir は ~/.config/mushi を指します
	ConfigDir string
	// CacheDir は ~/.cache/mushi/github-gitignore を指します
	CacheDir string
	// CommonIgnorePath は共通無視ファイルのパスです
	CommonIgnorePath string
	// Version は mushi のバージョンです
	Version string = "v0.2.2"
)

// Config は mushi の設定を保持する構造体です
type Config struct {
	NoUpdate bool `mapstructure:"no_update"`
}

var config Config

var RootCmd = &cobra.Command{
	Use:     "mushi",
	Short:   "mushi is a gitignore template generator",
	Version: Version,
}

// デフォルトの無視ルール
//
//go:embed default.txt
var defaultContent []byte

func init() {
	// 共通フラグやサブコマンドの初期化
	cobra.OnInitialize(initConfig, initViper)
}

func initViper() {
	// Viperの設定
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(ConfigDir)

	// デフォルト値の設定
	viper.SetDefault("no_update", false)

	// 設定ファイルの読み込み
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// 設定ファイルが見つからない以外のエラーは警告
			fmt.Fprintf(os.Stderr, "Warning: Error reading config file: %v\n", err)
		}
		// 設定ファイルが見つからない場合は新規作成
		if err := createConfigFile(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Error creating config file: %v\n", err)
		}
	}

	// 設定を構造体にバインド
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config: %v\n", err)
		os.Exit(1)
	}
}

// createConfigFile creates a default config file
func createConfigFile() error {
	configContent := `# mushi configuration file

# Whether to skip updating the local cache
# no_update = false
`

	return os.WriteFile(filepath.Join(ConfigDir, "config.toml"), []byte(configContent), 0644)
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
	CommonIgnorePath = filepath.Join(ConfigDir, "common.gitignore")

	// 共通無視ファイルが存在しない場合はデフォルトの無視ルールで作成
	if _, err := os.Stat(CommonIgnorePath); os.IsNotExist(err) {
		if err := createDefaultIgnore(CommonIgnorePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating common.gitignore file %s: %v\n", CommonIgnorePath, err)
			os.Exit(1)
		}
	}
}

// createDefaultIgnore creates a common.gitignore file with default ignore rules
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

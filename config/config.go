package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config はボットの設定値を保持する構造体
type Config struct {
	DiscordToken string
	GithubToken  string
	DBPath       string
}

// Load は .env ファイルと環境変数から設定を読み込む
func Load() (*Config, error) {
	// .env ファイルがあれば読み込む（なくてもエラーにしない）
	_ = godotenv.Load()

	discordToken := os.Getenv("DISCORD_TOKEN")
	if discordToken == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN is required")
	}

	githubToken := os.Getenv("GITHUB_TOKEN")

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "kosobot.db"
	}

	return &Config{
		DiscordToken: discordToken,
		GithubToken:  githubToken,
		DBPath:       dbPath,
	}, nil
}

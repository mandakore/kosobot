package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/atashiro/kosobot/bot"
	"github.com/atashiro/kosobot/config"
	"github.com/atashiro/kosobot/db"
)

func main() {
	// 設定を読み込む
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// データベースを初期化
	database, err := db.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()
	log.Println("Database initialized successfully.")

	// ボットを作成・起動
	b, err := bot.New(cfg.DiscordToken, database)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if err := b.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
	defer b.Stop()

	// シグナルを待機してグレースフルシャットダウン
	log.Println("Press Ctrl+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	log.Println("Received shutdown signal.")
}

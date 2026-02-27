package bot

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Bot は Discord ボットのインスタンスを保持する構造体
type Bot struct {
	Session *discordgo.Session
	DB      *sql.DB
}

// New は新しい Bot インスタンスを作成する
func New(token string, db *sql.DB) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	b := &Bot{
		Session: session,
		DB:      db,
	}

	// Ready イベントハンドラを登録
	session.AddHandler(b.onReady)

	// 必要な Intent を設定
	session.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions

	return b, nil
}

// Start は Discord ボットを起動する
func (b *Bot) Start() error {
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord session: %w", err)
	}
	log.Println("Bot is now running.")
	return nil
}

// Stop は Discord ボットを停止する
func (b *Bot) Stop() error {
	log.Println("Shutting down bot...")
	return b.Session.Close()
}

// onReady は Bot が Discord に接続完了したときに呼ばれるハンドラ
func (b *Bot) onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Bot is ready! Logged in as: %s#%s", event.User.Username, event.User.Discriminator)
}

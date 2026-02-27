package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// InitDB は SQLite データベースを初期化し、必要なテーブルを作成する
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 接続確認
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// WAL モードを有効化（パフォーマンス向上）
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	// テーブル作成
	if err := createTables(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// createTables は必要なテーブルを作成する
func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		discord_id                TEXT PRIMARY KEY,
		github_id                 TEXT,
		current_xp                INTEGER NOT NULL DEFAULT 0,
		current_level             INTEGER NOT NULL DEFAULT 0,
		last_github_contributions INTEGER NOT NULL DEFAULT 0,
		updated_at                DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	return nil
}

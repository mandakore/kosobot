package db

import (
	"testing"
	"time"
)

func TestInitDB(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	// DB が nil でないことを確認
	if db == nil {
		t.Fatal("InitDB returned nil db")
	}
}

func TestUsersTableSchema(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	// テーブル情報を取得
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		t.Fatalf("Failed to query table info: %v", err)
	}
	defer rows.Close()

	expectedColumns := map[string]string{
		"discord_id":                "TEXT",
		"github_id":                 "TEXT",
		"current_xp":                "INTEGER",
		"current_level":             "INTEGER",
		"last_github_contributions": "INTEGER",
		"updated_at":                "DATETIME",
	}

	foundColumns := make(map[string]string)
	for rows.Next() {
		var cid int
		var name, colType string
		var notNull int
		var dfltValue *string
		var pk int
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		foundColumns[name] = colType
	}

	// 期待するカラムがすべて存在するか確認
	for colName, expectedType := range expectedColumns {
		actualType, ok := foundColumns[colName]
		if !ok {
			t.Errorf("Column %q not found in users table", colName)
			continue
		}
		if actualType != expectedType {
			t.Errorf("Column %q: expected type %q, got %q", colName, expectedType, actualType)
		}
	}

	// discord_id が PRIMARY KEY であることを確認
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='discord_id' AND pk=1").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to check PK: %v", err)
	}
	if count != 1 {
		t.Error("discord_id should be PRIMARY KEY")
	}
}

func TestInsertAndQueryUser(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	// ユーザーを挿入
	_, err = db.Exec(
		"INSERT INTO users (discord_id, github_id, current_xp, current_level, last_github_contributions, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		"123456789", "testuser", 100, 1, 50, time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// ユーザーを取得
	var discordID, githubID string
	var xp, level, lastContrib int
	var updatedAt string
	err = db.QueryRow("SELECT discord_id, github_id, current_xp, current_level, last_github_contributions, updated_at FROM users WHERE discord_id = ?", "123456789").
		Scan(&discordID, &githubID, &xp, &level, &lastContrib, &updatedAt)
	if err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}

	if discordID != "123456789" {
		t.Errorf("Expected discord_id '123456789', got %q", discordID)
	}
	if githubID != "testuser" {
		t.Errorf("Expected github_id 'testuser', got %q", githubID)
	}
	if xp != 100 {
		t.Errorf("Expected current_xp 100, got %d", xp)
	}
	if level != 1 {
		t.Errorf("Expected current_level 1, got %d", level)
	}
	if lastContrib != 50 {
		t.Errorf("Expected last_github_contributions 50, got %d", lastContrib)
	}
}

func TestInsertDuplicateUser(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	// 同じ discord_id で2回挿入 → PRIMARY KEY 違反
	_, err = db.Exec("INSERT INTO users (discord_id) VALUES (?)", "123456789")
	if err != nil {
		t.Fatalf("First insert failed: %v", err)
	}

	_, err = db.Exec("INSERT INTO users (discord_id) VALUES (?)", "123456789")
	if err == nil {
		t.Error("Expected error on duplicate discord_id insert, but got nil")
	}
}

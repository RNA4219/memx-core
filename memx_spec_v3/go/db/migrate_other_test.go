package db

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestMigrateJournal(t *testing.T) {
	// 一時ディレクトリを作成
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "journal.db")

	db, err := openDB("file:" + dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// マイグレーション実行
	if err := migrateJournal(db); err != nil {
		t.Fatalf("migrateJournal failed: %v", err)
	}

	// user_version が設定されているか確認
	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		t.Fatalf("failed to check user_version: %v", err)
	}
	if version != 1 {
		t.Errorf("expected user_version=1, got: %d", version)
	}

	// notes テーブルが存在するか確認（sqlite_master で確認）
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='notes';").Scan(&tableName)
	if err != nil {
		t.Errorf("notes table should exist: %v", err)
	}
	if tableName != "notes" {
		t.Errorf("expected table name 'notes', got: %s", tableName)
	}

	// tags テーブルが存在するか確認
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='tags';").Scan(&tableName)
	if err != nil {
		t.Errorf("tags table should exist: %v", err)
	}
	if tableName != "tags" {
		t.Errorf("expected table name 'tags', got: %s", tableName)
	}
}

func TestMigrateJournal_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "journal.db")

	db, err := openDB("file:" + dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// 1回目のマイグレーション
	if err := migrateJournal(db); err != nil {
		t.Fatalf("first migrateJournal failed: %v", err)
	}

	// 2回目のマイグレーション（再実行安全性の確認）
	if err := migrateJournal(db); err != nil {
		t.Fatalf("second migrateJournal failed: %v", err)
	}

	// user_version が 1 のままであることを確認
	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		t.Fatalf("failed to check user_version: %v", err)
	}
	if version != 1 {
		t.Errorf("expected user_version=1, got: %d", version)
	}
}

func TestMigrateKnowledge(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "knowledge.db")

	db, err := openDB("file:" + dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := migrateKnowledge(db); err != nil {
		t.Fatalf("migrateKnowledge failed: %v", err)
	}

	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		t.Fatalf("failed to check user_version: %v", err)
	}
	if version != 1 {
		t.Errorf("expected user_version=1, got: %d", version)
	}
}

func TestMigrateKnowledge_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "knowledge.db")

	db, err := openDB("file:" + dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := migrateKnowledge(db); err != nil {
		t.Fatalf("first migrateKnowledge failed: %v", err)
	}

	if err := migrateKnowledge(db); err != nil {
		t.Fatalf("second migrateKnowledge failed: %v", err)
	}

	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		t.Fatalf("failed to check user_version: %v", err)
	}
	if version != 1 {
		t.Errorf("expected user_version=1, got: %d", version)
	}
}

func TestMigrateArchive(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "archive.db")

	db, err := openDB("file:" + dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := migrateArchive(db); err != nil {
		t.Fatalf("migrateArchive failed: %v", err)
	}

	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		t.Fatalf("failed to check user_version: %v", err)
	}
	if version != 1 {
		t.Errorf("expected user_version=1, got: %d", version)
	}

	// archive には FTS がないことを確認
	// notes_fts テーブルは存在しないはず
	var exists int
	err = db.QueryRow("SELECT 1 FROM notes_fts LIMIT 1;").Scan(&exists)
	if err == nil {
		t.Error("archive should not have notes_fts table")
	}
}

func TestMigrateArchive_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "archive.db")

	db, err := openDB("file:" + dbPath)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := migrateArchive(db); err != nil {
		t.Fatalf("first migrateArchive failed: %v", err)
	}

	if err := migrateArchive(db); err != nil {
		t.Fatalf("second migrateArchive failed: %v", err)
	}

	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		t.Fatalf("failed to check user_version: %v", err)
	}
	if version != 1 {
		t.Errorf("expected user_version=1, got: %d", version)
	}
}

func TestOpenAll(t *testing.T) {
	tmpDir := t.TempDir()

	paths := Paths{
		Short:     filepath.Join(tmpDir, "short.db"),
		Journal: filepath.Join(tmpDir, "journal.db"),
		Knowledge: filepath.Join(tmpDir, "knowledge.db"),
		Archive:   filepath.Join(tmpDir, "archive.db"),
	}

	conn, err := OpenAll(paths)
	if err != nil {
		t.Fatalf("OpenAll failed: %v", err)
	}
	defer conn.Close()

	// 各ストアの user_version を確認
	stores := []struct {
		name  string
		query string
	}{
		{"main", "PRAGMA main.user_version;"},
		{"journal", "PRAGMA journal.user_version;"},
		{"knowledge", "PRAGMA knowledge.user_version;"},
		{"archive", "PRAGMA archive.user_version;"},
	}

	for _, store := range stores {
		var version int
		if err := conn.DB.QueryRow(store.query).Scan(&version); err != nil {
			t.Errorf("failed to check %s.user_version: %v", store.name, err)
		} else if version != 1 {
			t.Errorf("expected %s.user_version=1, got: %d", store.name, version)
		}
	}
}

func openDB(dsn string) (*sql.DB, error) {
	return sql.Open("sqlite", dsn)
}
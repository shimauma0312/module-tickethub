package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"           // PostgreSQLドライバ
	_ "github.com/mattn/go-sqlite3" // SQLiteドライバ
)

// InitDatabase はデータベース接続を初期化します
func InitDatabase(config *DatabaseConfig) (*sql.DB, error) {
	var db *sql.DB
	var err error

	if config.Type == "sqlite" {
		// SQLite接続
		db, err = sql.Open("sqlite3", config.DSN())
		if err != nil {
			return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
		}

		// SQLite設定
		_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
		if err != nil {
			return nil, fmt.Errorf("failed to set SQLite pragma: %w", err)
		}
	} else {
		// PostgreSQL接続
		db, err = sql.Open("postgres", config.DSN())
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}
	}

	// 接続確認
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// マイグレーション実行
	if err := migrateDatabase(db, string(config.Type)); err != nil {
		return nil, fmt.Errorf("database migration failed: %w", err)
	}

	return db, nil
}

// migrateDatabase はデータベースのマイグレーションを実行します
func migrateDatabase(db *sql.DB, dbType string) error {
	log.Println("Running database migrations...")

	var err error
	if dbType == "sqlite" {
		// SQLiteスキーマ作成
		err = createSQLiteSchema(db)
	} else {
		// PostgreSQLスキーマ作成
		err = createPostgreSQLSchema(db)
	}

	if err != nil {
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}

// createSQLiteSchema はSQLite用のテーブルスキーマを作成します
func createSQLiteSchema(db *sql.DB) error {
	// Issues テーブル作成
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS issues (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			body TEXT NOT NULL,
			status TEXT NOT NULL,
			assignee_id INTEGER,
			creator_id INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			is_draft BOOLEAN NOT NULL DEFAULT FALSE,
			milestone_id INTEGER
		);
	`)
	if err != nil {
		return err
	}

	// Labels テーブル作成
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS issue_labels (
			issue_id INTEGER NOT NULL,
			label TEXT NOT NULL,
			PRIMARY KEY (issue_id, label),
			FOREIGN KEY (issue_id) REFERENCES issues(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	// 全文検索インデックス作成
	_, err = db.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS issue_fts USING fts5(
			title, body, content='issues', content_rowid='id'
		);
	`)
	if err != nil {
		return err
	}

	// トリガー作成（自動インデックス更新）
	_, err = db.Exec(`
		CREATE TRIGGER IF NOT EXISTS issues_ai AFTER INSERT ON issues BEGIN
			INSERT INTO issue_fts(rowid, title, body) VALUES (new.id, new.title, new.body);
		END;
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TRIGGER IF NOT EXISTS issues_ad AFTER DELETE ON issues BEGIN
			INSERT INTO issue_fts(issue_fts, rowid, title, body) VALUES('delete', old.id, old.title, old.body);
		END;
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TRIGGER IF NOT EXISTS issues_au AFTER UPDATE ON issues BEGIN
			INSERT INTO issue_fts(issue_fts, rowid, title, body) VALUES('delete', old.id, old.title, old.body);
			INSERT INTO issue_fts(rowid, title, body) VALUES (new.id, new.title, new.body);
		END;
	`)
	return err
}

// createPostgreSQLSchema はPostgreSQL用のテーブルスキーマを作成します
func createPostgreSQLSchema(db *sql.DB) error {
	// Issues テーブル作成
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS issues (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			body TEXT NOT NULL,
			status TEXT NOT NULL,
			assignee_id BIGINT,
			creator_id BIGINT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			is_draft BOOLEAN NOT NULL DEFAULT FALSE,
			milestone_id BIGINT
		);
	`)
	if err != nil {
		return err
	}

	// Labels テーブル作成
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS issue_labels (
			issue_id BIGINT NOT NULL,
			label TEXT NOT NULL,
			PRIMARY KEY (issue_id, label),
			FOREIGN KEY (issue_id) REFERENCES issues(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	// 全文検索用インデックス作成
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_issues_title_body ON issues USING GIN(
			to_tsvector('english', title || ' ' || body)
		);
	`)
	return err
}

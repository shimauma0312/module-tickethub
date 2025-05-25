package config

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitGormDB はGORMを使用してデータベース接続を初期化します
func InitGormDB(config *DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// ロガー設定
	loggerConfig := logger.Config{
		SlowThreshold: 200, // 200ms以上のクエリを「遅い」と判定
		LogLevel:      logger.Warn,
		Colorful:      true,
	}

	gormConfig := &gorm.Config{
		Logger: logger.New(
			defaultLogger{},
			loggerConfig,
		),
	}

	dialector, err := config.GetGormDialector()
	if err != nil {
		return nil, fmt.Errorf("failed to get gorm dialector: %w", err)
	}

	db, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		// エラーメッセージをより具体的に
		return nil, fmt.Errorf("failed to connect to database using GORM (%s): %w", config.Type, err)
	}

	// SQLite固有の設定
	if config.Type == SQLite {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get *sql.DB for SQLite specific settings: %w", err)
		}
		// 外部キー制約を有効化
		if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
			return nil, fmt.Errorf("failed to set SQLite pragma (foreign_keys): %w", err)
		}
		// WALモードを有効化（パフォーマンス向上）
		if _, err := sqlDB.Exec("PRAGMA journal_mode = WAL;"); err != nil {
			return nil, fmt.Errorf("failed to set WAL mode for SQLite: %w", err)
		}
	}

	// 接続確認
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 接続プールの設定
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

// デフォルトのロガー実装
type defaultLogger struct{}

func (l defaultLogger) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

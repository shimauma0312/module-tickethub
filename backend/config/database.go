package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBType はサポートされているデータベースタイプを表す型です
type DBType string

const (
	// SQLite はSQLiteデータベースを表します
	SQLite DBType = "sqlite"
	// Postgres はPostgreSQLデータベースを表します
	Postgres DBType = "postgres"
	// SQLServer はMicrosoft SQL Serverデータベースを表します
	SQLServer DBType = "sqlserver"
)

// DatabaseConfig はデータベース接続に必要な設定を保持する構造体
type DatabaseConfig struct {
	// データベースタイプ (sqlite, postgres, sqlserver)
	Type DBType
	// SQLite設定
	SQLitePath string
	// PostgreSQL/SQL Server共通設定
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	// PostgreSQL固有設定
	SSLMode string
	// SQL Server固有設定
	Instance        string
	TrustServerCert bool
}

// NewDatabaseConfig は環境変数から設定を読み込み、DatabaseConfigを生成します
func NewDatabaseConfig() (*DatabaseConfig, error) {
	dbTypeStr := os.Getenv("DB_TYPE")
	if dbTypeStr == "" {
		dbTypeStr = "sqlite" // デフォルトはSQLite
	}

	// 文字列をDBType型に変換
	var dbType DBType
	switch dbTypeStr {
	case "sqlite":
		dbType = SQLite
	case "postgres":
		dbType = Postgres
	case "sqlserver":
		dbType = SQLServer
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbTypeStr)
	}

	config := &DatabaseConfig{
		Type: dbType,
	}

	switch dbType {
	case SQLite:
		// SQLite設定
		sqlitePath := os.Getenv("SQLITE_DB_PATH")
		if sqlitePath == "" {
			sqlitePath = "./data/tickethub.db" // デフォルトパス
		}

		// ディレクトリが存在しない場合は作成
		dir := filepath.Dir(sqlitePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create directory for SQLite: %w", err)
			}
		}

		config.SQLitePath = sqlitePath

	case Postgres:
		// PostgreSQL設定
		config.Host = os.Getenv("DB_HOST")
		if config.Host == "" {
			config.Host = "localhost"
		}

		config.Port = os.Getenv("DB_PORT")
		if config.Port == "" {
			config.Port = "5432"
		}

		config.User = os.Getenv("DB_USER")
		if config.User == "" {
			config.User = "postgres"
		}

		config.Password = os.Getenv("DB_PASSWORD")
		if config.Password == "" {
			config.Password = "postgres"
		}

		config.DBName = os.Getenv("DB_NAME")
		if config.DBName == "" {
			config.DBName = "tickethub"
		}

		config.SSLMode = os.Getenv("DB_SSLMODE")
		if config.SSLMode == "" {
			config.SSLMode = "disable"
		}

	case SQLServer:
		// SQL Server設定
		config.Host = os.Getenv("DB_HOST")
		if config.Host == "" {
			config.Host = "localhost"
		}

		config.Port = os.Getenv("DB_PORT")
		if config.Port == "" {
			config.Port = "1433" // SQL Serverデフォルトポート
		}

		config.User = os.Getenv("DB_USER")
		if config.User == "" {
			config.User = "sa"
		}

		config.Password = os.Getenv("DB_PASSWORD")
		// 空のパスワードも許容

		config.DBName = os.Getenv("DB_NAME")
		if config.DBName == "" {
			config.DBName = "tickethub"
		}

		config.Instance = os.Getenv("DB_INSTANCE")

		// 文字列からbool値に変換
		trustServerCertStr := os.Getenv("DB_TRUST_SERVER_CERT")
		if trustServerCertStr != "" {
			trustServerCert, err := strconv.ParseBool(trustServerCertStr)
			if err == nil {
				config.TrustServerCert = trustServerCert
			}
		} else {
			config.TrustServerCert = true // デフォルトでは証明書を信頼
		}
	}

	return config, nil
}

// DSN はPostgreSQL用のデータソース名を生成します
func (c *DatabaseConfig) DSN() string {
	if c.Type != Postgres {
		return "" // PostgreSQL以外の場合は空文字列を返す
	}
	// 環境変数からSSLModeを取得、デフォルトは "disable"
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable" // デフォルトSSLモード
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, sslMode)
}

// GetGormDialector は設定に基づいて適切なGORM Dialectorを返します
func (c *DatabaseConfig) GetGormDialector() (gorm.Dialector, error) {
	switch c.Type {
	case SQLite:
		return sqlite.Open(c.SQLitePath), nil
	case Postgres:
		return postgres.Open(c.DSN()), nil
	case SQLServer:
		// SQL Server 接続文字列の構築例 (環境変数から読み込むことを推奨)
		// dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		// 	c.User, c.Password, c.Host, c.Port, c.DBName)
		// if c.Instance != "" {
		// 	dsn = fmt.Sprintf("sqlserver://%s:%s@%s\\\\%s:%s?database=%s",
		// 		c.User, c.Password, c.Host, c.Instance, c.Port, c.DBName)
		// }
		// return sqlserver.Open(dsn), nil // sqlserver driver のインポートが必要
		return nil, fmt.Errorf("SQL Server is not yet supported") // 現状維持
	default:
		return nil, fmt.Errorf("unsupported database type: %s", c.Type)
	}
}

// ConnectDB は設定に基づいてデータベース接続を確立します
// この関数はGORMの導入により、InitGormDB に置き換えられるべきです。
// 互換性のために残す場合でも、内容はGORMベースに修正するか、非推奨とすべきです。
func ConnectDB(config *DatabaseConfig) (*sql.DB, error) {
	// GORMを使用する場合、この関数は直接的には使用されません。
	// InitGormDB を使用して *gorm.DB を取得し、そこから *sql.DB を取得できます。
	// 例: gormDB, err := InitGormDB(config)
	//     if err != nil { return nil, err }
	//     sqlDB, err := gormDB.DB()
	//     if err != nil { return nil, err }
	//     return sqlDB, nil
	return nil, fmt.Errorf("ConnectDB is deprecated, use InitGormDB and then db.DB()")
}

// CloseDB はデータベース接続を安全に閉じます
func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

// GetDB は現在のデータベース接続を返します（この関数はGORM導入により不要になる可能性があります）
func GetDB() *sql.DB {
	// グローバルなDB接続変数を使用する場合はここで返す
	return nil
}

// SetDB は現在のデータベース接続を設定します（この関数はGORM導入により不要になる可能性があります）
func SetDB(db *sql.DB) {
	// グローバルなDB接続変数を設定する場合はここで行う
}

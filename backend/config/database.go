package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

// DSN はデータベース接続文字列を生成します
func (c *DatabaseConfig) DSN() string {
	switch c.Type {
	case SQLite:
		return c.SQLitePath

	case Postgres:
		// PostgreSQL接続文字列
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)

	case SQLServer:
		// SQL Server接続文字列
		dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s",
			c.Host, c.User, c.Password, c.DBName)

		if c.Port != "" && c.Port != "1433" {
			dsn += fmt.Sprintf(";port=%s", c.Port)
		}

		if c.Instance != "" {
			dsn += fmt.Sprintf(";instance=%s", c.Instance)
		}

		if c.TrustServerCert {
			dsn += ";TrustServerCertificate=true"
		}

		return dsn

	default:
		return ""
	}
}

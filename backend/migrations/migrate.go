package migrations

import (
	"fmt"
	"log"

	"github.com/shimauma0312/module-tickethub/config" // モジュールパスを修正
	"gorm.io/gorm"
)

// MigrateDB はデータベースのマイグレーションを実行します
// GORMの自動マイグレーション機能を使用するように変更
func MigrateDB(db *gorm.DB, dbConfig *config.DatabaseConfig, direction string) error {
	log.Println("Attempting to run GORM database migrations...")

	// GORMのマイグレーションは通常 "up" のみ。direction は将来的な拡張性のために残すことも検討。
	// ここでは GormMigrate が方向を意識しない前提で実装。
	if direction == "up" {
		err := GormMigrate(db) // gorm_migrate.go の GormMigrate を呼び出す
		if err != nil {
			return fmt.Errorf("failed to apply GORM migrations: %w", err)
		}
		log.Println("Successfully applied GORM migrations")
	} else if direction == "down" {
		// GORMの自動マイグレーションには直接的な "down" がないため、
		// 必要であれば手動でのロールバック処理や別のマイグレーションツールとの連携を検討。
		log.Println("GORM auto-migration does not support 'down' direction directly. Manual intervention may be required.")
		return fmt.Errorf("GORM auto-migration does not support 'down' direction directly")
	} else {
		return fmt.Errorf("invalid migration direction: %s. Only 'up' is meaningfully supported with GORM auto-migration", direction)
	}

	return nil
}

// setupSQLiteMigration, setupPostgresMigration, getMigrationsPath は GORM移行により不要となるため削除
/*
// setupSQLiteMigration はSQLite用のマイグレーション設定を行います
// ... (省略) ...
// getMigrationsPath はデータベースタイプに応じたマイグレーションファイルのパスを返します
// ... (省略) ...
*/

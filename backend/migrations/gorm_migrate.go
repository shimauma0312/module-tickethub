package migrations

import (
	"fmt"
	"log"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"gorm.io/gorm"
)

// GormMigrate はGORMを使用したマイグレーションを実行します
func GormMigrate(db *gorm.DB) error {
	log.Println("Running GORM database migrations...")

	// Issue関連のマイグレーション
	if err := models.AutoMigrateIssue(db); err != nil {
		return fmt.Errorf("failed to migrate issue tables: %w", err)
	}

	// システム設定のマイグレーション
	if err := models.AutoMigrateSystemSettings(db); err != nil {
		return fmt.Errorf("failed to migrate system settings table: %w", err)
	}

	// アクティビティログとバックアップ情報のマイグレーション
	if err := models.AutoMigrateActivityLog(db); err != nil {
		return fmt.Errorf("failed to migrate activity log tables: %w", err)
	}

	log.Println("GORM database migration completed successfully")
	return nil
}

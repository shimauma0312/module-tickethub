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

	// 他のモデルも同様に追加

	log.Println("GORM database migration completed successfully")
	return nil
}

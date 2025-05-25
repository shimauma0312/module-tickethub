package gorm

import (
	"context"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
	"gorm.io/gorm"
)

// systemSettingsRepository はGORMを使用したSystemSettingsRepositoryの実装
type systemSettingsRepository struct {
	db *gorm.DB
}

// NewSystemSettingsRepository は新しいSystemSettingsRepositoryインスタンスを作成
func NewSystemSettingsRepository(db *gorm.DB) repositories.SystemSettingsRepository {
	return &systemSettingsRepository{
		db: db,
	}
}

// Get はシステム設定を取得します
func (r *systemSettingsRepository) Get(ctx context.Context) (*models.SystemSettings, error) {
	var settings models.SystemSettings
	if err := r.db.WithContext(ctx).First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// デフォルト設定を返す
			return models.NewDefaultSystemSettings(), nil
		}
		return nil, err
	}
	return &settings, nil
}

// CreateOrUpdate はシステム設定を作成または更新します
func (r *systemSettingsRepository) CreateOrUpdate(ctx context.Context, settings *models.SystemSettings) error {
	var existing models.SystemSettings

	if err := r.db.WithContext(ctx).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// レコードが存在しない場合は新規作成
			return r.db.WithContext(ctx).Create(settings).Error
		}
		return err
	}

	// 既存レコードを更新
	settings.ID = existing.ID
	return r.db.WithContext(ctx).Save(settings).Error
}

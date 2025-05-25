package gorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserSettingsRepository はGORMベースのUserSettingsリポジトリ実装
type UserSettingsRepository struct {
	db *gorm.DB
}

// NewUserSettingsRepository は新しいUserSettingsRepositoryを作成します
func NewUserSettingsRepository(db *gorm.DB) (*UserSettingsRepository, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &UserSettingsRepository{db: db}, nil
}

// GetByUserID はユーザーIDによってUserSettingsを取得します
func (r *UserSettingsRepository) GetByUserID(ctx context.Context, userID int64) (*models.UserSettings, error) {
	var settings models.UserSettings
	// ユーザーIDで検索し、見つからなければデフォルト値を設定して返す (またはエラーを返すかは要件次第)
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// レコードが存在しない場合、デフォルト設定で新しいレコードを作成または返す
			// ここではエラーを返す代わりに、空またはデフォルトのUserSettingsを返すことも考えられる
			return &models.UserSettings{UserID: userID}, nil // またはエラーを返す: fmt.Errorf("user settings not found for user_id %d: %w", userID, err)
		}
		return nil, fmt.Errorf("failed to get user settings by user_id: %w", err)
	}
	return &settings, nil
}

// CreateOrUpdate はUserSettingsを作成または更新します
func (r *UserSettingsRepository) CreateOrUpdate(ctx context.Context, settings *models.UserSettings) error {
	// OnConflictを使用して、存在すれば更新、存在しなければ作成する
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{ // PostgreSQL, SQLiteで動作
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"language", "timezone", "notifications_enabled", "email_notifications_enabled", "push_notifications_enabled", "updated_at"}),
	}).Create(settings).Error; err != nil {
		return fmt.Errorf("failed to create or update user settings: %w", err)
	}
	return nil
}

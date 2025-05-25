package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shimauma0312/module-tickethub/models"
	"gorm.io/gorm"
)

// PasswordResetRepository はGORMベースのPasswordResetリポジトリ実装
type PasswordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository は新しいPasswordResetRepositoryを作成します
func NewPasswordResetRepository(db *gorm.DB) (*PasswordResetRepository, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &PasswordResetRepository{db: db}, nil
}

// Create は新しいPasswordResetを作成します
func (r *PasswordResetRepository) Create(ctx context.Context, reset *models.PasswordReset) error {
	if err := r.db.WithContext(ctx).Create(reset).Error; err != nil {
		return fmt.Errorf("failed to create password reset: %w", err)
	}
	return nil
}

// GetByID はIDによってPasswordResetを取得します
func (r *PasswordResetRepository) GetByID(ctx context.Context, id int64) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	if err := r.db.WithContext(ctx).First(&reset, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("password reset with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get password reset by id: %w", err)
	}
	return &reset, nil
}

// GetByToken はトークン文字列によってPasswordResetを取得します
func (r *PasswordResetRepository) GetByToken(ctx context.Context, token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&reset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("password reset with token %s not found", token)
		}
		return nil, fmt.Errorf("failed to get password reset by token: %w", err)
	}
	return &reset, nil
}

// GetByUserID はユーザーIDによってPasswordResetの一覧を取得します (有効なもののみを返す例)
func (r *PasswordResetRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.PasswordReset, error) {
	var resets []*models.PasswordReset
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND used_at IS NULL AND expires_at > ?", userID, time.Now()).
		Find(&resets).Error; err != nil {
		return nil, fmt.Errorf("failed to get password resets by user_id: %w", err)
	}
	return resets, nil
}

// Update は既存のPasswordResetを更新します
func (r *PasswordResetRepository) Update(ctx context.Context, reset *models.PasswordReset) error {
	if err := r.db.WithContext(ctx).Save(reset).Error; err != nil {
		return fmt.Errorf("failed to update password reset: %w", err)
	}
	return nil
}

// Delete はPasswordResetを削除します (物理削除)
func (r *PasswordResetRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&models.PasswordReset{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete password reset: %w", err)
	}
	return nil
}

// RevokeAllForUser はユーザーの全PasswordResetトークンを無効化します (UsedAtを更新)
func (r *PasswordResetRepository) RevokeAllForUser(ctx context.Context, userID int64) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&models.PasswordReset{}).
		Where("user_id = ? AND used_at IS NULL", userID).
		Update("used_at", now).Error; err != nil {
		return fmt.Errorf("failed to revoke all password resets for user_id %d: %w", userID, err)
	}
	return nil
}

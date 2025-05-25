package gorm

import (
	"context"
	"time"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
	"gorm.io/gorm"
)

// backupRepository はGORMを使用したBackupRepositoryの実装
type backupRepository struct {
	db *gorm.DB
}

// NewBackupRepository は新しいBackupRepositoryインスタンスを作成
func NewBackupRepository(db *gorm.DB) repositories.BackupRepository {
	return &backupRepository{
		db: db,
	}
}

// Create は新しいBackupInfoを作成します
func (r *backupRepository) Create(ctx context.Context, backup *models.BackupInfo) error {
	return r.db.WithContext(ctx).Create(backup).Error
}

// List はBackupInfoの一覧を取得します
func (r *backupRepository) List(ctx context.Context, page, limit int) ([]*models.BackupInfo, int, error) {
	var backups []*models.BackupInfo
	var total int64

	// 総件数を取得
	if err := r.db.WithContext(ctx).Model(&models.BackupInfo{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&backups).Error; err != nil {
		return nil, 0, err
	}

	return backups, int(total), nil
}

// GetByID はIDによってBackupInfoを取得します
func (r *backupRepository) GetByID(ctx context.Context, id int64) (*models.BackupInfo, error) {
	var backup models.BackupInfo
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&backup).Error; err != nil {
		return nil, err
	}
	return &backup, nil
}

// Update はBackupInfoを更新します
func (r *backupRepository) Update(ctx context.Context, backup *models.BackupInfo) error {
	return r.db.WithContext(ctx).Save(backup).Error
}

// Delete はBackupInfoを削除します
func (r *backupRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.BackupInfo{}).Error
}

// DeleteOldBackups は古いバックアップを削除します
func (r *backupRepository) DeleteOldBackups(ctx context.Context, retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	return r.db.WithContext(ctx).Where("created_at < ? AND status = ?", cutoffDate, "completed").Delete(&models.BackupInfo{}).Error
}

// GetLatestBackup は最新のバックアップを取得します
func (r *backupRepository) GetLatestBackup(ctx context.Context) (*models.BackupInfo, error) {
	var backup models.BackupInfo
	if err := r.db.WithContext(ctx).Where("status = ?", "completed").Order("created_at DESC").First(&backup).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // バックアップが存在しない場合はnilを返す
		}
		return nil, err
	}
	return &backup, nil
}

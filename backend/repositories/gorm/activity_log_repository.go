package gorm

import (
	"context"
	"runtime"
	"time"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
	"gorm.io/gorm"
)

// activityLogRepository はGORMを使用したActivityLogRepositoryの実装
type activityLogRepository struct {
	db *gorm.DB
}

// NewActivityLogRepository は新しいActivityLogRepositoryインスタンスを作成
func NewActivityLogRepository(db *gorm.DB) repositories.ActivityLogRepository {
	return &activityLogRepository{
		db: db,
	}
}

// Create は新しいActivityLogを作成します
func (r *activityLogRepository) Create(ctx context.Context, log *models.ActivityLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// List は条件に一致するActivityLogの一覧を取得します
func (r *activityLogRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ActivityLog, int, error) {
	var logs []*models.ActivityLog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ActivityLog{})

	// フィルターを適用
	if userID, ok := filter["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if action, ok := filter["action"]; ok {
		query = query.Where("action = ?", action)
	}
	if resourceType, ok := filter["resource_type"]; ok {
		query = query.Where("resource_type = ?", resourceType)
	}
	if resourceID, ok := filter["resource_id"]; ok {
		query = query.Where("resource_id = ?", resourceID)
	}
	if fromDate, ok := filter["from_date"]; ok {
		query = query.Where("created_at >= ?", fromDate)
	}
	if toDate, ok := filter["to_date"]; ok {
		query = query.Where("created_at <= ?", toDate)
	}

	// 総件数を取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, int(total), nil
}

// GetMetrics はシステムメトリクスを取得します
func (r *activityLogRepository) GetMetrics(ctx context.Context) (*models.SystemMetrics, error) {
	metrics := &models.SystemMetrics{
		GeneratedAt: time.Now(),
	}

	// ユーザー関連メトリクス
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&metrics.TotalUsers).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("is_active = ?", true).Count(&metrics.ActiveUsers).Error; err != nil {
		return nil, err
	}

	// Issue関連メトリクス
	if err := r.db.WithContext(ctx).Model(&models.Issue{}).Count(&metrics.TotalIssues).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&models.Issue{}).Where("status = ?", "open").Count(&metrics.OpenIssues).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&models.Issue{}).Where("status = ?", "closed").Count(&metrics.ClosedIssues).Error; err != nil {
		return nil, err
	}

	// Discussion関連メトリクス
	if err := r.db.WithContext(ctx).Model(&models.Discussion{}).Count(&metrics.TotalDiscussions).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&models.Discussion{}).Where("status = ?", "open").Count(&metrics.OpenDiscussions).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&models.Discussion{}).Where("status IN ?", []string{"closed", "answered"}).Count(&metrics.ClosedDiscussions).Error; err != nil {
		return nil, err
	}

	// コメント関連メトリクス
	if err := r.db.WithContext(ctx).Model(&models.Comment{}).Count(&metrics.TotalComments).Error; err != nil {
		return nil, err
	}

	// システムリソース情報
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	metrics.MemoryUsage = int64(memStats.Alloc)

	// データベースサイズの取得（SQLiteの場合）
	// 実際の実装では、データベースの種類に応じて異なる方法を使用する必要があります
	metrics.DatabaseSize = 0
	metrics.DatabaseConnections = 1

	// 最後のバックアップ時刻を取得
	var lastBackup models.BackupInfo
	if err := r.db.WithContext(ctx).Model(&models.BackupInfo{}).
		Where("status = ?", "completed").
		Order("created_at DESC").
		First(&lastBackup).Error; err == nil {
		metrics.LastBackupAt = lastBackup.CreatedAt
	}

	return metrics, nil
}

// DeleteOldLogs は古いログを削除します
func (r *activityLogRepository) DeleteOldLogs(ctx context.Context, retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	return r.db.WithContext(ctx).Where("created_at < ?", cutoffDate).Delete(&models.ActivityLog{}).Error
}

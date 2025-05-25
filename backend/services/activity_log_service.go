package services

import (
	"context"
	"encoding/json"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
)

// ActivityLogService はアクティビティログサービス
type ActivityLogService struct {
	activityLogRepo repositories.ActivityLogRepository
}

// NewActivityLogService は新しいActivityLogServiceを作成します
func NewActivityLogService(activityLogRepo repositories.ActivityLogRepository) *ActivityLogService {
	return &ActivityLogService{
		activityLogRepo: activityLogRepo,
	}
}

// LogActivity はアクティビティをログに記録します
func (s *ActivityLogService) LogActivity(ctx context.Context, userID int64, username string, action models.LogAction, resourceType models.ResourceType, resourceID int64, ipAddress, userAgent string, details interface{}) error {
	var detailsJSON string
	if details != nil {
		if detailsBytes, err := json.Marshal(details); err == nil {
			detailsJSON = string(detailsBytes)
		}
	}

	log := models.NewActivityLog(userID, username, action, resourceType, resourceID, ipAddress, userAgent, detailsJSON)
	return s.activityLogRepo.Create(ctx, log)
}

// GetActivityLogs はアクティビティログの一覧を取得します
func (s *ActivityLogService) GetActivityLogs(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ActivityLog, int, error) {
	return s.activityLogRepo.List(ctx, filter, page, limit)
}

// GetSystemMetrics はシステムメトリクスを取得します
func (s *ActivityLogService) GetSystemMetrics(ctx context.Context) (*models.SystemMetrics, error) {
	return s.activityLogRepo.GetMetrics(ctx)
}

// CleanOldLogs は古いログを削除します
func (s *ActivityLogService) CleanOldLogs(ctx context.Context, retentionDays int) error {
	return s.activityLogRepo.DeleteOldLogs(ctx, retentionDays)
}

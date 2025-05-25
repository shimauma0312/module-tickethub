package models

import (
	"time"

	"gorm.io/gorm"
)

// ActivityLog はシステム内でのアクティビティログを表す構造体
type ActivityLog struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Username     string    `json:"username"`
	Action       string    `json:"action"`        // 実行されたアクション
	ResourceType string    `json:"resource_type"` // 操作対象のリソースタイプ（user/issue/discussion/etc）
	ResourceID   int64     `json:"resource_id"`   // 操作対象のリソースID
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	Details      string    `json:"details,omitempty"` // 追加の詳細情報（JSON文字列）
	CreatedAt    time.Time `json:"created_at"`
}

// LogAction は事前定義されたアクション名
type LogAction string

const (
	// ユーザー関連
	ActionUserCreated         LogAction = "user.created"
	ActionUserUpdated         LogAction = "user.updated"
	ActionUserDeleted         LogAction = "user.deleted"
	ActionUserActivated       LogAction = "user.activated"
	ActionUserDeactivated     LogAction = "user.deactivated"
	ActionUserLogin           LogAction = "user.login"
	ActionUserLogout          LogAction = "user.logout"
	ActionUserPasswordChanged LogAction = "user.password_changed"

	// Issue関連
	ActionIssueCreated    LogAction = "issue.created"
	ActionIssueUpdated    LogAction = "issue.updated"
	ActionIssueDeleted    LogAction = "issue.deleted"
	ActionIssueClosed     LogAction = "issue.closed"
	ActionIssueReopened   LogAction = "issue.reopened"
	ActionIssueAssigned   LogAction = "issue.assigned"
	ActionIssueUnassigned LogAction = "issue.unassigned"

	// Discussion関連
	ActionDiscussionCreated  LogAction = "discussion.created"
	ActionDiscussionUpdated  LogAction = "discussion.updated"
	ActionDiscussionDeleted  LogAction = "discussion.deleted"
	ActionDiscussionClosed   LogAction = "discussion.closed"
	ActionDiscussionReopened LogAction = "discussion.reopened"

	// コメント関連
	ActionCommentCreated LogAction = "comment.created"
	ActionCommentUpdated LogAction = "comment.updated"
	ActionCommentDeleted LogAction = "comment.deleted"

	// ラベル関連
	ActionLabelCreated LogAction = "label.created"
	ActionLabelUpdated LogAction = "label.updated"
	ActionLabelDeleted LogAction = "label.deleted"

	// マイルストーン関連
	ActionMilestoneCreated LogAction = "milestone.created"
	ActionMilestoneUpdated LogAction = "milestone.updated"
	ActionMilestoneDeleted LogAction = "milestone.deleted"

	// システム関連
	ActionSystemSettingsUpdated    LogAction = "system.settings_updated"
	ActionSystemBackupCreated      LogAction = "system.backup_created"
	ActionSystemBackupRestored     LogAction = "system.backup_restored"
	ActionSystemMaintenanceToggled LogAction = "system.maintenance_toggled"
)

// ResourceType は事前定義されたリソースタイプ
type ResourceType string

const (
	ResourceUser       ResourceType = "user"
	ResourceIssue      ResourceType = "issue"
	ResourceDiscussion ResourceType = "discussion"
	ResourceComment    ResourceType = "comment"
	ResourceLabel      ResourceType = "label"
	ResourceMilestone  ResourceType = "milestone"
	ResourceSystem     ResourceType = "system"
)

// NewActivityLog は新しいアクティビティログを作成する
func NewActivityLog(userID int64, username string, action LogAction, resourceType ResourceType, resourceID int64, ipAddress, userAgent, details string) *ActivityLog {
	return &ActivityLog{
		UserID:       userID,
		Username:     username,
		Action:       string(action),
		ResourceType: string(resourceType),
		ResourceID:   resourceID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Details:      details,
		CreatedAt:    time.Now(),
	}
}

// SystemMetrics はシステムメトリクスを表す構造体
type SystemMetrics struct {
	TotalUsers          int64     `json:"total_users"`
	ActiveUsers         int64     `json:"active_users"`
	TotalIssues         int64     `json:"total_issues"`
	OpenIssues          int64     `json:"open_issues"`
	ClosedIssues        int64     `json:"closed_issues"`
	TotalDiscussions    int64     `json:"total_discussions"`
	OpenDiscussions     int64     `json:"open_discussions"`
	ClosedDiscussions   int64     `json:"closed_discussions"`
	TotalComments       int64     `json:"total_comments"`
	DatabaseSize        int64     `json:"database_size"` // バイト単位
	DatabaseConnections int       `json:"database_connections"`
	MemoryUsage         int64     `json:"memory_usage"`   // バイト単位
	CPUUsage            float64   `json:"cpu_usage"`      // パーセンテージ
	DiskUsage           int64     `json:"disk_usage"`     // バイト単位
	DiskAvailable       int64     `json:"disk_available"` // バイト単位
	LastBackupAt        time.Time `json:"last_backup_at,omitempty"`
	UptimeSeconds       int64     `json:"uptime_seconds"`
	GeneratedAt         time.Time `json:"generated_at"`
}

// BackupInfo はバックアップ情報を表す構造体
type BackupInfo struct {
	ID          int64     `json:"id"`
	Filename    string    `json:"filename"`
	FilePath    string    `json:"file_path"`
	FileSize    int64     `json:"file_size"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"` // creating/completed/failed
}

// NewBackupInfo は新しいバックアップ情報を作成する
func NewBackupInfo(filename, filepath string, createdBy int64, description string) *BackupInfo {
	return &BackupInfo{
		Filename:    filename,
		FilePath:    filepath,
		CreatedBy:   createdBy,
		Description: description,
		Status:      "creating",
		CreatedAt:   time.Now(),
	}
}

// Complete はバックアップの完了を記録する
func (b *BackupInfo) Complete(fileSize int64) {
	b.Status = "completed"
	b.FileSize = fileSize
}

// Fail はバックアップの失敗を記録する
func (b *BackupInfo) Fail() {
	b.Status = "failed"
}

// AutoMigrateActivityLog はActivityLog、SystemMetrics、BackupInfoテーブルのマイグレーションを実行します
func AutoMigrateActivityLog(db *gorm.DB) error {
	if err := db.AutoMigrate(&ActivityLog{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&BackupInfo{}); err != nil {
		return err
	}
	return nil
}

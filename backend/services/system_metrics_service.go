package services

import (
	"context"
	"runtime"
	"time"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
)

// SystemMetricsService はシステムメトリクス収集サービス
type SystemMetricsService struct {
	userRepo       repositories.UserRepository
	issueRepo      repositories.IssueRepository
	discussionRepo repositories.DiscussionRepository
	commentRepo    repositories.CommentRepository
	backupRepo     repositories.BackupRepository
	startTime      time.Time
}

// NewSystemMetricsService は新しいSystemMetricsServiceを作成します
func NewSystemMetricsService(
	userRepo repositories.UserRepository,
	issueRepo repositories.IssueRepository,
	discussionRepo repositories.DiscussionRepository,
	commentRepo repositories.CommentRepository,
	backupRepo repositories.BackupRepository,
) *SystemMetricsService {
	return &SystemMetricsService{
		userRepo:       userRepo,
		issueRepo:      issueRepo,
		discussionRepo: discussionRepo,
		commentRepo:    commentRepo,
		backupRepo:     backupRepo,
		startTime:      time.Now(),
	}
}

// GetSystemMetrics はシステムメトリクスを収集して返します
func (s *SystemMetricsService) GetSystemMetrics() (*models.SystemMetrics, error) {
	ctx := context.Background()
	
	metrics := &models.SystemMetrics{
		GeneratedAt: time.Now(),
	}

	// ユーザー統計
	totalUsers, err := s.userRepo.CountUsers(ctx)
	if err != nil {
		return nil, err
	}
	metrics.TotalUsers = totalUsers

	// アクティブユーザー（30日以内にログインしたユーザー）
	activeUsers, err := s.userRepo.CountActiveUsers(ctx, 30)
	if err != nil {
		return nil, err
	}
	metrics.ActiveUsers = activeUsers

	// Issue統計
	totalIssues, err := s.issueRepo.CountIssues(ctx)
	if err != nil {
		return nil, err
	}
	metrics.TotalIssues = totalIssues

	openIssues, err := s.issueRepo.CountOpenIssues(ctx)
	if err != nil {
		return nil, err
	}
	metrics.OpenIssues = openIssues
	metrics.ClosedIssues = totalIssues - openIssues

	// Discussion統計
	totalDiscussions, err := s.discussionRepo.CountDiscussions(ctx)
	if err != nil {
		return nil, err
	}
	metrics.TotalDiscussions = totalDiscussions

	openDiscussions, err := s.discussionRepo.CountOpenDiscussions(ctx)
	if err != nil {
		return nil, err
	}
	metrics.OpenDiscussions = openDiscussions
	metrics.ClosedDiscussions = totalDiscussions - openDiscussions

	// コメント統計
	totalComments, err := s.commentRepo.CountComments(ctx)
	if err != nil {
		return nil, err
	}
	metrics.TotalComments = totalComments

	// 最後のバックアップ時刻
	lastBackup, err := s.backupRepo.GetLatestBackup(ctx)
	if err == nil && lastBackup != nil {
		metrics.LastBackupAt = lastBackup.CreatedAt
	}

	// システムリソース情報
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	metrics.MemoryUsage = int64(m.Alloc)

	// CPU使用率（簡易的な値、より正確な値は外部ライブラリが必要）
	metrics.CPUUsage = 0.0 // TODO: 実際のCPU使用率を取得

	// データベースサイズ（簡易的な値）
	metrics.DatabaseSize = 0 // TODO: 実際のデータベースサイズを取得

	// データベース接続数（簡易的な値）
	metrics.DatabaseConnections = 1 // TODO: 実際の接続数を取得

	// ディスク使用量（簡易的な値）
	metrics.DiskUsage = 0     // TODO: 実際のディスク使用量を取得
	metrics.DiskAvailable = 0 // TODO: 実際の利用可能ディスク容量を取得

	// アップタイム
	metrics.UptimeSeconds = int64(time.Since(s.startTime).Seconds())

	return metrics, nil
}

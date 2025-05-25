package gorm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"gorm.io/gorm"
)

// IssueRepository はGORMを使用したIssueリポジトリの実装
type IssueRepository struct {
	db *gorm.DB
}

// NewIssueRepository は新しいIssueRepositoryを生成します
func NewIssueRepository(db *gorm.DB) *IssueRepository {
	return &IssueRepository{
		db: db,
	}
}

// Create は新しいIssueを作成します
func (r *IssueRepository) Create(ctx context.Context, issue *models.Issue) error {
	// 現在時刻を設定
	now := time.Now()
	if issue.CreatedAt.IsZero() {
		issue.CreatedAt = now
	}
	if issue.UpdatedAt.IsZero() {
		issue.UpdatedAt = now
	}

	// モデル変換
	gormIssue := models.IssueFromModel(issue)

	// トランザクション開始
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Issueの保存
	if err := tx.Create(gormIssue).Error; err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}

	// IDを設定
	issue.ID = gormIssue.ID

	// ラベルの保存
	if len(issue.Labels) > 0 {
		// ラベルのGORMモデルを作成
		labels := make([]models.IssueLabel, 0, len(issue.Labels))
		for _, label := range issue.Labels {
			labels = append(labels, models.IssueLabel{
				IssueID: issue.ID,
				Label:   label,
			})
		}

		if err := tx.Create(&labels).Error; err != nil {
			return fmt.Errorf("failed to create issue labels: %w", err)
		}
	}

	// トランザクションをコミット
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByID はIDによってIssueを取得します
func (r *IssueRepository) GetByID(ctx context.Context, id int64) (*models.Issue, error) {
	var gormIssue models.IssueGorm

	// Issueの取得（ラベルも同時に取得）
	err := r.db.WithContext(ctx).
		Preload("Labels").
		First(&gormIssue, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("issue with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	// 標準モデルに変換
	return gormIssue.ToModel(), nil
}

// List は条件に一致するIssueの一覧を取得します
func (r *IssueRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Issue, int, error) {
	var gormIssues []models.IssueGorm
	var total int64

	// クエリビルダー
	query := r.db.WithContext(ctx).Model(&models.IssueGorm{})

	// フィルター適用
	if filter != nil {
		for k, v := range filter {
			switch k {
			case "status":
				query = query.Where("status = ?", v)
			case "assignee_id":
				query = query.Where("assignee_id = ?", v)
			case "creator_id":
				query = query.Where("creator_id = ?", v)
			case "milestone_id":
				query = query.Where("milestone_id = ?", v)
			case "is_draft":
				query = query.Where("is_draft = ?", v)
			case "label":
				// ラベルでフィルタリング（サブクエリ使用）
				query = query.Where("id IN (SELECT issue_id FROM issue_labels WHERE label = ?)", v)
			}
		}
	}

	// 総件数を取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count issues: %w", err)
	}

	// データ取得
	err := query.
		Preload("Labels").
		Order("updated_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&gormIssues).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list issues: %w", err)
	}

	// 標準モデルに変換
	issues := make([]*models.Issue, len(gormIssues))
	for i, gormIssue := range gormIssues {
		issues[i] = gormIssue.ToModel()
	}

	return issues, int(total), nil
}

// Update は既存のIssueを更新します
func (r *IssueRepository) Update(ctx context.Context, issue *models.Issue) error {
	// 更新日時を更新
	issue.UpdatedAt = time.Now()

	// トランザクション開始
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// GORMモデルに変換
	gormIssue := models.IssueFromModel(issue)

	// Issueの更新
	if err := tx.Model(&models.IssueGorm{}).
		Where("id = ?", issue.ID).
		Updates(map[string]interface{}{
			"title":        issue.Title,
			"body":         issue.Body,
			"status":       issue.Status,
			"assignee_id":  gormIssue.AssigneeID,
			"updated_at":   issue.UpdatedAt,
			"is_draft":     issue.IsDraft,
			"milestone_id": gormIssue.MilestoneID,
		}).Error; err != nil {
		return fmt.Errorf("failed to update issue: %w", err)
	}

	// ラベルの更新（古いラベルを削除して新しいラベルを挿入）
	if err := tx.Where("issue_id = ?", issue.ID).Delete(&models.IssueLabel{}).Error; err != nil {
		return fmt.Errorf("failed to delete old labels: %w", err)
	}

	// 新しいラベルを挿入
	if len(issue.Labels) > 0 {
		labels := make([]models.IssueLabel, 0, len(issue.Labels))
		for _, label := range issue.Labels {
			labels = append(labels, models.IssueLabel{
				IssueID: issue.ID,
				Label:   label,
			})
		}

		if err := tx.Create(&labels).Error; err != nil {
			return fmt.Errorf("failed to create issue labels: %w", err)
		}
	}

	// トランザクションをコミット
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Delete はIssueを削除します
func (r *IssueRepository) Delete(ctx context.Context, id int64) error {
	// GORMの論理削除を使用
	if err := r.db.WithContext(ctx).Delete(&models.IssueGorm{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete issue: %w", err)
	}
	return nil
}

// Search はIssueの全文検索を行います
func (r *IssueRepository) Search(ctx context.Context, query string, page, limit int) ([]*models.Issue, int, error) {
	var gormIssues []models.IssueGorm
	var total int64

	// クエリの特殊文字をエスケープ
	searchQuery := "%" + strings.Replace(query, "%", "\\%", -1) + "%"

	// LIKEを使用して検索（データベース非依存の方法）
	dbQuery := r.db.WithContext(ctx).
		Model(&models.IssueGorm{}).
		Where("title LIKE ? OR body LIKE ?", searchQuery, searchQuery)

	// 総件数を取得
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// データ取得
	err := dbQuery.
		Preload("Labels").
		Order("updated_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&gormIssues).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to search issues: %w", err)
	}

	// 標準モデルに変換
	issues := make([]*models.Issue, len(gormIssues))
	for i, gormIssue := range gormIssues {
		issues[i] = gormIssue.ToModel()
	}

	return issues, int(total), nil
}

// GetAll はすべてのIssueを取得します（検索インデックス構築用）
func (r *IssueRepository) GetAll(ctx context.Context) ([]*models.Issue, error) {
	var gormIssues []models.IssueGorm

	// すべてのIssueを取得
	err := r.db.WithContext(ctx).
		Preload("Labels").
		Find(&gormIssues).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get all issues: %w", err)
	}

	// 標準モデルに変換
	issues := make([]*models.Issue, len(gormIssues))
	for i, gormIssue := range gormIssues {
		issues[i] = gormIssue.ToModel()
	}

	return issues, nil
}

// CountIssues は総Issue数を取得します
func (r *IssueRepository) CountIssues(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.IssueGorm{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count issues: %w", err)
	}
	return count, nil
}

// CountOpenIssues はオープンなIssue数を取得します
func (r *IssueRepository) CountOpenIssues(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.IssueGorm{}).Where("status = ?", "open").Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count open issues: %w", err)
	}
	return count, nil
}

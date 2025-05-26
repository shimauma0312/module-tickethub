package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shimauma0312/module-tickethub/backend/models"
)

// SearchRepository はSQLiteデータベースを使用した検索リポジトリの実装
type SearchRepository struct {
	db *sql.DB
}

// NewSearchRepository は新しいSearchRepositoryインスタンスを作成します
func NewSearchRepository(db *sql.DB) (*SearchRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return &SearchRepository{db: db}, nil
}

// SearchIssues は指定した検索クエリに基づいて課題を検索します
func (r *SearchRepository) SearchIssues(ctx context.Context, query string, page, limit int) ([]*models.Issue, int, error) {
	// SQLiteでの検索実装 (基本的な実装例)
	// 実際のアプリケーション要件に合わせて拡張する必要があります

	// カウントクエリ
	countQuery := `
		SELECT COUNT(*) FROM issues 
		WHERE title LIKE ? OR body LIKE ?
	`
	countParams := []interface{}{
		"%" + query + "%",
		"%" + query + "%",
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, countParams...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// 検索クエリ
	searchQuery := `
		SELECT id, title, body, status, creator_id, created_at, updated_at
		FROM issues 
		WHERE title LIKE ? OR body LIKE ?
		ORDER BY updated_at DESC
		LIMIT ? OFFSET ?
	`
	searchParams := []interface{}{
		"%" + query + "%",
		"%" + query + "%",
		limit,
		(page - 1) * limit,
	}

	rows, err := r.db.QueryContext(ctx, searchQuery, searchParams...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search issues: %w", err)
	}
	defer rows.Close()

	var issues []*models.Issue
	for rows.Next() {
		issue := &models.Issue{}
		err := rows.Scan(
			&issue.ID,
			&issue.Title,
			&issue.Body,
			&issue.Status,
			&issue.CreatorID,
			&issue.CreatedAt,
			&issue.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan issue row: %w", err)
		}
		issues = append(issues, issue)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating search results: %w", err)
	}

	return issues, total, nil
}

// SearchComments はコメント内のテキストを検索します
func (r *SearchRepository) SearchComments(ctx context.Context, query string, page, limit int) ([]*models.Comment, int, error) {
	// SQLiteでのコメント検索実装（簡略化）
	// 実際のアプリケーション要件に合わせて拡張する必要があります

	// カウントクエリ
	countQuery := `
		SELECT COUNT(*) FROM comments 
		WHERE body LIKE ?
	`

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, "%"+query+"%").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comment search results: %w", err)
	}

	// 検索クエリ
	searchQuery := `
		SELECT id, body, created_at, updated_at, creator_id, target_id
		FROM comments 
		WHERE body LIKE ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, searchQuery, "%"+query+"%", limit, (page-1)*limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search comments: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.Body,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.CreatorID,
			&comment.TargetID,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan comment row: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating comment search results: %w", err)
	}

	return comments, total, nil
}

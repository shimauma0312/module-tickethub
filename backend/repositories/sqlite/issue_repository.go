package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/shimauma0312/module-tickethub/backend/models"
)

// IssueRepository はSQLite固有のIssueリポジトリ実装です
type IssueRepository struct {
	db *sql.DB
}

// NewIssueRepository は新しいSQLite IssueRepositoryを生成します
func NewIssueRepository(db *sql.DB) *IssueRepository {
	return &IssueRepository{
		db: db,
	}
}

// Create は新しいIssueを作成します
func (r *IssueRepository) Create(ctx context.Context, issue *models.Issue) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Issuesテーブルに挿入
	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO issues (title, body, status, assignee_id, creator_id, created_at, updated_at, is_draft, milestone_id) 
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		issue.Title,
		issue.Body,
		issue.Status,
		sql.NullInt64{Int64: issue.AssigneeID, Valid: issue.AssigneeID > 0},
		issue.CreatorID,
		issue.CreatedAt,
		issue.UpdatedAt,
		issue.IsDraft,
		sql.NullInt64{Int64: issue.MilestoneID, Valid: issue.MilestoneID > 0},
	)
	if err != nil {
		return fmt.Errorf("failed to insert issue: %w", err)
	}

	// 生成されたIDを取得
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	issue.ID = id

	// ラベルを挿入
	for _, label := range issue.Labels {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO issue_labels (issue_id, label) VALUES (?, ?)`,
			issue.ID,
			label,
		)
		if err != nil {
			return fmt.Errorf("failed to insert label: %w", err)
		}
	}

	return tx.Commit()
}

// GetByID はIDによってIssueを取得します
func (r *IssueRepository) GetByID(ctx context.Context, id int64) (*models.Issue, error) {
	// まずIssueの基本情報を取得
	issue := &models.Issue{}
	var assigneeID, milestoneID sql.NullInt64
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, title, body, status, assignee_id, creator_id, created_at, updated_at, is_draft, milestone_id
		 FROM issues WHERE id = ?`,
		id,
	).Scan(
		&issue.ID,
		&issue.Title,
		&issue.Body,
		&issue.Status,
		&assigneeID,
		&issue.CreatorID,
		&issue.CreatedAt,
		&issue.UpdatedAt,
		&issue.IsDraft,
		&milestoneID,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("issue not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query issue: %w", err)
	}

	// NULLの場合は0を設定
	if assigneeID.Valid {
		issue.AssigneeID = assigneeID.Int64
	}
	if milestoneID.Valid {
		issue.MilestoneID = milestoneID.Int64
	}

	// ラベルを取得
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT label FROM issue_labels WHERE issue_id = ?`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query labels: %w", err)
	}
	defer rows.Close()

	issue.Labels = []string{}
	for rows.Next() {
		var label string
		if err := rows.Scan(&label); err != nil {
			return nil, fmt.Errorf("failed to scan label: %w", err)
		}
		issue.Labels = append(issue.Labels, label)
	}

	return issue, nil
}

// List は条件に一致するIssueの一覧を取得します
func (r *IssueRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Issue, int, error) {
	// クエリの構築
	query := "SELECT id, title, body, status, assignee_id, creator_id, created_at, updated_at, is_draft, milestone_id FROM issues"
	countQuery := "SELECT COUNT(*) FROM issues"

	// WHERE句の条件を構築
	conditions := []string{}
	args := []interface{}{}

	for key, value := range filter {
		switch key {
		case "status":
			conditions = append(conditions, "status = ?")
			args = append(args, value)
		case "assignee_id":
			conditions = append(conditions, "assignee_id = ?")
			args = append(args, value)
		case "creator_id":
			conditions = append(conditions, "creator_id = ?")
			args = append(args, value)
		case "milestone_id":
			conditions = append(conditions, "milestone_id = ?")
			args = append(args, value)
		case "is_draft":
			conditions = append(conditions, "is_draft = ?")
			args = append(args, value)
		}
	}

	// WHERE句を追加
	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	// ページネーション
	query += " ORDER BY updated_at DESC LIMIT ? OFFSET ?"
	offset := (page - 1) * limit
	args = append(args, limit, offset)

	// 総件数のカウント
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count issues: %w", err)
	}

	// データの取得
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query issues: %w", err)
	}
	defer rows.Close()

	issues := []*models.Issue{}
	issueIDs := []int64{}

	for rows.Next() {
		issue := &models.Issue{}
		var assigneeID, milestoneID sql.NullInt64

		if err := rows.Scan(
			&issue.ID,
			&issue.Title,
			&issue.Body,
			&issue.Status,
			&assigneeID,
			&issue.CreatorID,
			&issue.CreatedAt,
			&issue.UpdatedAt,
			&issue.IsDraft,
			&milestoneID,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan issue: %w", err)
		}

		if assigneeID.Valid {
			issue.AssigneeID = assigneeID.Int64
		}
		if milestoneID.Valid {
			issue.MilestoneID = milestoneID.Int64
		}

		issue.Labels = []string{}
		issues = append(issues, issue)
		issueIDs = append(issueIDs, issue.ID)
	}

	// ラベルの一括取得（N+1問題を回避）
	if len(issues) > 0 {
		placeholders := strings.Repeat("?,", len(issueIDs))
		placeholders = placeholders[:len(placeholders)-1] // 最後のカンマを除去

		labelQuery := fmt.Sprintf("SELECT issue_id, label FROM issue_labels WHERE issue_id IN (%s)", placeholders)
		args := make([]interface{}, len(issueIDs))
		for i, id := range issueIDs {
			args[i] = id
		}

		labelRows, err := r.db.QueryContext(ctx, labelQuery, args...)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to query labels: %w", err)
		}
		defer labelRows.Close()

		// ID -> Issue のマップを作成
		issueMap := make(map[int64]*models.Issue)
		for _, issue := range issues {
			issueMap[issue.ID] = issue
		}

		// ラベルをマッピング
		for labelRows.Next() {
			var issueID int64
			var label string
			if err := labelRows.Scan(&issueID, &label); err != nil {
				return nil, 0, fmt.Errorf("failed to scan label: %w", err)
			}

			if issue, ok := issueMap[issueID]; ok {
				issue.Labels = append(issue.Labels, label)
			}
		}
	}

	return issues, total, nil
}

// Update は既存のIssueを更新します
func (r *IssueRepository) Update(ctx context.Context, issue *models.Issue) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 更新時間を現在時刻に設定
	issue.UpdatedAt = time.Now()

	// Issuesテーブルを更新
	_, err = tx.ExecContext(
		ctx,
		`UPDATE issues 
		 SET title = ?, body = ?, status = ?, assignee_id = ?, updated_at = ?, is_draft = ?, milestone_id = ?
		 WHERE id = ?`,
		issue.Title,
		issue.Body,
		issue.Status,
		sql.NullInt64{Int64: issue.AssigneeID, Valid: issue.AssigneeID > 0},
		issue.UpdatedAt,
		issue.IsDraft,
		sql.NullInt64{Int64: issue.MilestoneID, Valid: issue.MilestoneID > 0},
		issue.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update issue: %w", err)
	}

	// ラベルを削除してから再挿入
	_, err = tx.ExecContext(ctx, "DELETE FROM issue_labels WHERE issue_id = ?", issue.ID)
	if err != nil {
		return fmt.Errorf("failed to delete labels: %w", err)
	}

	// 新しいラベルを挿入
	for _, label := range issue.Labels {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO issue_labels (issue_id, label) VALUES (?, ?)`,
			issue.ID,
			label,
		)
		if err != nil {
			return fmt.Errorf("failed to insert label: %w", err)
		}
	}

	return tx.Commit()
}

// Delete はIssueを削除します
func (r *IssueRepository) Delete(ctx context.Context, id int64) error {
	// issue_labelsはON DELETE CASCADEで自動的に削除される
	_, err := r.db.ExecContext(ctx, "DELETE FROM issues WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete issue: %w", err)
	}
	return nil
}

// Search はIssueの全文検索を行います
func (r *IssueRepository) Search(ctx context.Context, query string, page, limit int) ([]*models.Issue, int, error) {
	// SQLiteのFTS5を使用した検索
	searchQuery := `
		SELECT i.id, i.title, i.body, i.status, i.assignee_id, i.creator_id, i.created_at, i.updated_at, i.is_draft, i.milestone_id
		FROM issue_fts f
		JOIN issues i ON f.rowid = i.id
		WHERE issue_fts MATCH ?
		ORDER BY rank
		LIMIT ? OFFSET ?
	`

	countQuery := `
		SELECT COUNT(*)
		FROM issue_fts f
		JOIN issues i ON f.rowid = i.id
		WHERE issue_fts MATCH ?
	`

	// 検索クエリの整形
	searchPattern := strings.Join(strings.Fields(query), " OR ")

	// 総件数の取得
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// ページネーションの設定
	offset := (page - 1) * limit

	// 検索実行
	rows, err := r.db.QueryContext(ctx, searchQuery, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search issues: %w", err)
	}
	defer rows.Close()

	issues := []*models.Issue{}
	issueIDs := []int64{}

	for rows.Next() {
		issue := &models.Issue{}
		var assigneeID, milestoneID sql.NullInt64

		if err := rows.Scan(
			&issue.ID,
			&issue.Title,
			&issue.Body,
			&issue.Status,
			&assigneeID,
			&issue.CreatorID,
			&issue.CreatedAt,
			&issue.UpdatedAt,
			&issue.IsDraft,
			&milestoneID,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan issue: %w", err)
		}

		if assigneeID.Valid {
			issue.AssigneeID = assigneeID.Int64
		}
		if milestoneID.Valid {
			issue.MilestoneID = milestoneID.Int64
		}

		issue.Labels = []string{}
		issues = append(issues, issue)
		issueIDs = append(issueIDs, issue.ID)
	}

	// ラベルの一括取得（N+1問題を回避）
	if len(issues) > 0 {
		placeholders := strings.Repeat("?,", len(issueIDs))
		placeholders = placeholders[:len(placeholders)-1] // 最後のカンマを除去

		labelQuery := fmt.Sprintf("SELECT issue_id, label FROM issue_labels WHERE issue_id IN (%s)", placeholders)
		args := make([]interface{}, len(issueIDs))
		for i, id := range issueIDs {
			args[i] = id
		}

		labelRows, err := r.db.QueryContext(ctx, labelQuery, args...)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to query labels: %w", err)
		}
		defer labelRows.Close()

		// ID -> Issue のマップを作成
		issueMap := make(map[int64]*models.Issue)
		for _, issue := range issues {
			issueMap[issue.ID] = issue
		}

		// ラベルをマッピング
		for labelRows.Next() {
			var issueID int64
			var label string
			if err := labelRows.Scan(&issueID, &label); err != nil {
				return nil, 0, fmt.Errorf("failed to scan label: %w", err)
			}

			if issue, ok := issueMap[issueID]; ok {
				issue.Labels = append(issue.Labels, label)
			}
		}
	}

	return issues, total, nil
}

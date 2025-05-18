package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/shimauma0312/module-tickethub/backend/models"
)

// IssueRepository はPostgreSQL固有のIssueリポジトリ実装です
type IssueRepository struct {
	db *sql.DB
}

// NewIssueRepository は新しいPostgreSQL IssueRepositoryを生成します
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
	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO issues (title, body, status, assignee_id, creator_id, created_at, updated_at, is_draft, milestone_id) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		issue.Title,
		issue.Body,
		issue.Status,
		sql.NullInt64{Int64: issue.AssigneeID, Valid: issue.AssigneeID > 0},
		issue.CreatorID,
		issue.CreatedAt,
		issue.UpdatedAt,
		issue.IsDraft,
		sql.NullInt64{Int64: issue.MilestoneID, Valid: issue.MilestoneID > 0},
	).Scan(&issue.ID)

	if err != nil {
		return fmt.Errorf("failed to insert issue: %w", err)
	}

	// ラベルを挿入
	for _, label := range issue.Labels {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO issue_labels (issue_id, label) VALUES ($1, $2)`,
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
		 FROM issues WHERE id = $1`,
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
		`SELECT label FROM issue_labels WHERE issue_id = $1`,
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
	paramCount := 1

	for key, value := range filter {
		switch key {
		case "status":
			conditions = append(conditions, fmt.Sprintf("status = $%d", paramCount))
			args = append(args, value)
			paramCount++
		case "assignee_id":
			conditions = append(conditions, fmt.Sprintf("assignee_id = $%d", paramCount))
			args = append(args, value)
			paramCount++
		case "creator_id":
			conditions = append(conditions, fmt.Sprintf("creator_id = $%d", paramCount))
			args = append(args, value)
			paramCount++
		case "milestone_id":
			conditions = append(conditions, fmt.Sprintf("milestone_id = $%d", paramCount))
			args = append(args, value)
			paramCount++
		case "is_draft":
			conditions = append(conditions, fmt.Sprintf("is_draft = $%d", paramCount))
			args = append(args, value)
			paramCount++
		}
	}

	// WHERE句を追加
	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	// ページネーション
	query += fmt.Sprintf(" ORDER BY updated_at DESC LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
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
		// PostgreSQLではANY演算子を使用してIN句を実装
		labelQuery := "SELECT issue_id, label FROM issue_labels WHERE issue_id = ANY($1)"

		labelRows, err := r.db.QueryContext(ctx, labelQuery, issueIDs)
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
		 SET title = $1, body = $2, status = $3, assignee_id = $4, updated_at = $5, is_draft = $6, milestone_id = $7
		 WHERE id = $8`,
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
	_, err = tx.ExecContext(ctx, "DELETE FROM issue_labels WHERE issue_id = $1", issue.ID)
	if err != nil {
		return fmt.Errorf("failed to delete labels: %w", err)
	}

	// 新しいラベルを挿入
	for _, label := range issue.Labels {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO issue_labels (issue_id, label) VALUES ($1, $2)`,
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
	_, err := r.db.ExecContext(ctx, "DELETE FROM issues WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete issue: %w", err)
	}
	return nil
}

// Search はIssueの全文検索を行います
func (r *IssueRepository) Search(ctx context.Context, query string, page, limit int) ([]*models.Issue, int, error) {
	// PostgreSQLのto_tsqueryを使用した全文検索
	// 複数の単語をORで結合
	terms := strings.Fields(query)
	for i, term := range terms {
		terms[i] = term + ":*" // 前方一致検索のためにワイルドカードを追加
	}
	tsQuery := strings.Join(terms, " | ")

	searchQuery := `
		SELECT i.id, i.title, i.body, i.status, i.assignee_id, i.creator_id, i.created_at, i.updated_at, i.is_draft, i.milestone_id
		FROM issues i
		WHERE to_tsvector('english', i.title || ' ' || i.body) @@ to_tsquery('english', $1)
		ORDER BY ts_rank(to_tsvector('english', i.title || ' ' || i.body), to_tsquery('english', $1)) DESC
		LIMIT $2 OFFSET $3
	`

	countQuery := `
		SELECT COUNT(*)
		FROM issues i
		WHERE to_tsvector('english', i.title || ' ' || i.body) @@ to_tsquery('english', $1)
	`

	// 総件数の取得
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, tsQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// ページネーションの設定
	offset := (page - 1) * limit

	// 検索実行
	rows, err := r.db.QueryContext(ctx, searchQuery, tsQuery, limit, offset)
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
		// PostgreSQLではANY演算子を使用してIN句を実装
		labelQuery := "SELECT issue_id, label FROM issue_labels WHERE issue_id = ANY($1)"

		labelRows, err := r.db.QueryContext(ctx, labelQuery, issueIDs)
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

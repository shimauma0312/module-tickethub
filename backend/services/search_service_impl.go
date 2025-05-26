package services

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
	"gorm.io/gorm"
)

// searchServiceImpl は SearchService インターフェースの実装
type searchServiceImpl struct {
	db          *gorm.DB
	sqlDB       *sql.DB
	issueRepo   repositories.IssueRepository
	commentRepo repositories.CommentRepository
}

// NewSearchService は SearchService の新しいインスタンスを作成する
func NewSearchService(db *gorm.DB, issueRepo repositories.IssueRepository, commentRepo repositories.CommentRepository) (SearchService, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	service := &searchServiceImpl{
		db:          db,
		sqlDB:       sqlDB,
		issueRepo:   issueRepo,
		commentRepo: commentRepo,
	}

	// 初期化時にFTS5テーブルを作成
	if err := service.initFTS5Tables(); err != nil {
		return nil, fmt.Errorf("failed to initialize FTS5 tables: %w", err)
	}

	return service, nil
}

// initFTS5Tables はFTS5検索テーブルを初期化する
func (s *searchServiceImpl) initFTS5Tables() error {
	// FTS5拡張モジュールが有効かチェック
	var fts5Enabled bool
	row := s.sqlDB.QueryRow("SELECT * FROM sqlite_master WHERE type='table' AND name='sqlite_master'")
	if err := row.Scan(&fts5Enabled); err != nil {
		if err != sql.ErrNoRows {
			return fmt.Errorf("failed to check FTS5 availability: %w", err)
		}
	}

	// FTS5テーブルの作成
	// issue_search テーブル
	_, err := s.sqlDB.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS issue_search USING fts5(
			doc_id UNINDEXED,
			title,
			body,
			tokenize = 'unicode61 remove_diacritics 1'
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create issue_search table: %w", err)
	}

	// comment_search テーブル
	_, err = s.sqlDB.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS comment_search USING fts5(
			doc_id UNINDEXED,
			target_id UNINDEXED,
			body,
			tokenize = 'unicode61 remove_diacritics 1'
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create comment_search table: %w", err)
	}

	return nil
}

// Search は指定されたクエリに基づいてコンテンツを検索する
func (s *searchServiceImpl) Search(ctx context.Context, query models.SearchQuery) (*models.SearchResults, error) {
	results := &models.SearchResults{
		Query:       query.Query,
		CurrentPage: query.Offset/query.Limit + 1,
		Results:     []models.SearchResult{},
	}

	if query.Limit <= 0 {
		query.Limit = 20 // デフォルト値
	}

	// SQLクエリのパラメータ
	ftsQuery := s.formatFTSQuery(query.Query)
	statusCondition := ""
	if query.Status != "" && query.Status != "all" {
		statusCondition = fmt.Sprintf("AND i.status = '%s'", query.Status)
	}

	// ラベル条件の構築
	labelCondition := ""
	labelArgs := []interface{}{}
	if len(query.Labels) > 0 {
		placeholders := make([]string, len(query.Labels))
		for i, label := range query.Labels {
			placeholders[i] = "?"
			labelArgs = append(labelArgs, label)
		}
		labelCondition = fmt.Sprintf("AND (SELECT COUNT(*) FROM json_each(i.labels) WHERE json_each.value IN (%s)) = %d",
			strings.Join(placeholders, ","), len(query.Labels))
	}

	// 担当者条件
	assigneeCondition := ""
	if query.AssigneeID > 0 {
		assigneeCondition = fmt.Sprintf("AND i.assignee_id = %d", query.AssigneeID)
	}

	// 作成者条件
	creatorCondition := ""
	if query.CreatorID > 0 {
		creatorCondition = fmt.Sprintf("AND i.creator_id = %d", query.CreatorID)
	}

	// Issue検索
	issueQuery := fmt.Sprintf(`
		SELECT 
			i.id, i.title, i.body, i.status, i.labels, i.assignee_id, 
			i.creator_id, i.created_at, i.updated_at,
			highlight(issue_search, 1, '<mark>', '</mark>') as title_highlight,
			highlight(issue_search, 2, '<mark>', '</mark>') as body_highlight,
			rank
		FROM issue_search
		JOIN issues i ON issue_search.doc_id = i.id
		WHERE issue_search MATCH ?
		%s %s %s %s
		ORDER BY rank
		LIMIT ? OFFSET ?
	`, statusCondition, labelCondition, assigneeCondition, creatorCondition)

	// Issue検索の実行
	issueArgs := append([]interface{}{ftsQuery}, labelArgs...)
	issueArgs = append(issueArgs, query.Limit, query.Offset)

	issueRows, err := s.sqlDB.QueryContext(ctx, issueQuery, issueArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to search issues: %w", err)
	}
	defer issueRows.Close()

	// Issue結果の処理
	for issueRows.Next() {
		var result models.SearchResult
		var titleHighlight, bodyHighlight string
		var labelsJSON string
		var createdAt, updatedAt string
		var rank float64

		err := issueRows.Scan(
			&result.ID, &result.Title, &result.Body, &result.Status, &labelsJSON,
			&result.AssigneeID, &result.CreatorID, &createdAt, &updatedAt,
			&titleHighlight, &bodyHighlight, &rank,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan issue result: %w", err)
		}

		// 日時のパース
		result.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		result.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		// ラベルの変換
		result.Labels = []string{}
		// JSON文字列からラベル配列への変換処理

		// 検索結果のメタデータ設定
		result.Type = models.SearchResultTypeIssue
		result.Rank = rank
		result.Highlighted = s.combineHighlights(titleHighlight, bodyHighlight)
		result.Snippet = s.createSnippet(result.Highlighted)

		results.Results = append(results.Results, result)
	}

	// コメント検索クエリ
	commentQuery := fmt.Sprintf(`
		SELECT 
			c.id, c.body, c.creator_id, c.created_at, c.updated_at, c.target_id,
			highlight(comment_search, 2, '<mark>', '</mark>') as body_highlight,
			rank
		FROM comment_search
		JOIN comments c ON comment_search.doc_id = c.id
		WHERE comment_search MATCH ?
		AND c.type = 'issue'
		%s
		ORDER BY rank
		LIMIT ? OFFSET ?
	`, creatorCondition)

	// コメント検索の実行
	commentArgs := []interface{}{ftsQuery, query.Limit, query.Offset}

	commentRows, err := s.sqlDB.QueryContext(ctx, commentQuery, commentArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to search comments: %w", err)
	}
	defer commentRows.Close()

	// コメント結果の処理
	for commentRows.Next() {
		var result models.SearchResult
		var bodyHighlight string
		var createdAt, updatedAt string
		var rank float64

		err := commentRows.Scan(
			&result.ID, &result.Body, &result.CreatorID, &createdAt, &updatedAt,
			&result.TargetID, &bodyHighlight, &rank,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment result: %w", err)
		}

		// 日時のパース
		result.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		result.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		// 関連するIssueのタイトルを取得
		issue, err := s.issueRepo.GetByID(ctx, result.TargetID)
		if err == nil && issue != nil {
			result.Title = fmt.Sprintf("Comment on: %s", issue.Title)
		} else {
			result.Title = "Comment"
		}

		// 検索結果のメタデータ設定
		result.Type = models.SearchResultTypeComment
		result.Rank = rank
		result.Highlighted = bodyHighlight
		result.Snippet = s.createSnippet(bodyHighlight)

		results.Results = append(results.Results, result)
	}

	// 総結果数のカウント取得
	var totalCount int
	countQuery := `
		SELECT 
			(SELECT COUNT(*) FROM issue_search WHERE issue_search MATCH ?) +
			(SELECT COUNT(*) FROM comment_search WHERE comment_search MATCH ?)
	`
	err = s.sqlDB.QueryRowContext(ctx, countQuery, ftsQuery, ftsQuery).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	results.TotalCount = totalCount
	results.TotalPages = (totalCount + query.Limit - 1) / query.Limit

	// 検索結果をランクでソート
	s.sortResultsByRank(results.Results)

	return results, nil
}

// IndexIssue は指定されたIssueをインデックスに追加または更新する
func (s *searchServiceImpl) IndexIssue(ctx context.Context, issue *models.Issue) error {
	// 既存のインデックスを削除
	if err := s.DeleteFromIndex(ctx, "issue", issue.ID); err != nil {
		return fmt.Errorf("failed to delete existing index: %w", err)
	}

	// 新しいインデックスを追加
	_, err := s.sqlDB.ExecContext(ctx, `
		INSERT INTO issue_search (doc_id, title, body)
		VALUES (?, ?, ?)
	`, issue.ID, issue.Title, issue.Body)

	if err != nil {
		return fmt.Errorf("failed to index issue: %w", err)
	}

	return nil
}

// IndexComment は指定されたCommentをインデックスに追加または更新する
func (s *searchServiceImpl) IndexComment(ctx context.Context, comment *models.Comment) error {
	// コメントタイプが"issue"のみをインデックス化
	if comment.Type != "issue" {
		return nil
	}

	// 既存のインデックスを削除
	if err := s.DeleteFromIndex(ctx, "comment", comment.ID); err != nil {
		return fmt.Errorf("failed to delete existing index: %w", err)
	}

	// 新しいインデックスを追加
	_, err := s.sqlDB.ExecContext(ctx, `
		INSERT INTO comment_search (doc_id, target_id, body)
		VALUES (?, ?, ?)
	`, comment.ID, comment.TargetID, comment.Body)

	if err != nil {
		return fmt.Errorf("failed to index comment: %w", err)
	}

	return nil
}

// DeleteFromIndex は指定されたドキュメントをインデックスから削除する
func (s *searchServiceImpl) DeleteFromIndex(ctx context.Context, docType string, docID int64) error {
	var query string

	switch docType {
	case "issue":
		query = "DELETE FROM issue_search WHERE doc_id = ?"
	case "comment":
		query = "DELETE FROM comment_search WHERE doc_id = ?"
	default:
		return fmt.Errorf("unknown document type: %s", docType)
	}

	_, err := s.sqlDB.ExecContext(ctx, query, docID)
	if err != nil {
		return fmt.Errorf("failed to delete from index: %w", err)
	}

	return nil
}

// RebuildIndex はすべてのインデックスを再構築する
func (s *searchServiceImpl) RebuildIndex(ctx context.Context) error {
	// トランザクション開始
	tx, err := s.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// インデックスをクリア
	_, err = tx.ExecContext(ctx, "DELETE FROM issue_search")
	if err != nil {
		return fmt.Errorf("failed to clear issue index: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM comment_search")
	if err != nil {
		return fmt.Errorf("failed to clear comment index: %w", err)
	}

	// すべてのIssueを再インデックス化
	issues, err := s.issueRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all issues: %w", err)
	}

	for _, issue := range issues {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO issue_search (doc_id, title, body)
			VALUES (?, ?, ?)
		`, issue.ID, issue.Title, issue.Body)

		if err != nil {
			return fmt.Errorf("failed to index issue %d: %w", issue.ID, err)
		}
	}

	// すべてのIssueコメントを再インデックス化
	comments, err := s.commentRepo.GetAllOfType(ctx, "issue")
	if err != nil {
		return fmt.Errorf("failed to get all issue comments: %w", err)
	}

	for _, comment := range comments {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO comment_search (doc_id, target_id, body)
			VALUES (?, ?, ?)
		`, comment.ID, comment.TargetID, comment.Body)

		if err != nil {
			return fmt.Errorf("failed to index comment %d: %w", comment.ID, err)
		}
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ParseQuery は検索クエリ文字列を解析する
func (s *searchServiceImpl) ParseQuery(queryString string) models.SearchQuery {
	query := models.SearchQuery{
		Query:  queryString,
		Status: "all", // デフォルト値
		Limit:  20,    // デフォルト値
	}

	// ラベルフィルタの解析: label:bug label:feature
	labelRegex := regexp.MustCompile(`label:(\S+)`)
	labelMatches := labelRegex.FindAllStringSubmatch(queryString, -1)
	for _, match := range labelMatches {
		if len(match) > 1 {
			query.Labels = append(query.Labels, match[1])
			queryString = strings.Replace(queryString, match[0], "", 1)
		}
	}

	// ステータスフィルタの解析: status:open / status:closed
	statusRegex := regexp.MustCompile(`status:(\S+)`)
	statusMatch := statusRegex.FindStringSubmatch(queryString)
	if len(statusMatch) > 1 {
		query.Status = statusMatch[1]
		queryString = strings.Replace(queryString, statusMatch[0], "", 1)
	}

	// 担当者フィルタの解析: assignee:123
	assigneeRegex := regexp.MustCompile(`assignee:(\d+)`)
	assigneeMatch := assigneeRegex.FindStringSubmatch(queryString)
	if len(assigneeMatch) > 1 {
		query.AssigneeID = s.parseID(assigneeMatch[1])
		queryString = strings.Replace(queryString, assigneeMatch[0], "", 1)
	}

	// 作成者フィルタの解析: creator:123
	creatorRegex := regexp.MustCompile(`creator:(\d+)`)
	creatorMatch := creatorRegex.FindStringSubmatch(queryString)
	if len(creatorMatch) > 1 {
		query.CreatorID = s.parseID(creatorMatch[1])
		queryString = strings.Replace(queryString, creatorMatch[0], "", 1)
	}

	// 残りのテキストをクエリとして設定
	query.Query = strings.TrimSpace(queryString)

	return query
}

// 補助関数

// formatFTSQuery はFTS5用の検索クエリを整形する
func (s *searchServiceImpl) formatFTSQuery(query string) string {
	// 基本的なクエリの整形
	formattedQuery := strings.TrimSpace(query)

	// 空のクエリの場合は全件一致を返す
	if formattedQuery == "" {
		return "*"
	}

	// 特殊文字のエスケープ
	escapeChars := []string{"'", "\"", "\\"}
	for _, char := range escapeChars {
		formattedQuery = strings.ReplaceAll(formattedQuery, char, "\\"+char)
	}

	// AND/OR演算子のサポート
	formattedQuery = strings.ReplaceAll(formattedQuery, " AND ", " AND ")
	formattedQuery = strings.ReplaceAll(formattedQuery, " OR ", " OR ")

	// 部分一致検索のためにワイルドカードを追加
	words := strings.Fields(formattedQuery)
	for i, word := range words {
		if !strings.HasPrefix(word, "\"") && !strings.HasSuffix(word, "\"") {
			words[i] = word + "*"
		}
	}

	return strings.Join(words, " ")
}

// parseID は文字列からID値を解析する
func (s *searchServiceImpl) parseID(idStr string) int64 {
	id := int64(0)
	fmt.Sscanf(idStr, "%d", &id)
	return id
}

// combineHighlights はタイトルと本文のハイライトを結合する
func (s *searchServiceImpl) combineHighlights(titleHighlight, bodyHighlight string) string {
	result := ""
	if titleHighlight != "" {
		result += "Title: " + titleHighlight + "\n"
	}
	if bodyHighlight != "" {
		result += "Body: " + bodyHighlight
	}
	return result
}

// createSnippet はハイライトされたテキストからスニペットを作成する
func (s *searchServiceImpl) createSnippet(highlighted string) string {
	// マークアップタグを除去
	snippet := strings.ReplaceAll(highlighted, "<mark>", "")
	snippet = strings.ReplaceAll(snippet, "</mark>", "")

	// 長さを制限
	maxLen := 200
	if len(snippet) > maxLen {
		snippet = snippet[:maxLen] + "..."
	}

	return snippet
}

// sortResultsByRank は検索結果をランクでソートする
func (s *searchServiceImpl) sortResultsByRank(results []models.SearchResult) {
	// 実装省略 - より高度なランキングロジックを実装する場合に使用
}

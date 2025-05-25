package models

import (
	"time"
)

// SearchResultType は検索結果の種類を表す
type SearchResultType string

const (
	// SearchResultTypeIssue は課題検索結果を表す
	SearchResultTypeIssue SearchResultType = "issue"
	// SearchResultTypeComment はコメント検索結果を表す
	SearchResultTypeComment SearchResultType = "comment"
)

// SearchQuery は検索クエリを表す構造体
type SearchQuery struct {
	Query      string   `json:"query"`       // 検索キーワード
	Labels     []string `json:"labels"`      // ラベルによるフィルタ
	Status     string   `json:"status"`      // ステータスによるフィルタ (open/closed/all)
	AssigneeID int64    `json:"assignee_id"` // 担当者IDによるフィルタ
	CreatorID  int64    `json:"creator_id"`  // 作成者IDによるフィルタ
	Limit      int      `json:"limit"`       // 結果の上限数
	Offset     int      `json:"offset"`      // 結果のオフセット
}

// SearchResult は検索結果を表す構造体
type SearchResult struct {
	Type        SearchResultType `json:"type"`        // 結果の種類 (issue/comment)
	ID          int64            `json:"id"`          // コンテンツのID
	Title       string           `json:"title"`       // タイトル (Issueの場合)
	Body        string           `json:"body"`        // 本文
	Snippet     string           `json:"snippet"`     // 検索キーワードを含むスニペット
	Labels      []string         `json:"labels"`      // ラベル (Issueの場合)
	Status      string           `json:"status"`      // ステータス (Issueの場合)
	AssigneeID  int64            `json:"assignee_id"` // 担当者ID (Issueの場合)
	CreatorID   int64            `json:"creator_id"`  // 作成者ID
	CreatedAt   time.Time        `json:"created_at"`  // 作成日時
	UpdatedAt   time.Time        `json:"updated_at"`  // 更新日時
	TargetID    int64            `json:"target_id"`   // 関連対象ID (Commentの場合)
	Rank        float64          `json:"rank"`        // 検索ランキングスコア
	Highlighted string           `json:"highlighted"` // ハイライト付きテキスト
}

// SearchResults は検索結果のリストを表す構造体
type SearchResults struct {
	Results     []SearchResult `json:"results"`      // 検索結果
	TotalCount  int            `json:"total_count"`  // 総結果数
	CurrentPage int            `json:"current_page"` // 現在のページ
	TotalPages  int            `json:"total_pages"`  // 総ページ数
	Query       string         `json:"query"`        // 検索クエリ
}

// SearchIndex はFTS5インデックスのスキーマを表す構造体
type SearchIndex struct {
	ID        int64     `json:"id"`
	DocID     int64     `json:"doc_id"`     // 対象ドキュメントのID
	DocType   string    `json:"doc_type"`   // ドキュメントタイプ (issue/comment)
	Title     string    `json:"title"`      // タイトル (Issueの場合)
	Body      string    `json:"body"`       // 本文
	CreatedAt time.Time `json:"created_at"` // インデックス作成日時
	UpdatedAt time.Time `json:"updated_at"` // インデックス更新日時
}

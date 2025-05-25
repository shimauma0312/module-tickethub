package services

import (
	"context"

	"github.com/shimauma0312/module-tickethub/backend/models"
)

// SearchService は検索機能を提供するインターフェース
type SearchService interface {
	// Search は指定されたクエリに基づいてコンテンツを検索する
	Search(ctx context.Context, query models.SearchQuery) (*models.SearchResults, error)

	// IndexIssue は指定されたIssueをインデックスに追加または更新する
	IndexIssue(ctx context.Context, issue *models.Issue) error

	// IndexComment は指定されたCommentをインデックスに追加または更新する
	IndexComment(ctx context.Context, comment *models.Comment) error

	// DeleteFromIndex は指定されたドキュメントをインデックスから削除する
	DeleteFromIndex(ctx context.Context, docType string, docID int64) error

	// RebuildIndex はすべてのインデックスを再構築する
	RebuildIndex(ctx context.Context) error

	// ParseQuery は検索クエリ文字列を解析する
	ParseQuery(queryString string) models.SearchQuery
}

package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/services"
)

// SearchHandler は検索APIハンドラ
type SearchHandler struct {
	searchService services.SearchService
}

// NewSearchHandler は SearchHandler の新しいインスタンスを作成する
func NewSearchHandler(searchService services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// RegisterRoutes はAPIルートを登録する
func (h *SearchHandler) RegisterRoutes(router *gin.Engine) {
	search := router.Group("/api/search")
	{
		search.GET("", h.Search)
		search.POST("/rebuild-index", h.RebuildIndex)
	}
}

// Search はコンテンツを検索するハンドラ
// @Summary コンテンツを検索する
// @Description 指定したクエリに基づいてイシューやコメントを検索する
// @Tags 検索
// @Accept json
// @Produce json
// @Param query query string false "検索クエリ"
// @Param limit query int false "結果の上限数" default(20)
// @Param offset query int false "結果のオフセット" default(0)
// @Param labels query string false "ラベルフィルタ (カンマ区切り)"
// @Param status query string false "ステータスフィルタ (open/closed/all)" default(all)
// @Param assignee_id query int false "担当者IDフィルタ"
// @Param creator_id query int false "作成者IDフィルタ"
// @Success 200 {object} models.SearchResults "検索結果"
// @Failure 400 {object} ErrorResponse "不正なリクエスト"
// @Failure 500 {object} ErrorResponse "サーバエラー"
// @Router /api/search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	// クエリパラメータの取得
	queryString := c.Query("query")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	labelsStr := c.Query("labels")
	status := c.DefaultQuery("status", "all")
	assigneeID, _ := strconv.ParseInt(c.Query("assignee_id"), 10, 64)
	creatorID, _ := strconv.ParseInt(c.Query("creator_id"), 10, 64)

	// ラベルの解析
	var labels []string
	if labelsStr != "" {
		labels = splitAndTrim(labelsStr, ",")
	}

	// 検索クエリオブジェクトの作成
	query := models.SearchQuery{
		Query:      queryString,
		Labels:     labels,
		Status:     status,
		AssigneeID: assigneeID,
		CreatorID:  creatorID,
		Limit:      limit,
		Offset:     offset,
	}

	// 検索サービスの呼び出し
	results, err := h.searchService.Search(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "検索に失敗しました",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

// RebuildIndex は検索インデックスを再構築するハンドラ
// @Summary 検索インデックスを再構築する
// @Description すべての検索インデックスを最新の状態に再構築する
// @Tags 検索
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse "インデックス再構築成功"
// @Failure 500 {object} ErrorResponse "サーバエラー"
// @Router /api/search/rebuild-index [post]
func (h *SearchHandler) RebuildIndex(c *gin.Context) {
	// 管理者権限のチェック (実際の実装ではミドルウェアで行うことが多い)
	// if !isAdmin(c) {
	// 	c.JSON(http.StatusForbidden, ErrorResponse{
	// 		Error: "権限がありません",
	// 	})
	// 	return
	// }

	err := h.searchService.RebuildIndex(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "インデックス再構築に失敗しました",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "インデックスの再構築が完了しました",
	})
}

// splitAndTrim は文字列をデリミタで分割し、各要素をトリムする
func splitAndTrim(s, delimiter string) []string {
	if s == "" {
		return []string{}
	}

	parts := make([]string, 0)
	for _, part := range strings.Split(s, delimiter) {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

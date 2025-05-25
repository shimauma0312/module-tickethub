package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
)

// DiscussionHandler はDiscussion関連のハンドラーを管理する構造体
type DiscussionHandler struct {
	discussionRepo repositories.DiscussionRepository
	labelRepo      repositories.LabelRepository
	userRepo       repositories.UserRepository
}

// NewDiscussionHandler は新しいDiscussionHandlerを作成します
func NewDiscussionHandler(
	discussionRepo repositories.DiscussionRepository,
	labelRepo repositories.LabelRepository,
	userRepo repositories.UserRepository,
) *DiscussionHandler {
	return &DiscussionHandler{
		discussionRepo: discussionRepo,
		labelRepo:      labelRepo,
		userRepo:       userRepo,
	}
}

// DiscussionRequest はDiscussion作成・更新リクエストのデータ構造
type DiscussionRequest struct {
	Title    string   `json:"title" binding:"required"`
	Body     string   `json:"body"`
	Category string   `json:"category" binding:"required"`
	Labels   []string `json:"labels"`
	IsDraft  bool     `json:"is_draft"`
}

// @Summary Discussion一覧の取得
// @Description Discussionの一覧を取得します
// @Tags discussions
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(10)
// @Param status query string false "ステータス (open/closed/answered)" default("open")
// @Param category query string false "カテゴリ (general/question/announcement/idea)"
// @Param label query string false "ラベル名"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/discussions [get]
func (h *DiscussionHandler) ListDiscussions(c *gin.Context) {
	// クエリパラメータの取得
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.DefaultQuery("status", "open")
	category := c.Query("category")
	label := c.Query("label")

	// フィルタの作成
	filter := map[string]interface{}{
		"status": status,
	}

	// カテゴリが指定されている場合
	if category != "" {
		filter["category"] = category
	}

	// ラベルが指定されている場合は別途処理が必要
	// （実際の実装ではリポジトリレイヤーでラベルフィルタリングを行う）

	// データベースから取得
	discussions, total, err := h.discussionRepo.List(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"discussions": discussions,
		"total":       total,
		"page":        page,
		"limit":       limit,
	})
}

// @Summary Discussionの取得
// @Description 指定されたIDのDiscussionを取得します
// @Tags discussions
// @Accept json
// @Produce json
// @Param id path int true "Discussion ID"
// @Success 200 {object} models.Discussion
// @Router /api/v1/discussions/{id} [get]
func (h *DiscussionHandler) GetDiscussion(c *gin.Context) {
	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion ID format"})
		return
	}

	// データベースから取得
	discussion, err := h.discussionRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if discussion == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Discussion not found"})
		return
	}

	c.JSON(http.StatusOK, discussion)
}

// @Summary Discussionの作成
// @Description 新しいDiscussionを作成します
// @Tags discussions
// @Accept json
// @Produce json
// @Param discussion body DiscussionRequest true "Discussion情報"
// @Success 201 {object} models.Discussion
// @Router /api/v1/discussions [post]
func (h *DiscussionHandler) CreateDiscussion(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// リクエストの解析
	var req DiscussionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Discussionの作成
	discussion := models.NewDiscussion(req.Title, req.Body, req.Category, userID.(int64))
	discussion.IsDraft = req.IsDraft

	// ラベルの設定
	for _, label := range req.Labels {
		discussion.AddLabel(label)
	}

	// データベースに保存
	err := h.discussionRepo.Create(c.Request.Context(), discussion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, discussion)
}

// @Summary Discussionの更新
// @Description 指定されたIDのDiscussionを更新します
// @Tags discussions
// @Accept json
// @Produce json
// @Param id path int true "Discussion ID"
// @Param discussion body DiscussionRequest true "Discussion情報"
// @Success 200 {object} models.Discussion
// @Router /api/v1/discussions/{id} [put]
func (h *DiscussionHandler) UpdateDiscussion(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion ID format"})
		return
	}

	// データベースから取得
	discussion, err := h.discussionRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if discussion == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Discussion not found"})
		return
	}

	// 作成者のみ編集可能（または管理者）
	isAdmin, _ := c.Get("is_admin")
	if discussion.CreatorID != userID.(int64) && !(isAdmin != nil && isAdmin.(bool)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this discussion"})
		return
	}

	// リクエストの解析
	var req DiscussionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Discussionの更新
	discussion.Title = req.Title
	discussion.Body = req.Body
	discussion.Category = req.Category
	discussion.IsDraft = req.IsDraft
	discussion.UpdatedAt = models.CurrentTime()

	// ラベルの更新
	discussion.Labels = []string{} // 既存のラベルをクリア
	for _, label := range req.Labels {
		discussion.AddLabel(label)
	}

	// データベースに保存
	err = h.discussionRepo.Update(c.Request.Context(), discussion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, discussion)
}

// @Summary Discussionの削除
// @Description 指定されたIDのDiscussionを削除します
// @Tags discussions
// @Accept json
// @Produce json
// @Param id path int true "Discussion ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/discussions/{id} [delete]
func (h *DiscussionHandler) DeleteDiscussion(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion ID format"})
		return
	}

	// データベースから取得
	discussion, err := h.discussionRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if discussion == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Discussion not found"})
		return
	}

	// 作成者のみ削除可能（または管理者）
	isAdmin, _ := c.Get("is_admin")
	if discussion.CreatorID != userID.(int64) && !(isAdmin != nil && isAdmin.(bool)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this discussion"})
		return
	}

	// データベースから削除
	err = h.discussionRepo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Discussion deleted successfully"})
}

// @Summary Discussionのステータス変更
// @Description 指定されたIDのDiscussionのステータスを変更します
// @Tags discussions
// @Accept json
// @Produce json
// @Param id path int true "Discussion ID"
// @Param status body map[string]string true "ステータス情報"
// @Success 200 {object} models.Discussion
// @Router /api/v1/discussions/{id}/status [patch]
func (h *DiscussionHandler) UpdateDiscussionStatus(c *gin.Context) {
	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion ID format"})
		return
	}

	// データベースから取得
	discussion, err := h.discussionRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if discussion == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Discussion not found"})
		return
	}

	// リクエストの解析
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ステータスの更新
	switch req.Status {
	case "open":
		discussion.Reopen()
	case "closed":
		discussion.Close()
	case "answered":
		discussion.MarkAsAnswered()
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be 'open', 'closed', or 'answered'"})
		return
	}

	// データベースに保存
	err = h.discussionRepo.Update(c.Request.Context(), discussion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, discussion)
}

// @Summary Discussionのドラフト状態変更
// @Description 指定されたIDのDiscussionのドラフト状態を変更します
// @Tags discussions
// @Accept json
// @Produce json
// @Param id path int true "Discussion ID"
// @Param draft body map[string]bool true "ドラフト情報"
// @Success 200 {object} models.Discussion
// @Router /api/v1/discussions/{id}/draft [patch]
func (h *DiscussionHandler) UpdateDiscussionDraftStatus(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion ID format"})
		return
	}

	// データベースから取得
	discussion, err := h.discussionRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if discussion == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Discussion not found"})
		return
	}

	// 作成者のみ編集可能（または管理者）
	isAdmin, _ := c.Get("is_admin")
	if discussion.CreatorID != userID.(int64) && !(isAdmin != nil && isAdmin.(bool)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this discussion"})
		return
	}

	// リクエストの解析
	var req struct {
		IsDraft bool `json:"is_draft"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ドラフト状態の更新
	discussion.SetDraft(req.IsDraft)

	// データベースに保存
	err = h.discussionRepo.Update(c.Request.Context(), discussion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, discussion)
}

// @Summary Discussionの検索
// @Description Discussionを検索します
// @Tags discussions
// @Accept json
// @Produce json
// @Param q query string true "検索クエリ"
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/discussions/search [get]
func (h *DiscussionHandler) SearchDiscussions(c *gin.Context) {
	// クエリパラメータの取得
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// データベースから検索
	discussions, total, err := h.discussionRepo.Search(c.Request.Context(), query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"discussions": discussions,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"query":       query,
	})
}

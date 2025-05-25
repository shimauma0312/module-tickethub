package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
)

// MilestoneHandler はMilestone関連のハンドラーを管理する構造体
type MilestoneHandler struct {
	milestoneRepo repositories.MilestoneRepository
}

// NewMilestoneHandler は新しいMilestoneHandlerを作成します
func NewMilestoneHandler(milestoneRepo repositories.MilestoneRepository) *MilestoneHandler {
	return &MilestoneHandler{
		milestoneRepo: milestoneRepo,
	}
}

// MilestoneRequest はMilestone作成・更新リクエストのデータ構造
type MilestoneRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"` // ISO 8601形式 (YYYY-MM-DD)
}

// @Summary マイルストーン一覧の取得
// @Description マイルストーンの一覧を取得します
// @Tags milestones
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(10)
// @Param status query string false "ステータス (open/closed)" default("open")
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/milestones [get]
func (h *MilestoneHandler) ListMilestones(c *gin.Context) {
	// クエリパラメータの取得
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.DefaultQuery("status", "open")

	// フィルタの作成
	filter := map[string]interface{}{
		"status": status,
	}

	// データベースから取得
	milestones, total, err := h.milestoneRepo.List(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"milestones": milestones,
		"total":      total,
		"page":       page,
		"limit":      limit,
	})
}

// @Summary マイルストーンの取得
// @Description 指定されたIDのマイルストーンを取得します
// @Tags milestones
// @Accept json
// @Produce json
// @Param id path int true "マイルストーンID"
// @Success 200 {object} models.Milestone
// @Router /api/v1/milestones/{id} [get]
func (h *MilestoneHandler) GetMilestone(c *gin.Context) {
	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID format"})
		return
	}

	// データベースから取得
	milestone, err := h.milestoneRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if milestone == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found"})
		return
	}

	c.JSON(http.StatusOK, milestone)
}

// @Summary マイルストーンの作成
// @Description 新しいマイルストーンを作成します
// @Tags milestones
// @Accept json
// @Produce json
// @Param milestone body MilestoneRequest true "マイルストーン情報"
// @Success 201 {object} models.Milestone
// @Router /api/v1/milestones [post]
func (h *MilestoneHandler) CreateMilestone(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// リクエストの解析
	var req MilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 期日の解析
	var dueDate time.Time
	if req.DueDate != "" {
		var err error
		dueDate, err = time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date format. Use YYYY-MM-DD"})
			return
		}
	}

	// マイルストーンの作成
	milestone := models.NewMilestone(req.Title, req.Description, dueDate, userID.(int64))

	// 検証
	if !milestone.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone data"})
		return
	}

	// データベースに保存
	err := h.milestoneRepo.Create(c.Request.Context(), milestone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, milestone)
}

// @Summary マイルストーンの更新
// @Description 指定されたIDのマイルストーンを更新します
// @Tags milestones
// @Accept json
// @Produce json
// @Param id path int true "マイルストーンID"
// @Param milestone body MilestoneRequest true "マイルストーン情報"
// @Success 200 {object} models.Milestone
// @Router /api/v1/milestones/{id} [put]
func (h *MilestoneHandler) UpdateMilestone(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID format"})
		return
	}

	// データベースから取得
	milestone, err := h.milestoneRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if milestone == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found"})
		return
	}

	// 作成者のみ編集可能（または管理者）
	isAdmin, _ := c.Get("is_admin")
	if milestone.CreatorID != userID.(int64) && !(isAdmin != nil && isAdmin.(bool)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this milestone"})
		return
	}

	// リクエストの解析
	var req MilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 期日の解析
	var dueDate time.Time
	if req.DueDate != "" {
		var err error
		dueDate, err = time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date format. Use YYYY-MM-DD"})
			return
		}
	}

	// マイルストーンの更新
	milestone.Update(req.Title, req.Description, dueDate)

	// 検証
	if !milestone.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone data"})
		return
	}

	// データベースに保存
	err = h.milestoneRepo.Update(c.Request.Context(), milestone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, milestone)
}

// @Summary マイルストーンの削除
// @Description 指定されたIDのマイルストーンを削除します
// @Tags milestones
// @Accept json
// @Produce json
// @Param id path int true "マイルストーンID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/milestones/{id} [delete]
func (h *MilestoneHandler) DeleteMilestone(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID format"})
		return
	}

	// データベースから取得
	milestone, err := h.milestoneRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if milestone == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found"})
		return
	}

	// 作成者のみ削除可能（または管理者）
	isAdmin, _ := c.Get("is_admin")
	if milestone.CreatorID != userID.(int64) && !(isAdmin != nil && isAdmin.(bool)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this milestone"})
		return
	}

	// データベースから削除
	err = h.milestoneRepo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Milestone deleted successfully"})
}

// @Summary マイルストーンのステータス変更
// @Description 指定されたIDのマイルストーンのステータスを変更します
// @Tags milestones
// @Accept json
// @Produce json
// @Param id path int true "マイルストーンID"
// @Param status body map[string]string true "ステータス情報"
// @Success 200 {object} models.Milestone
// @Router /api/v1/milestones/{id}/status [patch]
func (h *MilestoneHandler) UpdateMilestoneStatus(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID format"})
		return
	}

	// データベースから取得
	milestone, err := h.milestoneRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if milestone == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found"})
		return
	}

	// 作成者のみ編集可能（または管理者）
	isAdmin, _ := c.Get("is_admin")
	if milestone.CreatorID != userID.(int64) && !(isAdmin != nil && isAdmin.(bool)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this milestone"})
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
		milestone.Reopen()
	case "closed":
		milestone.Close()
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be 'open' or 'closed'"})
		return
	}

	// データベースに保存
	err = h.milestoneRepo.Update(c.Request.Context(), milestone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, milestone)
}

package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
)

// LabelHandler はLabel関連のハンドラーを管理する構造体
type LabelHandler struct {
	labelRepo repositories.LabelRepository
}

// NewLabelHandler は新しいLabelHandlerを作成します
func NewLabelHandler(labelRepo repositories.LabelRepository) *LabelHandler {
	return &LabelHandler{
		labelRepo: labelRepo,
	}
}

// LabelRequest はLabel作成・更新リクエストのデータ構造
type LabelRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color" binding:"required"`
	Type        string `json:"type" binding:"required"` // issue/discussion/both
}

// @Summary ラベル一覧の取得
// @Description ラベルの一覧を取得します
// @Tags labels
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(50)
// @Param type query string false "ラベルタイプ (issue/discussion/both)"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/labels [get]
func (h *LabelHandler) ListLabels(c *gin.Context) {
	// クエリパラメータの取得
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	labelType := c.Query("type")

	// フィルタの作成
	filter := map[string]interface{}{}
	if labelType != "" {
		filter["type"] = labelType
	}

	// データベースから取得
	labels, total, err := h.labelRepo.List(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"labels": labels,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// @Summary ラベルの取得
// @Description 指定されたIDのラベルを取得します
// @Tags labels
// @Accept json
// @Produce json
// @Param id path int true "ラベルID"
// @Success 200 {object} models.Label
// @Router /api/v1/labels/{id} [get]
func (h *LabelHandler) GetLabel(c *gin.Context) {
	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID format"})
		return
	}

	// データベースから取得
	label, err := h.labelRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if label == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
		return
	}

	c.JSON(http.StatusOK, label)
}

// @Summary ラベルの作成
// @Description 新しいラベルを作成します
// @Tags labels
// @Accept json
// @Produce json
// @Param label body LabelRequest true "ラベル情報"
// @Success 201 {object} models.Label
// @Router /api/v1/labels [post]
func (h *LabelHandler) CreateLabel(c *gin.Context) {
	// 管理者権限の確認
	isAdmin, exists := c.Get("is_admin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
		return
	}

	// リクエストの解析
	var req LabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 同名ラベルの存在確認
	existingLabel, err := h.labelRepo.GetByName(c.Request.Context(), req.Name, req.Type)
	if err != nil && err.Error() != "label not found" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingLabel != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Label with this name already exists"})
		return
	}

	// ラベルの作成
	label := models.NewLabel(req.Name, req.Description, req.Color, req.Type)

	// 検証
	if !label.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label data"})
		return
	}

	// データベースに保存
	err = h.labelRepo.Create(c.Request.Context(), label)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, label)
}

// @Summary ラベルの更新
// @Description 指定されたIDのラベルを更新します
// @Tags labels
// @Accept json
// @Produce json
// @Param id path int true "ラベルID"
// @Param label body LabelRequest true "ラベル情報"
// @Success 200 {object} models.Label
// @Router /api/v1/labels/{id} [put]
func (h *LabelHandler) UpdateLabel(c *gin.Context) {
	// 管理者権限の確認
	isAdmin, exists := c.Get("is_admin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID format"})
		return
	}

	// データベースから取得
	label, err := h.labelRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if label == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
		return
	}

	// リクエストの解析
	var req LabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 同名ラベルの存在確認（更新対象以外）
	if req.Name != label.Name {
		existingLabel, err := h.labelRepo.GetByName(c.Request.Context(), req.Name, req.Type)
		if err != nil && err.Error() != "label not found" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if existingLabel != nil && existingLabel.ID != id {
			c.JSON(http.StatusConflict, gin.H{"error": "Label with this name already exists"})
			return
		}
	}

	// ラベルの更新
	label.Update(req.Name, req.Description, req.Color, req.Type)

	// 検証
	if !label.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label data"})
		return
	}

	// データベースに保存
	err = h.labelRepo.Update(c.Request.Context(), label)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, label)
}

// @Summary ラベルの削除
// @Description 指定されたIDのラベルを削除します
// @Tags labels
// @Accept json
// @Produce json
// @Param id path int true "ラベルID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/labels/{id} [delete]
func (h *LabelHandler) DeleteLabel(c *gin.Context) {
	// 管理者権限の確認
	isAdmin, exists := c.Get("is_admin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID format"})
		return
	}

	// データベースから取得
	label, err := h.labelRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if label == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
		return
	}

	// データベースから削除
	err = h.labelRepo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Label deleted successfully"})
}

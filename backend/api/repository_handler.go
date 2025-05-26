package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
	"github.com/shimauma0312/module-tickethub/backend/services"
)

// RepositoryHandler はリポジトリ管理のAPIハンドラー
type RepositoryHandler struct {
	repoRepo        repositories.RepositoryRepository
	activityService *services.ActivityLogService
}

// NewRepositoryHandler は新しいRepositoryHandlerを作成します
func NewRepositoryHandler(
	repoRepo repositories.RepositoryRepository,
	activityService *services.ActivityLogService,
) *RepositoryHandler {
	return &RepositoryHandler{
		repoRepo:        repoRepo,
		activityService: activityService,
	}
}

// GetRepositories はリポジトリ一覧を取得します
// @Summary リポジトリ一覧取得
// @Description 管理者がシステム内のリポジトリ一覧を取得します
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(20)
// @Param name query string false "リポジトリ名で絞り込み"
// @Param type query string false "リポジトリタイプで絞り込み"
// @Param is_archived query bool false "アーカイブ状態で絞り込み"
// @Success 200 {object} map[string]interface{} "リポジトリ一覧"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/repositories [get]
// @Security BearerAuth
func (h *RepositoryHandler) GetRepositories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filter["name"] = name
	}
	if repoType := c.Query("type"); repoType != "" {
		filter["type"] = repoType
	}
	if isArchivedStr := c.Query("is_archived"); isArchivedStr != "" {
		isArchived, err := strconv.ParseBool(isArchivedStr)
		if err == nil {
			filter["is_archived"] = isArchived
		}
	}

	repositories, total, err := h.repoRepo.List(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get repositories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"repositories": repositories,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}

// GetRepository はリポジトリ詳細を取得します
// @Summary リポジトリ詳細取得
// @Description 管理者がリポジトリの詳細情報を取得します
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "リポジトリID"
// @Success 200 {object} models.Repository "リポジトリ詳細"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 404 {object} map[string]string "リポジトリが見つからない"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/repositories/{id} [get]
// @Security BearerAuth
func (h *RepositoryHandler) GetRepository(c *gin.Context) {
	repoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	repo, err := h.repoRepo.GetByID(c.Request.Context(), repoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	c.JSON(http.StatusOK, repo)
}

// CreateRepository は新しいリポジトリを作成します
// @Summary リポジトリ作成
// @Description 管理者が新しいリポジトリを作成します
// @Tags admin
// @Accept json
// @Produce json
// @Param repository body object true "リポジトリ情報"
// @Success 201 {object} models.Repository "作成されたリポジトリ"
// @Failure 400 {object} map[string]string "リクエストエラー"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/repositories [post]
// @Security BearerAuth
func (h *RepositoryHandler) CreateRepository(c *gin.Context) {
	var requestData struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Type        string `json:"type" binding:"required,oneof=public private internal"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 作成者（管理者）のIDを取得
	ownerID := c.GetInt64("user_id")

	// リポジトリタイプを変換
	var repoType models.RepositoryType
	switch requestData.Type {
	case "public":
		repoType = models.PublicRepo
	case "private":
		repoType = models.PrivateRepo
	case "internal":
		repoType = models.InternalRepo
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository type"})
		return
	}

	// リポジトリの作成
	repo := models.NewRepository(requestData.Name, requestData.Description, repoType, ownerID)
	if err := h.repoRepo.Create(c.Request.Context(), repo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create repository"})
		return
	}

	// アクティビティログに記録
	username := c.GetString("username")
	h.activityService.LogActivity(
		c.Request.Context(),
		ownerID,
		username,
		models.ActionRepositoryCreated,
		models.ResourceRepository,
		repo.ID,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		map[string]interface{}{
			"name":        repo.Name,
			"description": repo.Description,
			"type":        repo.Type,
		},
	)

	c.JSON(http.StatusCreated, repo)
}

// UpdateRepository はリポジトリを更新します
// @Summary リポジトリ更新
// @Description 管理者がリポジトリ情報を更新します
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "リポジトリID"
// @Param repository body object true "更新するリポジトリ情報"
// @Success 200 {object} models.Repository "更新されたリポジトリ"
// @Failure 400 {object} map[string]string "リクエストエラー"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 404 {object} map[string]string "リポジトリが見つからない"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/repositories/{id} [put]
// @Security BearerAuth
func (h *RepositoryHandler) UpdateRepository(c *gin.Context) {
	repoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	repo, err := h.repoRepo.GetByID(c.Request.Context(), repoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	var requestData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type" binding:"omitempty,oneof=public private internal"`
		IsArchived  *bool  `json:"is_archived"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// リポジトリタイプを変換
	var repoType models.RepositoryType
	if requestData.Type != "" {
		switch requestData.Type {
		case "public":
			repoType = models.PublicRepo
		case "private":
			repoType = models.PrivateRepo
		case "internal":
			repoType = models.InternalRepo
		}
	}

	// リポジトリ情報を更新
	repo.Update(requestData.Name, requestData.Description, repoType)

	// アーカイブ状態を更新
	if requestData.IsArchived != nil {
		if *requestData.IsArchived {
			repo.Archive()
		} else {
			repo.Unarchive()
		}
	}

	if err := h.repoRepo.Update(c.Request.Context(), repo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update repository"})
		return
	}

	// アクティビティログに記録
	userID := c.GetInt64("user_id")
	username := c.GetString("username")
	h.activityService.LogActivity(
		c.Request.Context(),
		userID,
		username,
		models.ActionRepositoryUpdated,
		models.ResourceRepository,
		repoID,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		map[string]interface{}{
			"name":        requestData.Name,
			"description": requestData.Description,
			"type":        requestData.Type,
			"is_archived": requestData.IsArchived,
		},
	)

	c.JSON(http.StatusOK, repo)
}

// DeleteRepository はリポジトリを削除します
// @Summary リポジトリ削除
// @Description 管理者がリポジトリを削除します
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "リポジトリID"
// @Success 200 {object} map[string]string "削除成功メッセージ"
// @Failure 400 {object} map[string]string "リクエストエラー"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 404 {object} map[string]string "リポジトリが見つからない"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/repositories/{id} [delete]
// @Security BearerAuth
func (h *RepositoryHandler) DeleteRepository(c *gin.Context) {
	repoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	// リポジトリが存在するか確認
	repo, err := h.repoRepo.GetByID(c.Request.Context(), repoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	// リポジトリを削除
	if err := h.repoRepo.Delete(c.Request.Context(), repoID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete repository"})
		return
	}

	// アクティビティログに記録
	userID := c.GetInt64("user_id")
	username := c.GetString("username")
	h.activityService.LogActivity(
		c.Request.Context(),
		userID,
		username,
		models.ActionRepositoryDeleted,
		models.ResourceRepository,
		repoID,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		map[string]interface{}{
			"name": repo.Name,
		},
	)

	c.JSON(http.StatusOK, gin.H{"message": "Repository deleted successfully"})
}

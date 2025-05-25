package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
	"github.com/shimauma0312/module-tickethub/backend/services"
)

// AdminHandler は管理者機能のAPIハンドラー
type AdminHandler struct {
	userRepo        repositories.UserRepository
	systemRepo      repositories.SystemSettingsRepository
	activityService *services.ActivityLogService
	backupService   *services.BackupService
	metricsService  *services.SystemMetricsService
}

// NewAdminHandler は新しいAdminHandlerを作成します
func NewAdminHandler(
	userRepo repositories.UserRepository,
	systemRepo repositories.SystemSettingsRepository,
	activityService *services.ActivityLogService,
	backupService *services.BackupService,
	metricsService *services.SystemMetricsService,
) *AdminHandler {
	return &AdminHandler{
		userRepo:        userRepo,
		systemRepo:      systemRepo,
		activityService: activityService,
		backupService:   backupService,
		metricsService:  metricsService,
	}
}

// GetUsers はユーザー一覧を取得します
// @Summary ユーザー一覧取得
// @Description 管理者がシステム内のユーザー一覧を取得します
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(20)
// @Param active query bool false "アクティブユーザーのみ"
// @Success 200 {object} map[string]interface{} "ユーザー一覧"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/users [get]
// @Security BearerAuth
func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := make(map[string]interface{})
	if activeStr := c.Query("active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			filter["is_active"] = active
		}
	}

	users, total, err := h.userRepo.List(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// UpdateUser はユーザー情報を更新します
// @Summary ユーザー情報更新
// @Description 管理者がユーザー情報を更新します
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "ユーザーID"
// @Param user body object true "更新するユーザー情報"
// @Success 200 {object} models.User "更新されたユーザー情報"
// @Failure 400 {object} map[string]string "リクエストエラー"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 404 {object} map[string]string "ユーザーが見つからない"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/users/{id} [put]
// @Security BearerAuth
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData struct {
		FullName string `json:"full_name"`
		IsAdmin  *bool  `json:"is_admin"`
		IsActive *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// ユーザー情報を更新
	if updateData.FullName != "" {
		user.FullName = updateData.FullName
	}
	if updateData.IsAdmin != nil {
		user.SetAdmin(*updateData.IsAdmin)
	}
	if updateData.IsActive != nil {
		if *updateData.IsActive {
			user.Activate()
		} else {
			user.Deactivate()
		}
	}

	if err := h.userRepo.Update(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// アクティビティログに記録
	currentUserID := c.GetInt64("user_id")
	currentUsername := c.GetString("username")
	h.activityService.LogActivity(
		c.Request.Context(),
		currentUserID,
		currentUsername,
		models.ActionUserUpdated,
		models.ResourceUser,
		userID,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		map[string]interface{}{
			"updated_fields": updateData,
		},
	)

	c.JSON(http.StatusOK, user)
}

// GetSystemSettings はシステム設定を取得します
// @Summary システム設定取得
// @Description 管理者がシステム設定を取得します
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} models.SystemSettings "システム設定"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/settings [get]
// @Security BearerAuth
func (h *AdminHandler) GetSystemSettings(c *gin.Context) {
	settings, err := h.systemRepo.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get system settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSystemSettings はシステム設定を更新します
// @Summary システム設定更新
// @Description 管理者がシステム設定を更新します
// @Tags admin
// @Accept json
// @Produce json
// @Param settings body map[string]interface{} true "更新する設定"
// @Success 200 {object} models.SystemSettings "更新されたシステム設定"
// @Failure 400 {object} map[string]string "リクエストエラー"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/settings [put]
// @Security BearerAuth
func (h *AdminHandler) UpdateSystemSettings(c *gin.Context) {
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	settings, err := h.systemRepo.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current settings"})
		return
	}

	settings.Update(updateData)

	if err := h.systemRepo.CreateOrUpdate(c.Request.Context(), settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	// アクティビティログに記録
	userID := c.GetInt64("user_id")
	username := c.GetString("username")
	h.activityService.LogActivity(
		c.Request.Context(),
		userID,
		username,
		models.ActionSystemSettingsUpdated,
		models.ResourceSystem,
		0,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		updateData,
	)

	c.JSON(http.StatusOK, settings)
}

// GetActivityLogs はアクティビティログを取得します
// @Summary アクティビティログ取得
// @Description 管理者がシステムのアクティビティログを取得します
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(20)
// @Param user_id query int false "ユーザーID"
// @Param action query string false "アクション"
// @Param resource_type query string false "リソースタイプ"
// @Success 200 {object} map[string]interface{} "アクティビティログ一覧"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/activity-logs [get]
// @Security BearerAuth
func (h *AdminHandler) GetActivityLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := make(map[string]interface{})
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			filter["user_id"] = userID
		}
	}
	if action := c.Query("action"); action != "" {
		filter["action"] = action
	}
	if resourceType := c.Query("resource_type"); resourceType != "" {
		filter["resource_type"] = resourceType
	}

	logs, total, err := h.activityService.GetActivityLogs(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activity logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetSystemMetrics はシステムメトリクスを取得します
// @Summary システムメトリクス取得
// @Description 管理者がシステムのメトリクス情報を取得します
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} models.SystemMetrics "システムメトリクス"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/metrics [get]
// @Security BearerAuth
func (h *AdminHandler) GetSystemMetrics(c *gin.Context) {
	metrics, err := h.metricsService.GetSystemMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get system metrics"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// CreateBackup はデータベースバックアップを作成します
// @Summary バックアップ作成
// @Description 管理者がデータベースバックアップを作成します
// @Tags admin
// @Accept json
// @Produce json
// @Param backup body object true "バックアップ情報"
// @Success 200 {object} models.BackupInfo "作成されたバックアップ情報"
// @Failure 400 {object} map[string]string "リクエストエラー"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/backups [post]
// @Security BearerAuth
func (h *AdminHandler) CreateBackup(c *gin.Context) {
	var requestData struct {
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := c.GetInt64("user_id")
	backup, err := h.backupService.CreateBackup(c.Request.Context(), userID, requestData.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create backup"})
		return
	}

	// アクティビティログに記録
	username := c.GetString("username")
	h.activityService.LogActivity(
		c.Request.Context(),
		userID,
		username,
		models.ActionSystemBackupCreated,
		models.ResourceSystem,
		backup.ID,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		map[string]interface{}{
			"description": requestData.Description,
		},
	)

	c.JSON(http.StatusOK, backup)
}

// GetBackups はバックアップ一覧を取得します
// @Summary バックアップ一覧取得
// @Description 管理者がバックアップ一覧を取得します
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(20)
// @Success 200 {object} map[string]interface{} "バックアップ一覧"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/backups [get]
// @Security BearerAuth
func (h *AdminHandler) GetBackups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	backups, total, err := h.backupService.GetBackups(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get backups"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"backups": backups,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// RestoreBackup はバックアップからデータベースを復元します
// @Summary バックアップ復元
// @Description 管理者がバックアップからデータベースを復元します
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "バックアップID"
// @Success 200 {object} map[string]string "復元成功メッセージ"
// @Failure 400 {object} map[string]string "リクエストエラー"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 404 {object} map[string]string "バックアップが見つからない"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/backups/{id}/restore [post]
// @Security BearerAuth
func (h *AdminHandler) RestoreBackup(c *gin.Context) {
	backupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
		return
	}

	if err := h.backupService.RestoreBackup(c.Request.Context(), backupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore backup"})
		return
	}

	// アクティビティログに記録
	userID := c.GetInt64("user_id")
	username := c.GetString("username")
	h.activityService.LogActivity(
		c.Request.Context(),
		userID,
		username,
		models.ActionSystemBackupRestored,
		models.ResourceSystem,
		backupID,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		nil,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Backup restored successfully"})
}

// DeleteBackup はバックアップを削除します
// @Summary バックアップ削除
// @Description 管理者がバックアップを削除します
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "バックアップID"
// @Success 200 {object} map[string]string "削除成功メッセージ"
// @Failure 400 {object} map[string]string "リクエストエラー"
// @Failure 401 {object} map[string]string "認証エラー"
// @Failure 403 {object} map[string]string "権限エラー"
// @Failure 404 {object} map[string]string "バックアップが見つからない"
// @Failure 500 {object} map[string]string "サーバーエラー"
// @Router /api/admin/backups/{id} [delete]
// @Security BearerAuth
func (h *AdminHandler) DeleteBackup(c *gin.Context) {
	backupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
		return
	}

	if err := h.backupService.DeleteBackup(c.Request.Context(), backupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete backup"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup deleted successfully"})
}

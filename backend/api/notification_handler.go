package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/services"
)

// NotificationHandler は通知関連のAPIハンドラ
type NotificationHandler struct {
	notificationService *services.NotificationService
}

// NewNotificationHandler は新しいNotificationHandlerを作成します
func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// RegisterRoutes はルーティングを登録します
func (h *NotificationHandler) RegisterRoutes(router *gin.RouterGroup) {
	notifications := router.Group("/notifications")
	{
		notifications.GET("", h.GetNotifications)
		notifications.GET("/unread-count", h.GetUnreadCount)
		notifications.PUT("/:id/read", h.MarkAsRead)
		notifications.PUT("/read-all", h.MarkAllAsRead)
		notifications.POST("/subscribe", h.Subscribe)
		notifications.DELETE("/unsubscribe", h.Unsubscribe)
		notifications.GET("/vapid-public-key", h.GetVAPIDPublicKey)
	}

	settings := router.Group("/settings/notifications")
	{
		settings.GET("", h.GetNotificationSettings)
		settings.PUT("", h.UpdateNotificationSettings)
	}
}

// GetNotifications はユーザーの通知一覧を取得します
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// クエリパラメータの取得
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	isReadStr := c.Query("is_read")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	var isRead *bool
	if isReadStr != "" {
		isReadBool := isReadStr == "true"
		isRead = &isReadBool
	}

	// 通知一覧を取得
	notifications, total, err := h.notificationService.GetNotifications(c.Request.Context(), userID, isRead, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         total,
		"page":          page,
		"limit":         limit,
	})
}

// GetUnreadCount はユーザーの未読通知数を取得します
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 未読通知のみを取得
	isRead := false
	notifications, _, err := h.notificationService.GetNotifications(c.Request.Context(), userID, &isRead, 1, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get unread count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unread_count": len(notifications),
	})
}

// MarkAsRead は通知を既読状態に更新します
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// パスパラメータからID取得
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	// 既読に更新
	err = h.notificationService.MarkAsRead(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// MarkAllAsRead はユーザーの全通知を既読状態に更新します
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 全て既読に更新
	err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark all as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// SubscriptionRequest はプッシュ通知サブスクリプションリクエスト
type SubscriptionRequest struct {
	Subscription string `json:"subscription" binding:"required"`
}

// Subscribe はブラウザのプッシュ通知サブスクリプションを登録します
func (h *NotificationHandler) Subscribe(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// サブスクリプションを登録
	err := h.notificationService.AddPushSubscription(c.Request.Context(), userID, req.Subscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to subscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UnsubscribeRequest はプッシュ通知サブスクリプション解除リクエスト
type UnsubscribeRequest struct {
	Endpoint string `json:"endpoint" binding:"required"`
}

// Unsubscribe はブラウザのプッシュ通知サブスクリプションを解除します
func (h *NotificationHandler) Unsubscribe(c *gin.Context) {
	var req UnsubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// サブスクリプションを解除
	err := h.notificationService.RemovePushSubscription(c.Request.Context(), req.Endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetVAPIDPublicKey はVAPID公開鍵を取得します
func (h *NotificationHandler) GetVAPIDPublicKey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"vapid_public_key": h.notificationService.GetVAPIDPublicKey(),
	})
}

// GetNotificationSettings はユーザーの通知設定を取得します
func (h *NotificationHandler) GetNotificationSettings(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// ユーザー設定を取得
	userSettings, err := h.notificationService.GetUserSettings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email_notification": userSettings.EmailNotification,
		"push_notification":  userSettings.PushNotification,
		"notification_types": userSettings.NotificationTypes,
	})
}

// UpdateNotificationSettingsRequest は通知設定更新リクエスト
type UpdateNotificationSettingsRequest struct {
	EmailNotification *bool   `json:"email_notification"`
	PushNotification  *bool   `json:"push_notification"`
	NotificationTypes *string `json:"notification_types"`
}

// UpdateNotificationSettings はユーザーの通知設定を更新します
func (h *NotificationHandler) UpdateNotificationSettings(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateNotificationSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 通知設定を更新
	err := h.notificationService.UpdateUserNotificationSettings(
		c.Request.Context(),
		userID,
		req.EmailNotification,
		req.PushNotification,
		req.NotificationTypes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// getUserIDFromContext はコンテキストからユーザーIDを取得するヘルパー関数
func getUserIDFromContext(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	id, ok := userID.(int64)
	if !ok {
		return 0
	}

	return id
}

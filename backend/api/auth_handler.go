package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/services"
)

// AuthHandler は認証関連のハンドラー
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler は新しいAuthHandlerを作成します
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRequest はユーザー登録リクエストのデータ構造
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"omitempty,max=100"`
}

// LoginRequest はログインリクエストのデータ構造
type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// PasswordChangeRequest はパスワード変更リクエストのデータ構造
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// PasswordResetRequest はパスワードリセットリクエストのデータ構造
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetCompleteRequest はパスワードリセット完了リクエストのデータ構造
type PasswordResetCompleteRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// TokenResponse はトークンレスポンスのデータ構造
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Register はユーザー登録ハンドラー
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Username, req.Email, req.Password, req.FullName)
	if err != nil {
		// エラーメッセージから具体的なエラーの種類を判断
		if strings.Contains(err.Error(), "username") && strings.Contains(err.Error(), "already taken") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
			return
		}
		if strings.Contains(err.Error(), "email") && strings.Contains(err.Error(), "already registered") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email is already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login はログインハンドラー
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	user, accessToken, refreshToken, err := h.authService.Login(
		c.Request.Context(),
		req.UsernameOrEmail,
		req.Password,
		userAgent,
		ipAddress,
	)

	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	// JWTトークンをCookieにも設定（クロスサイトリクエストから保護するため）
	c.SetCookie(
		"access_token",
		accessToken,
		int(30*time.Minute.Seconds()), // 30分
		"/",
		"",
		false, // 本番環境ではtrueに
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
		"token": TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int(30 * time.Minute.Seconds()),
			TokenType:    "Bearer",
		},
	})
}

// Logout はログアウトハンドラー
func (h *AuthHandler) Logout(c *gin.Context) {
	// ヘッダーからトークンを取得
	token := extractToken(c)

	// ユーザーIDはAuthMiddlewareで設定されることを前提
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	// Cookieも削除
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// LogoutAll は全セッションからログアウトするハンドラー
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	// ユーザーIDはAuthMiddlewareで設定されることを前提
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.authService.LogoutAll(c.Request.Context(), userID.(int64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout from all sessions"})
		return
	}

	// Cookieも削除
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out from all sessions"})
}

// RefreshToken はトークンを更新するハンドラー
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// リフレッシュトークンを取得
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	// 新しいアクセストークンを生成
	newAccessToken, newRefreshToken, err := h.authService.RefreshToken(
		c.Request.Context(),
		refreshToken,
		userAgent,
		ipAddress,
	)

	if err != nil {
		if err == services.ErrTokenInvalid || err == services.ErrTokenExpired || err == services.ErrTokenRevoked {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}

	// Cookieも更新
	c.SetCookie(
		"access_token",
		newAccessToken,
		int(30*time.Minute.Seconds()), // 30分
		"/",
		"",
		false, // 本番環境ではtrueに
		true,  // HttpOnly
	)

	// 新しいリフレッシュトークンをCookieに設定
	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		int((7 * 24 * time.Hour).Seconds()), // 7 days
		"/api/auth",
		"",
		true,  // Secure: HTTPSでのみ送信
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

// ChangePassword はパスワード変更ハンドラー
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req PasswordChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザーIDはAuthMiddlewareで設定されることを前提
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.authService.ChangePassword(
		c.Request.Context(),
		userID.(int64),
		req.CurrentPassword,
		req.NewPassword,
	); err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// InitiatePasswordReset はパスワードリセットを開始するハンドラー
func (h *AuthHandler) InitiatePasswordReset(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.InitiatePasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		if err == services.ErrUserNotFound {
			// セキュリティのため、ユーザーが存在しない場合でも成功したように見せる
			c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link has been sent"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate password reset"})
		return
	}

	// 本番環境では、ここでメール送信サービスを呼び出す
	// 開発環境ではトークンを直接返す
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset initiated",
		"token":   token, // 本番環境では削除
	})
}

// CompletePasswordReset はパスワードリセットを完了するハンドラー
func (h *AuthHandler) CompletePasswordReset(c *gin.Context) {
	var req PasswordResetCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.CompletePasswordReset(
		c.Request.Context(),
		req.Token,
		req.NewPassword,
	); err != nil {
		if err == services.ErrPasswordResetTokenInvalid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}

// ValidatePasswordResetToken はパスワードリセットトークンを検証するハンドラー
func (h *AuthHandler) ValidatePasswordResetToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	user, err := h.authService.ValidatePasswordResetToken(c.Request.Context(), token)
	if err != nil {
		if err == services.ErrPasswordResetTokenInvalid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token is valid",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// extractToken はリクエストからトークンを抽出します
func extractToken(c *gin.Context) string {
	// Authorization ヘッダーからトークンを取得
	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) > 7 && strings.ToUpper(bearerToken[0:7]) == "BEARER " {
		return bearerToken[7:]
	}

	// Cookie からトークンを取得
	token, _ := c.Cookie("access_token")
	if token != "" {
		return token
	}

	// クエリパラメータからトークンを取得
	return c.Query("token")
}

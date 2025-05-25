package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/services"
)

// AuthMiddleware は認証ミドルウェア
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// トークンの取得
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		// トークンの検証
		user, claims, err := authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			var statusCode int
			var message string

			switch err {
			case services.ErrTokenInvalid:
				statusCode = http.StatusUnauthorized
				message = "Invalid token"
			case services.ErrTokenExpired:
				statusCode = http.StatusUnauthorized
				message = "Token has expired"
			case services.ErrTokenRevoked:
				statusCode = http.StatusUnauthorized
				message = "Token has been revoked"
			case services.ErrUserNotFound:
				statusCode = http.StatusUnauthorized
				message = "User not found"
			default:
				statusCode = http.StatusInternalServerError
				message = "Failed to validate token"
			}

			c.JSON(statusCode, gin.H{"error": message})
			c.Abort()
			return
		}

		// アクセストークンのみ許可
		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
			c.Abort()
			return
		}

		// ユーザー情報をコンテキストに設定
		c.Set("user_id", user.ID)
		c.Set("username", user.Username)
		c.Set("is_admin", user.IsAdmin)

		c.Next()
	}
}

// AdminMiddleware は管理者権限ミドルウェア
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CSRFMiddleware はCSRF対策ミドルウェア
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET, HEAD, OPTIONS, TRACE はCSRF保護不要
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodHead ||
			c.Request.Method == http.MethodOptions ||
			c.Request.Method == http.MethodTrace {
			c.Next()
			return
		}

		// CSRFトークンの検証
		csrfToken := c.GetHeader("X-CSRF-Token")
		if csrfToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token is required"})
			c.Abort()
			return
		}

		// Cookieからトークンを取得して比較（実際の実装ではSecureなトークン検証を行う）
		cookie, err := c.Cookie("csrf_token")
		if err != nil || cookie == "" || cookie != csrfToken {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateCSRFToken はCSRFトークンを生成するハンドラー
func GenerateCSRFToken(c *gin.Context) {
	// 実際の実装では安全なランダムトークンを生成する
	token := "csrf-token-" + c.ClientIP() + "-" + c.GetHeader("User-Agent")

	// Cookieに設定
	c.SetCookie(
		"csrf_token",
		token,
		3600, // 1時間
		"/",
		"",
		false, // 本番環境ではtrueに
		false, // JavaScriptからアクセス可能にする
	)

	c.JSON(http.StatusOK, gin.H{
		"csrf_token": token,
	})
}

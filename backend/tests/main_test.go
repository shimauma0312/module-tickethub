package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "TicketHub API Server is running",
			"version": "0.1.0",
		})
	})
	return r
}

func TestPingRoute(t *testing.T) {
	// テスト用のルーターを設定
	r := setupRouter()

	// テストリクエストを作成
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	// レスポンスレコーダーを作成
	w := httptest.NewRecorder()

	// リクエストを実行
	r.ServeHTTP(w, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TicketHub API Server is running")
}

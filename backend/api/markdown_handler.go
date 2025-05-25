package api

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownHandler はMarkdown関連のハンドラーを管理する構造体
type MarkdownHandler struct {
	md goldmark.Markdown
}

// NewMarkdownHandler は新しいMarkdownHandlerを作成します
func NewMarkdownHandler() *MarkdownHandler {
	// Goldmarkインスタンスの設定
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,           // GitHub Flavored Markdown
			extension.Table,         // テーブル
			extension.Strikethrough, // 取り消し線
			extension.Linkify,       // URL自動リンク
			extension.TaskList,      // タスクリスト
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // 見出しID自動生成
			parser.WithAttribute(),     // 属性
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // 改行をbrタグに変換
			html.WithXHTML(),     // XHTML準拠
			html.WithUnsafe(),    // HTMLタグを許可
		),
	)

	return &MarkdownHandler{
		md: md,
	}
}

// MarkdownRequest はMarkdownレンダリングリクエストのデータ構造
type MarkdownRequest struct {
	Text string `json:"text" binding:"required"`
}

// @Summary Markdownのレンダリング
// @Description MarkdownをHTMLに変換します
// @Tags markdown
// @Accept json
// @Produce json
// @Param markdown body MarkdownRequest true "Markdownテキスト"
// @Success 200 {object} map[string]string
// @Router /api/v1/markdown [post]
func (h *MarkdownHandler) RenderMarkdown(c *gin.Context) {
	// リクエストの解析
	var req MarkdownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Markdownの変換
	var buf bytes.Buffer
	if err := h.md.Convert([]byte(req.Text), &buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render markdown: " + err.Error()})
		return
	}

	// レスポンスの返却
	c.JSON(http.StatusOK, gin.H{
		"html": buf.String(),
	})
}

// @Summary MarkdownのプレビューHTML取得（生）
// @Description MarkdownをHTMLに変換し、生のHTMLを返します
// @Tags markdown
// @Accept json
// @Produce html
// @Param markdown body MarkdownRequest true "Markdownテキスト"
// @Success 200 {string} string "HTML"
// @Router /api/v1/markdown/raw [post]
func (h *MarkdownHandler) RenderRawMarkdown(c *gin.Context) {
	// リクエストの解析
	var req MarkdownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Markdownの変換
	var buf bytes.Buffer
	if err := h.md.Convert([]byte(req.Text), &buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render markdown: " + err.Error()})
		return
	}

	// レスポンスの返却
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, buf.String())
}

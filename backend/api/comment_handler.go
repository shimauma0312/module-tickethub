package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
)

// CommentHandler はComment関連のハンドラーを管理する構造体
type CommentHandler struct {
	commentRepo    repositories.CommentRepository
	issueRepo      repositories.IssueRepository
	discussionRepo repositories.DiscussionRepository
	reactionRepo   repositories.ReactionRepository
	userRepo       repositories.UserRepository
}

// NewCommentHandler は新しいCommentHandlerを作成します
func NewCommentHandler(
	commentRepo repositories.CommentRepository,
	issueRepo repositories.IssueRepository,
	discussionRepo repositories.DiscussionRepository,
	reactionRepo repositories.ReactionRepository,
	userRepo repositories.UserRepository,
) *CommentHandler {
	return &CommentHandler{
		commentRepo:    commentRepo,
		issueRepo:      issueRepo,
		discussionRepo: discussionRepo,
		reactionRepo:   reactionRepo,
		userRepo:       userRepo,
	}
}

// CommentRequest はComment作成・更新リクエストのデータ構造
type CommentRequest struct {
	Body string `json:"body" binding:"required"`
}

// ReplyCommentRequest は返信コメント作成リクエストのデータ構造
type ReplyCommentRequest struct {
	Body            string `json:"body" binding:"required"`
	ParentCommentID int64  `json:"parent_comment_id" binding:"required"`
}

// @Summary コメント一覧の取得
// @Description 指定されたIssueまたはDiscussionのコメント一覧を取得します
// @Tags comments
// @Accept json
// @Produce json
// @Param target_type path string true "ターゲットタイプ (issue/discussion)"
// @Param target_id path int true "ターゲットID"
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(50)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/{target_type}/{target_id}/comments [get]
func (h *CommentHandler) ListComments(c *gin.Context) {
	// パラメータの取得
	targetType := c.Param("target_type")
	if targetType != "issue" && targetType != "discussion" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target type. Must be 'issue' or 'discussion'"})
		return
	}

	targetID, err := strconv.ParseInt(c.Param("target_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// ターゲットの存在確認
	if targetType == "issue" {
		issue, err := h.issueRepo.GetByID(c.Request.Context(), targetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if issue == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
			return
		}
	} else {
		discussion, err := h.discussionRepo.GetByID(c.Request.Context(), targetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if discussion == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Discussion not found"})
			return
		}
	}

	// コメント一覧の取得
	comments, total, err := h.commentRepo.ListByTarget(c.Request.Context(), targetID, targetType, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// コメントと返信をグループ化（返信はTreeNodeに変換）
	commentTree := buildCommentTree(comments)

	c.JSON(http.StatusOK, gin.H{
		"comments": commentTree,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// CommentTreeNode はコメントツリーのノードを表す構造体
type CommentTreeNode struct {
	Comment *models.Comment    `json:"comment"`
	Replies []*CommentTreeNode `json:"replies,omitempty"`
}

// buildCommentTree はコメントとその返信をツリー構造に変換する
func buildCommentTree(comments []*models.Comment) []*CommentTreeNode {
	// コメントをIDでマッピング
	commentMap := make(map[int64]*CommentTreeNode)
	var rootComments []*CommentTreeNode

	// 最初にすべてのコメントをノードに変換
	for _, comment := range comments {
		node := &CommentTreeNode{
			Comment: comment,
			Replies: []*CommentTreeNode{},
		}
		commentMap[comment.ID] = node

		// 親コメントがなければルートコメントに追加
		if !comment.IsReply() {
			rootComments = append(rootComments, node)
		}
	}

	// 返信を親コメントに追加
	for _, comment := range comments {
		if comment.IsReply() && comment.ParentCommentID > 0 {
			// 親コメントが見つかった場合のみ追加
			if parentNode, ok := commentMap[comment.ParentCommentID]; ok {
				node := commentMap[comment.ID]
				parentNode.Replies = append(parentNode.Replies, node)
			}
		}
	}

	return rootComments
}

// @Summary コメントの取得
// @Description 指定されたIDのコメントを取得します
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "コメントID"
// @Success 200 {object} models.Comment
// @Router /api/v1/comments/{id} [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID format"})
		return
	}

	// データベースから取得
	comment, err := h.commentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if comment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// @Summary コメントの作成
// @Description 指定されたIssueまたはDiscussionに新しいコメントを作成します
// @Tags comments
// @Accept json
// @Produce json
// @Param target_type path string true "ターゲットタイプ (issue/discussion)"
// @Param target_id path int true "ターゲットID"
// @Param comment body CommentRequest true "コメント情報"
// @Success 201 {object} models.Comment
// @Router /api/v1/{target_type}/{target_id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// パラメータの取得
	targetType := c.Param("target_type")
	if targetType != "issue" && targetType != "discussion" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target type. Must be 'issue' or 'discussion'"})
		return
	}

	targetID, err := strconv.ParseInt(c.Param("target_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
		return
	}

	// ターゲットの存在確認
	if targetType == "issue" {
		issue, err := h.issueRepo.GetByID(c.Request.Context(), targetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if issue == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
			return
		}
	} else {
		discussion, err := h.discussionRepo.GetByID(c.Request.Context(), targetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if discussion == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Discussion not found"})
			return
		}
	}

	// リクエストの解析
	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// コメントの作成
	comment := models.NewComment(req.Body, userID.(int64), targetID, targetType)

	// データベースに保存
	err = h.commentRepo.Create(c.Request.Context(), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// @Summary 返信コメントの作成
// @Description 指定されたコメントに返信コメントを作成します
// @Tags comments
// @Accept json
// @Produce json
// @Param target_type path string true "ターゲットタイプ (issue/discussion)"
// @Param target_id path int true "ターゲットID"
// @Param comment body ReplyCommentRequest true "返信コメント情報"
// @Success 201 {object} models.Comment
// @Router /api/v1/{target_type}/{target_id}/comments/reply [post]
func (h *CommentHandler) CreateReplyComment(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// パラメータの取得
	targetType := c.Param("target_type")
	if targetType != "issue" && targetType != "discussion" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target type. Must be 'issue' or 'discussion'"})
		return
	}

	targetID, err := strconv.ParseInt(c.Param("target_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
		return
	}

	// リクエストの解析
	var req ReplyCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 親コメントの存在確認
	parentComment, err := h.commentRepo.GetByID(c.Request.Context(), req.ParentCommentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if parentComment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parent comment not found"})
		return
	}

	// 親コメントが同じターゲットに属しているか確認
	if parentComment.TargetID != targetID || parentComment.Type != targetType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parent comment does not belong to the specified target"})
		return
	}

	// 返信コメントの作成
	comment := models.NewReply(req.Body, userID.(int64), targetID, req.ParentCommentID, targetType)

	// データベースに保存
	err = h.commentRepo.Create(c.Request.Context(), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// @Summary コメントの更新
// @Description 指定されたIDのコメントを更新します
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "コメントID"
// @Param comment body CommentRequest true "コメント情報"
// @Success 200 {object} models.Comment
// @Router /api/v1/comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID format"})
		return
	}

	// データベースから取得
	comment, err := h.commentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if comment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// 作成者のみ編集可能（または管理者）
	isAdmin, _ := c.Get("is_admin")
	if comment.CreatorID != userID.(int64) && !(isAdmin != nil && isAdmin.(bool)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this comment"})
		return
	}

	// リクエストの解析
	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// コメントの更新
	comment.Edit(req.Body)

	// データベースに保存
	err = h.commentRepo.Update(c.Request.Context(), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// @Summary コメントの削除
// @Description 指定されたIDのコメントを削除します
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "コメントID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID format"})
		return
	}

	// データベースから取得
	comment, err := h.commentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if comment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// 作成者のみ削除可能（または管理者）
	isAdmin, _ := c.Get("is_admin")
	if comment.CreatorID != userID.(int64) && !(isAdmin != nil && isAdmin.(bool)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this comment"})
		return
	}

	// データベースから削除
	err = h.commentRepo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// @Summary 返信コメント一覧の取得
// @Description 指定されたコメントへの返信コメント一覧を取得します
// @Tags comments
// @Accept json
// @Produce json
// @Param comment_id path int true "親コメントID"
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/comments/{comment_id}/replies [get]
func (h *CommentHandler) ListReplies(c *gin.Context) {
	// IDの取得
	commentID, err := strconv.ParseInt(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID format"})
		return
	}

	// 親コメントの存在確認
	parentComment, err := h.commentRepo.GetByID(c.Request.Context(), commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if parentComment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parent comment not found"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// 返信コメント一覧の取得
	replies, total, err := h.commentRepo.ListReplies(c.Request.Context(), commentID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"replies": replies,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

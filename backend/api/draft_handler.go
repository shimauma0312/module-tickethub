package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
)

// DraftHandler はドラフト機能のハンドラーを管理する構造体
type DraftHandler struct {
	issueRepo      repositories.IssueRepository
	discussionRepo repositories.DiscussionRepository
}

// NewDraftHandler は新しいDraftHandlerを作成します
func NewDraftHandler(
	issueRepo repositories.IssueRepository,
	discussionRepo repositories.DiscussionRepository,
) *DraftHandler {
	return &DraftHandler{
		issueRepo:      issueRepo,
		discussionRepo: discussionRepo,
	}
}

// IssueDraftRequest はIssueドラフトリクエストのデータ構造
type IssueDraftRequest struct {
	Title       string   `json:"title"`
	Body        string   `json:"body"`
	Labels      []string `json:"labels"`
	AssigneeID  int64    `json:"assignee_id"`
	MilestoneID int64    `json:"milestone_id"`
}

// DiscussionDraftRequest はDiscussionドラフトリクエストのデータ構造
type DiscussionDraftRequest struct {
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Category string   `json:"category"`
	Labels   []string `json:"labels"`
}

// @Summary Issueドラフトの保存
// @Description 新規または既存のIssueドラフトを保存します
// @Tags drafts
// @Accept json
// @Produce json
// @Param id path int false "Issue ID（新規の場合は省略）"
// @Param draft body IssueDraftRequest true "Issueドラフト情報"
// @Success 200 {object} models.Issue
// @Router /api/v1/drafts/issues [post]
// @Router /api/v1/drafts/issues/{id} [put]
func (h *DraftHandler) SaveIssueDraft(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// リクエストの解析
	var req IssueDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// IDの取得（パスパラメータから）
	idStr := c.Param("id")
	var issue *models.Issue
	var err error

	if idStr != "" {
		// 既存のIssueを更新
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID format"})
			return
		}

		// データベースから取得
		issue, err = h.issueRepo.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 見つからない場合
		if issue == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
			return
		}

		// 作成者のみ編集可能（または管理者）
		isAdmin, _ := c.Get("is_admin")
		if issue.CreatorID != userID.(int64) && !(isAdmin != nil && isAdmin.(bool)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this issue"})
			return
		}

		// Issueの更新
		if req.Title != "" {
			issue.Title = req.Title
		}
		if req.Body != "" {
			issue.Body = req.Body
		}
		issue.AssigneeID = req.AssigneeID
		issue.MilestoneID = req.MilestoneID
		issue.IsDraft = true
		issue.UpdatedAt = models.CurrentTime()

		// ラベルの更新（指定されている場合のみ）
		if req.Labels != nil {
			issue.Labels = []string{} // 既存のラベルをクリア
			for _, label := range req.Labels {
				issue.AddLabel(label)
			}
		}

		// データベースに保存
		err = h.issueRepo.Update(c.Request.Context(), issue)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 新規Issueを作成
		title := req.Title
		if title == "" {
			title = "Draft Issue"
		}

		issue = models.NewIssue(title, req.Body, userID.(int64))
		issue.IsDraft = true
		issue.AssigneeID = req.AssigneeID
		issue.MilestoneID = req.MilestoneID

		// ラベルの設定
		for _, label := range req.Labels {
			issue.AddLabel(label)
		}

		// データベースに保存
		err = h.issueRepo.Create(c.Request.Context(), issue)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, issue)
}

// @Summary Discussionドラフトの保存
// @Description 新規または既存のDiscussionドラフトを保存します
// @Tags drafts
// @Accept json
// @Produce json
// @Param id path int false "Discussion ID（新規の場合は省略）"
// @Param draft body DiscussionDraftRequest true "Discussionドラフト情報"
// @Success 200 {object} models.Discussion
// @Router /api/v1/drafts/discussions [post]
// @Router /api/v1/drafts/discussions/{id} [put]
func (h *DraftHandler) SaveDiscussionDraft(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// リクエストの解析
	var req DiscussionDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// IDの取得（パスパラメータから）
	idStr := c.Param("id")
	var discussion *models.Discussion
	var err error

	if idStr != "" {
		// 既存のDiscussionを更新
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discussion ID format"})
			return
		}

		// データベースから取得
		discussion, err = h.discussionRepo.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 見つからない場合
		if discussion == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Discussion not found"})
			return
		}

		// 作成者のみ編集可能（または管理者）
		isAdmin, _ := c.Get("is_admin")
		if discussion.CreatorID != userID.(int64) && !(isAdmin != nil && isAdmin.(bool)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this discussion"})
			return
		}

		// Discussionの更新
		if req.Title != "" {
			discussion.Title = req.Title
		}
		if req.Body != "" {
			discussion.Body = req.Body
		}
		if req.Category != "" {
			discussion.Category = req.Category
		}
		discussion.IsDraft = true
		discussion.UpdatedAt = models.CurrentTime()

		// ラベルの更新（指定されている場合のみ）
		if req.Labels != nil {
			discussion.Labels = []string{} // 既存のラベルをクリア
			for _, label := range req.Labels {
				discussion.AddLabel(label)
			}
		}

		// データベースに保存
		err = h.discussionRepo.Update(c.Request.Context(), discussion)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 新規Discussionを作成
		title := req.Title
		if title == "" {
			title = "Draft Discussion"
		}

		category := req.Category
		if category == "" {
			category = "general"
		}

		discussion = models.NewDiscussion(title, req.Body, category, userID.(int64))
		discussion.IsDraft = true

		// ラベルの設定
		for _, label := range req.Labels {
			discussion.AddLabel(label)
		}

		// データベースに保存
		err = h.discussionRepo.Create(c.Request.Context(), discussion)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, discussion)
}

// @Summary ドラフト一覧の取得
// @Description ユーザーのドラフト一覧を取得します
// @Tags drafts
// @Accept json
// @Produce json
// @Param type query string false "ドラフトタイプ (issues/discussions/all)" default("all")
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/drafts [get]
func (h *DraftHandler) ListDrafts(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// クエリパラメータの取得
	draftType := c.DefaultQuery("type", "all")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 返却用のデータ構造
	result := gin.H{
		"page":  page,
		"limit": limit,
	}

	// フィルタの作成（自分のドラフトのみ）
	filter := map[string]interface{}{
		"creator_id": userID.(int64),
		"is_draft":   true,
	}

	// ドラフトタイプに応じてデータを取得
	if draftType == "issues" || draftType == "all" {
		issues, total, err := h.issueRepo.List(c.Request.Context(), filter, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result["issues"] = issues
		result["issues_total"] = total
	}

	if draftType == "discussions" || draftType == "all" {
		discussions, total, err := h.discussionRepo.List(c.Request.Context(), filter, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result["discussions"] = discussions
		result["discussions_total"] = total
	}

	c.JSON(http.StatusOK, result)
}

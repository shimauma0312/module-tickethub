package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
)

// IssueHandler はIssue関連のハンドラーを管理する構造体
type IssueHandler struct {
	issueRepo     repositories.IssueRepository
	labelRepo     repositories.LabelRepository
	milestoneRepo repositories.MilestoneRepository
	userRepo      repositories.UserRepository
}

// NewIssueHandler は新しいIssueHandlerを作成します
func NewIssueHandler(
	issueRepo repositories.IssueRepository,
	labelRepo repositories.LabelRepository,
	milestoneRepo repositories.MilestoneRepository,
	userRepo repositories.UserRepository,
) *IssueHandler {
	return &IssueHandler{
		issueRepo:     issueRepo,
		labelRepo:     labelRepo,
		milestoneRepo: milestoneRepo,
		userRepo:      userRepo,
	}
}

// IssueRequest はIssue作成・更新リクエストのデータ構造
type IssueRequest struct {
	Title       string   `json:"title" binding:"required"`
	Body        string   `json:"body"`
	Labels      []string `json:"labels"`
	AssigneeID  int64    `json:"assignee_id"`
	MilestoneID int64    `json:"milestone_id"`
	IsDraft     bool     `json:"is_draft"`
}

// @Summary Issue一覧の取得
// @Description Issueの一覧を取得します
// @Tags issues
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(10)
// @Param status query string false "ステータス (open/closed)" default("open")
// @Param label query string false "ラベル名"
// @Param assignee query int false "担当者ID"
// @Param milestone query int false "マイルストーンID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/issues [get]
func (h *IssueHandler) ListIssues(c *gin.Context) {
	// クエリパラメータの取得
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.DefaultQuery("status", "open")
	
	// projectIDStr := c.Query("project_id")
	// var projectID *int64

	// フィルタの作成
	filter := map[string]interface{}{
		"status": status,
	}

	// 担当者IDが指定されている場合
	assigneeIDStr := c.Query("assignee")
	if assigneeIDStr != "" {
		assigneeID, err := strconv.ParseInt(assigneeIDStr, 10, 64)
		if err == nil && assigneeID > 0 {
			filter["assignee_id"] = assigneeID
		}
	}

	// マイルストーンIDが指定されている場合
	milestoneIDStr := c.Query("milestone")
	if milestoneIDStr != "" {
		milestoneID, err := strconv.ParseInt(milestoneIDStr, 10, 64)
		if err == nil && milestoneID > 0 {
			filter["milestone_id"] = milestoneID
		}
	}

	// プロジェクトIDが指定されている場合
	// if projectIDStr != "" {
	// 	projectIDParsed, err := strconv.ParseInt(projectIDStr, 10, 64)
	// 	if err == nil && projectIDParsed > 0 {
	// 		projectID = &projectIDParsed
	// 		filter["project_id"] = projectIDParsed
	// 	}
	// }

	// ラベルが指定されている場合は別途処理が必要
	// （実際の実装ではリポジトリレイヤーでラベルフィルタリングを行う）

	// データベースから取得
	issues, total, err := h.issueRepo.List(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"issues": issues,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// @Summary Issueの取得
// @Description 指定されたIDのIssueを取得します
// @Tags issues
// @Accept json
// @Produce json
// @Param id path int true "Issue ID"
// @Success 200 {object} models.Issue
// @Router /api/v1/issues/{id} [get]
func (h *IssueHandler) GetIssue(c *gin.Context) {
	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID format"})
		return
	}

	// データベースから取得
	issue, err := h.issueRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if issue == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}

	c.JSON(http.StatusOK, issue)
}

// @Summary Issueの作成
// @Description 新しいIssueを作成します
// @Tags issues
// @Accept json
// @Produce json
// @Param issue body IssueRequest true "Issue情報"
// @Success 201 {object} models.Issue
// @Router /api/v1/issues [post]
func (h *IssueHandler) CreateIssue(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// リクエストの解析
	var req IssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Issueの作成
	issue := models.NewIssue(req.Title, req.Body, userID.(int64))
	issue.IsDraft = req.IsDraft
	issue.AssigneeID = req.AssigneeID
	issue.MilestoneID = req.MilestoneID

	// ラベルの設定
	for _, label := range req.Labels {
		issue.AddLabel(label)
	}

	// データベースに保存
	err := h.issueRepo.Create(c.Request.Context(), issue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, issue)
}

// @Summary Issueの更新
// @Description 指定されたIDのIssueを更新します
// @Tags issues
// @Accept json
// @Produce json
// @Param id path int true "Issue ID"
// @Param issue body IssueRequest true "Issue情報"
// @Success 200 {object} models.Issue
// @Router /api/v1/issues/{id} [put]
func (h *IssueHandler) UpdateIssue(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID format"})
		return
	}

	// データベースから取得
	issue, err := h.issueRepo.GetByID(c.Request.Context(), id)
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

	// リクエストの解析
	var req IssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Issueの更新
	issue.Title = req.Title
	issue.Body = req.Body
	issue.IsDraft = req.IsDraft
	issue.AssigneeID = req.AssigneeID
	issue.MilestoneID = req.MilestoneID
	issue.UpdatedAt = models.CurrentTime()

	// ラベルの更新
	issue.Labels = []string{} // 既存のラベルをクリア
	for _, label := range req.Labels {
		issue.AddLabel(label)
	}

	// データベースに保存
	err = h.issueRepo.Update(c.Request.Context(), issue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, issue)
}

// @Summary Issueの削除
// @Description 指定されたIDのIssueを削除します
// @Tags issues
// @Accept json
// @Produce json
// @Param id path int true "Issue ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/issues/{id} [delete]
func (h *IssueHandler) DeleteIssue(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID format"})
		return
	}

	// データベースから取得
	issue, err := h.issueRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if issue == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}

	// 作成者のみ削除可能（または管理者）
	isAdmin, _ := c.Get("is_admin")
	if issue.CreatorID != userID.(int64) && !(isAdmin != nil && isAdmin.(bool)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this issue"})
		return
	}

	// データベースから削除
	err = h.issueRepo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Issue deleted successfully"})
}

// @Summary Issueのステータス変更
// @Description 指定されたIDのIssueのステータスを変更します
// @Tags issues
// @Accept json
// @Produce json
// @Param id path int true "Issue ID"
// @Param status body map[string]string true "ステータス情報"
// @Success 200 {object} models.Issue
// @Router /api/v1/issues/{id}/status [patch]
func (h *IssueHandler) UpdateIssueStatus(c *gin.Context) {
	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID format"})
		return
	}

	// データベースから取得
	issue, err := h.issueRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 見つからない場合
	if issue == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}

	// リクエストの解析
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ステータスの更新
	switch req.Status {
	case "open":
		issue.Reopen()
	case "closed":
		issue.Close()
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be 'open' or 'closed'"})
		return
	}

	// データベースに保存
	err = h.issueRepo.Update(c.Request.Context(), issue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, issue)
}

// @Summary Issueのドラフト状態変更
// @Description 指定されたIDのIssueのドラフト状態を変更します
// @Tags issues
// @Accept json
// @Produce json
// @Param id path int true "Issue ID"
// @Param draft body map[string]bool true "ドラフト情報"
// @Success 200 {object} models.Issue
// @Router /api/v1/issues/{id}/draft [patch]
func (h *IssueHandler) UpdateIssueDraftStatus(c *gin.Context) {
	// ユーザーIDの取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// IDの取得
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID format"})
		return
	}

	// データベースから取得
	issue, err := h.issueRepo.GetByID(c.Request.Context(), id)
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

	// リクエストの解析
	var req struct {
		IsDraft bool `json:"is_draft"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ドラフト状態の更新
	issue.IsDraft = req.IsDraft
	issue.UpdatedAt = models.CurrentTime()

	// データベースに保存
	err = h.issueRepo.Update(c.Request.Context(), issue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, issue)
}

// @Summary Issueの検索
// @Description Issueを検索します
// @Tags issues
// @Accept json
// @Produce json
// @Param q query string true "検索クエリ"
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "1ページあたりの件数" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/issues/search [get]
func (h *IssueHandler) SearchIssues(c *gin.Context) {
	// クエリパラメータの取得
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// データベースから検索
	issues, total, err := h.issueRepo.Search(c.Request.Context(), query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"issues": issues,
		"total":  total,
		"page":   page,
		"limit":  limit,
		"query":  query,
	})
}

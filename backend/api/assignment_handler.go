package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
)

// AssignmentHandler はアサイン機能のハンドラーを管理する構造体
type AssignmentHandler struct {
	issueRepo repositories.IssueRepository
	userRepo  repositories.UserRepository
}

// NewAssignmentHandler は新しいAssignmentHandlerを作成します
func NewAssignmentHandler(
	issueRepo repositories.IssueRepository,
	userRepo repositories.UserRepository,
) *AssignmentHandler {
	return &AssignmentHandler{
		issueRepo: issueRepo,
		userRepo:  userRepo,
	}
}

// @Summary Issueの担当者設定
// @Description 指定されたIssueの担当者を設定します
// @Tags assignments
// @Accept json
// @Produce json
// @Param id path int true "Issue ID"
// @Param assignment body map[string]int64 true "担当者情報"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/issues/{id}/assign [put]
func (h *AssignmentHandler) AssignIssue(c *gin.Context) {
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
		AssigneeID int64 `json:"assignee_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 担当者IDの検証（0は担当者なしを意味する）
	if req.AssigneeID > 0 {
		user, err := h.userRepo.GetByID(c.Request.Context(), req.AssigneeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if user == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Assignee user not found"})
			return
		}
	}

	// 担当者の更新
	issue.AssigneeID = req.AssigneeID
	issue.UpdatedAt = models.CurrentTime()

	// データベースに保存
	err = h.issueRepo.Update(c.Request.Context(), issue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Issue assigned successfully",
		"issue_id":    issue.ID,
		"assignee_id": issue.AssigneeID,
	})
}

// @Summary Issueの担当者削除
// @Description 指定されたIssueの担当者を削除します
// @Tags assignments
// @Accept json
// @Produce json
// @Param id path int true "Issue ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/issues/{id}/unassign [put]
func (h *AssignmentHandler) UnassignIssue(c *gin.Context) {
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

	// 担当者の削除
	issue.AssigneeID = 0
	issue.UpdatedAt = models.CurrentTime()

	// データベースに保存
	err = h.issueRepo.Update(c.Request.Context(), issue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Issue unassigned successfully",
		"issue_id": issue.ID,
	})
}

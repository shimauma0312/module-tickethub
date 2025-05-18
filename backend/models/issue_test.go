package models_test

import (
	"testing"
	"time"

	"github.com/shimauma0312/module-tickethub/models"
	"github.com/stretchr/testify/assert"
)

func TestNewIssue(t *testing.T) {
	title := "テストIssue"
	body := "これはテスト用のIssueです"
	creatorID := int64(1)

	issue := models.NewIssue(title, body, creatorID)

	assert.Equal(t, title, issue.Title)
	assert.Equal(t, body, issue.Body)
	assert.Equal(t, creatorID, issue.CreatorID)
	assert.Equal(t, "open", issue.Status)
	assert.False(t, issue.IsDraft)
	assert.NotEmpty(t, issue.CreatedAt)
	assert.NotEmpty(t, issue.UpdatedAt)
	assert.Empty(t, issue.Labels)
}

func TestIssue_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		issue models.Issue
		want  bool
	}{
		{
			name: "有効なIssue",
			issue: models.Issue{
				Title:     "有効なIssue",
				CreatorID: 1,
			},
			want: true,
		},
		{
			name: "タイトルなし",
			issue: models.Issue{
				Title:     "",
				CreatorID: 1,
			},
			want: false,
		},
		{
			name: "作成者IDなし",
			issue: models.Issue{
				Title:     "作成者IDなし",
				CreatorID: 0,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.issue.IsValid())
		})
	}
}

func TestIssue_Status(t *testing.T) {
	issue := models.NewIssue("テストIssue", "説明", 1)

	// 初期状態確認
	assert.Equal(t, "open", issue.Status)

	// クローズしてステータス確認
	beforeUpdate := issue.UpdatedAt
	time.Sleep(1 * time.Millisecond) // UpdatedAtの変更を確認するため
	issue.Close()
	assert.Equal(t, "closed", issue.Status)
	assert.True(t, issue.UpdatedAt.After(beforeUpdate))

	// 再オープンしてステータス確認
	beforeUpdate = issue.UpdatedAt
	time.Sleep(1 * time.Millisecond)
	issue.Reopen()
	assert.Equal(t, "open", issue.Status)
	assert.True(t, issue.UpdatedAt.After(beforeUpdate))
}

func TestIssue_Labels(t *testing.T) {
	issue := models.NewIssue("テストIssue", "説明", 1)

	// ラベル追加
	issue.AddLabel("バグ")
	assert.Contains(t, issue.Labels, "バグ")
	assert.Len(t, issue.Labels, 1)

	// 同じラベルを追加しても重複しないことを確認
	issue.AddLabel("バグ")
	assert.Len(t, issue.Labels, 1)

	// 別のラベルを追加
	issue.AddLabel("機能要望")
	assert.Contains(t, issue.Labels, "機能要望")
	assert.Len(t, issue.Labels, 2)

	// ラベル削除
	issue.RemoveLabel("バグ")
	assert.NotContains(t, issue.Labels, "バグ")
	assert.Contains(t, issue.Labels, "機能要望")
	assert.Len(t, issue.Labels, 1)

	// 存在しないラベルの削除は何も起こらない
	issue.RemoveLabel("存在しないラベル")
	assert.Len(t, issue.Labels, 1)
}

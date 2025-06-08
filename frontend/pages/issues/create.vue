<template>
  <div class="issue-create">
    <div class="issue-create-header">
      <h1 class="issue-create-title">新規Issue作成</h1>
    </div>
    
    <div class="issue-form">
      <div class="form-group">
        <label for="issue-title" class="form-label">タイトル</label>
        <Input 
          id="issue-title" 
          v-model="issue.title" 
          placeholder="Issueのタイトルを入力" 
          class="issue-title-input"
        />
      </div>
      
      <div class="form-group">
        <label for="issue-body" class="form-label">内容</label>
        <MarkdownEditor 
          id="issue-body"
          v-model="issue.body" 
          :placeholder="bodyPlaceholder" 
          :minHeight="300"
          :showPreview="true"
        />
      </div>
      
      <div class="issue-form-sidebar">
        <div class="form-sidebar-section">
          <h3 class="sidebar-heading">ラベル</h3>
          <div v-if="selectedLabels.length" class="selected-labels">
            <Label 
              v-for="label in selectedLabels" 
              :key="label.id" 
              :name="label.name" 
              :color="label.color"
              closable
              @close="removeLabel(label.id)"
            />
          </div>
          <div v-else class="no-selection">
            ラベルが選択されていません
          </div>
          <Select 
            placeholder="ラベルを選択"
            :options="labelOptions"
            @change="addLabel"
          />
        </div>
        
        <div class="form-sidebar-section">
          <h3 class="sidebar-heading">マイルストーン</h3>
          <div v-if="issue.milestone" class="selected-milestone">
            {{ issue.milestone.title }}
          </div>
          <div v-else class="no-selection">
            マイルストーンが選択されていません
          </div>
          <Select 
            placeholder="マイルストーンを選択"
            :options="milestoneOptions"
            @change="selectMilestone"
          />
        </div>
        
        <div class="form-sidebar-section">
          <h3 class="sidebar-heading">担当者</h3>
          <div v-if="selectedAssignees.length" class="selected-assignees">
            <div v-for="assignee in selectedAssignees" :key="assignee.id" class="assignee-item">
              {{ assignee.name }}
              <button class="remove-btn" @click="removeAssignee(assignee.id)">×</button>
            </div>
          </div>
          <div v-else class="no-selection">
            担当者が選択されていません
          </div>
          <Select 
            placeholder="担当者を選択"
            :options="userOptions"
            @change="addAssignee"
          />
        </div>
      </div>
      
      <div class="form-actions">
        <Button variant="outline" @click="saveDraft" class="mr-2">
          下書き保存
        </Button>
        <Button variant="primary" @click="submitIssue" :disabled="!issue.title">
          Issue作成
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Issue データ
const issue = ref({
  title: '',
  body: '',
  labels: [],
  milestone: null,
  assignees: []
})

// ドラフト保存用のID
const draftId = ref(null)

// 選択済みラベル
const selectedLabels = computed(() => {
  return issue.value.labels || []
})

// 選択済み担当者
const selectedAssignees = computed(() => {
  return issue.value.assignees || []
})

// プレースホルダー
const bodyPlaceholder = `## 概要
簡潔に問題・課題を説明してください。

## 詳細
さらに詳しい情報を記載してください。

## 期待される結果
どのような結果を期待しているかを記載してください。

## 参考リンク
関連する資料やリンクを記載してください。`

// ラベルオプション（デモデータ）
const labelOptions = [
  { value: '1', label: 'バグ', data: { id: 1, name: 'バグ', color: '#d73a4a' } },
  { value: '2', label: '機能追加', data: { id: 2, name: '機能追加', color: '#0366d6' } },
  { value: '3', label: 'ドキュメント', data: { id: 3, name: 'ドキュメント', color: '#0075ca' } },
  { value: '4', label: '質問', data: { id: 4, name: '質問', color: '#d876e3' } },
  { value: '5', label: '優先度:高', data: { id: 5, name: '優先度:高', color: '#e99695' } },
  { value: '6', label: '優先度:中', data: { id: 6, name: '優先度:中', color: '#fbca04' } },
  { value: '7', label: '優先度:低', data: { id: 7, name: '優先度:低', color: '#c2e0c6' } }
]

// マイルストーンオプション（デモデータ）
const milestoneOptions = [
  { value: '1', label: 'v1.0 初期リリース', data: { id: 1, title: 'v1.0 初期リリース' } },
  { value: '2', label: 'v1.1 機能強化', data: { id: 2, title: 'v1.1 機能強化' } },
  { value: '3', label: 'v1.2 パフォーマンス改善', data: { id: 3, title: 'v1.2 パフォーマンス改善' } }
]

// ユーザーオプション（デモデータ）
const userOptions = [
  { value: '1', label: 'shimauma0312', data: { id: 1, name: 'shimauma0312' } },
  { value: '2', label: 'tanuki456', data: { id: 2, name: 'tanuki456' } },
  { value: '3', label: 'kitsune789', data: { id: 3, name: 'kitsune789' } }
]

// ラベル追加
function addLabel(selectedOption) {
  if (!selectedOption) return
  
  const labelData = selectedOption.data
  
  // 既に選択済みかチェック
  if (!issue.value.labels) {
    issue.value.labels = []
  }
  
  if (!issue.value.labels.some(label => label.id === labelData.id)) {
    issue.value.labels.push(labelData)
  }
}

// ラベル削除
function removeLabel(labelId) {
  if (!issue.value.labels) return
  
  issue.value.labels = issue.value.labels.filter(label => label.id !== labelId)
}

// マイルストーン選択
function selectMilestone(selectedOption) {
  if (!selectedOption) {
    issue.value.milestone = null
    return
  }
  
  issue.value.milestone = selectedOption.data
}

// 担当者追加
function addAssignee(selectedOption) {
  if (!selectedOption) return
  
  const assigneeData = selectedOption.data
  
  // 既に選択済みかチェック
  if (!issue.value.assignees) {
    issue.value.assignees = []
  }
  
  if (!issue.value.assignees.some(assignee => assignee.id === assigneeData.id)) {
    issue.value.assignees.push(assigneeData)
  }
}

// 担当者削除
function removeAssignee(assigneeId) {
  if (!issue.value.assignees) return
  
  issue.value.assignees = issue.value.assignees.filter(assignee => assignee.id !== assigneeId)
}

// Issue作成
async function submitIssue() {
  if (!issue.value.title.trim()) return
  
  try {
    // APIに送信（実装予定）
    console.log('Issue作成:', issue.value)
    
    // 成功したら下書きを削除
    if (draftId.value) {
      await deleteDraft(draftId.value)
    }
    
    // Issue一覧ページに遷移
    router.push('/issues')
  } catch (error) {
    console.error('Issue作成中にエラーが発生しました:', error)
    // エラーハンドリング
  }
}

// 下書き保存
async function saveDraft() {
  try {
    // APIに下書き保存（実装予定）
    console.log('下書き保存:', issue.value)
    
    // 下書きIDを保存
    draftId.value = 'draft-' + Date.now()
    
    // 通知表示
    alert('下書きを保存しました')
  } catch (error) {
    console.error('下書き保存中にエラーが発生しました:', error)
    // エラーハンドリング
  }
}

// 下書き削除
async function deleteDraft(id) {
  try {
    // APIから下書き削除（実装予定）
    console.log('下書き削除:', id)
  } catch (error) {
    console.error('下書き削除中にエラーが発生しました:', error)
  }
}
</script>

<style scoped>
.issue-create {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.issue-create-header {
  margin-bottom: 2rem;
}

.issue-create-title {
  font-size: 1.75rem;
  font-weight: 600;
}

.issue-form {
  display: grid;
  grid-template-columns: 3fr 1fr;
  gap: 2rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  font-weight: 500;
  margin-bottom: 0.5rem;
}

.issue-title-input {
  width: 100%;
  font-size: 1.125rem;
}

.issue-form-sidebar {
  grid-column: 2;
  align-self: start;
}

.form-sidebar-section {
  background-color: white;
  border: 1px solid var(--color-border-primary);
  border-radius: 0.375rem;
  padding: 1rem;
  margin-bottom: 1rem;
}

.sidebar-heading {
  font-size: 0.875rem;
  font-weight: 600;
  margin-top: 0;
  margin-bottom: 0.75rem;
  color: var(--color-text-secondary);
}

.no-selection {
  color: var(--color-text-tertiary);
  font-size: 0.875rem;
  font-style: italic;
  margin-bottom: 0.75rem;
}

.selected-labels {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.selected-milestone {
  margin-bottom: 0.75rem;
  font-size: 0.875rem;
}

.selected-assignees {
  margin-bottom: 0.75rem;
}

.assignee-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
}

.remove-btn {
  background: none;
  border: none;
  color: var(--color-text-tertiary);
  cursor: pointer;
  font-size: 1rem;
  padding: 0 0.25rem;
}

.remove-btn:hover {
  color: var(--color-text-primary);
}

.form-actions {
  grid-column: 1 / 3;
  display: flex;
  justify-content: flex-end;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border-primary);
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .issue-form {
    grid-template-columns: 1fr;
  }
  
  .issue-form-sidebar {
    grid-column: 1;
  }
  
  .form-actions {
    grid-column: 1;
  }
}
</style>

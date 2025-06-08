<template>
  <div class="issue-detail">
    <!-- Issue ヘッダー部分 -->
    <div class="issue-header">
      <div class="issue-title-container">
        <h1 class="issue-title">
          {{ issue.title }} <span class="issue-number">#{{ issue.id }}</span>
        </h1>
        <div class="issue-status" :class="getStatusClass(issue.status)">
          {{ getStatusLabel(issue.status) }}
        </div>
      </div>
      
      <div class="issue-meta">
        <span class="issue-created">
          <strong>{{ issue.created_by }}</strong> さんが {{ formatDate(issue.created_at) }} に作成
        </span>
        <span class="issue-comments">コメント: {{ issue.comments_count || 0 }}</span>
      </div>
    </div>
    
    <!-- Issue 本文とラベル、担当者情報 -->
    <div class="issue-content-wrapper">
      <div class="issue-content">
        <!-- Issue 本文 -->
        <div class="issue-body markdown-body">
          <div v-if="issue.body" v-html="renderMarkdown(issue.body)"></div>
          <div v-else class="issue-body-empty">本文はありません</div>
        </div>
        
        <!-- コメントスレッド -->
        <div class="issue-comments-thread">
          <h3 class="comments-heading">コメント</h3>
          
          <div v-if="!comments.length" class="no-comments">
            まだコメントはありません。会話を始めましょう！
          </div>
          
          <div v-for="comment in comments" :key="comment.id" class="comment">
            <div class="comment-header">
              <strong>{{ comment.author }}</strong>
              <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
            </div>
            <div class="comment-body markdown-body" v-html="renderMarkdown(comment.body)"></div>
            <div class="comment-reactions">
              <button class="reaction-button" @click="showReactionPicker(comment.id)">
                <span class="reaction-icon">😀</span> リアクション
              </button>
              <div v-for="(reaction, emoji) in comment.reactions" :key="emoji" class="reaction-badge">
                {{ emoji }} {{ reaction.count }}
              </div>
            </div>
          </div>
          
          <!-- 新規コメント入力エリア -->
          <div class="new-comment">
            <h4>コメントを追加</h4>
            <MarkdownEditor 
              v-model="newComment" 
              :placeholder="'コメントを入力してください'" 
              :minHeight="150"
            />
            <div class="comment-actions">
              <Button 
                @click="submitComment" 
                :disabled="!newComment.trim()" 
                variant="primary"
              >
                コメントする
              </Button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- サイドバー（ラベル、マイルストーン、担当者） -->
      <div class="issue-sidebar">
        <div class="sidebar-section">
          <h3 class="sidebar-heading">ステータス</h3>
          <Select 
            v-model="issue.status" 
            :options="statusOptions" 
            @change="updateIssueStatus"
          />
        </div>
        
        <div class="sidebar-section">
          <h3 class="sidebar-heading">担当者</h3>
          <div v-if="issue.assignees && issue.assignees.length" class="assignees-list">
            <div v-for="assignee in issue.assignees" :key="assignee.id" class="assignee-item">
              {{ assignee.name }}
            </div>
          </div>
          <div v-else class="no-assignees">
            担当者が設定されていません
          </div>
          <Button variant="outline" size="sm" class="mt-2" @click="showAssigneeSelector">
            担当者を変更
          </Button>
        </div>
        
        <div class="sidebar-section">
          <h3 class="sidebar-heading">ラベル</h3>
          <div v-if="issue.labels && issue.labels.length" class="labels-list">
            <Label 
              v-for="label in issue.labels" 
              :key="label.id" 
              :color="label.color" 
              :name="label.name"
            />
          </div>
          <div v-else class="no-labels">
            ラベルが設定されていません
          </div>
          <Button variant="outline" size="sm" class="mt-2" @click="showLabelSelector">
            ラベルを変更
          </Button>
        </div>
        
        <div class="sidebar-section">
          <h3 class="sidebar-heading">マイルストーン</h3>
          <div v-if="issue.milestone" class="milestone-info">
            {{ issue.milestone.title }}
          </div>
          <div v-else class="no-milestone">
            マイルストーンが設定されていません
          </div>
          <Button variant="outline" size="sm" class="mt-2" @click="showMilestoneSelector">
            マイルストーンを変更
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const route = useRoute()
const router = useRouter()
const issueId = computed(() => route.params.id)

// Issue データ
const issue = ref({
  id: null,
  title: '',
  body: '',
  status: 'open',
  created_by: '',
  created_at: '',
  assignees: [],
  labels: [],
  milestone: null
})

// コメントデータ
const comments = ref([])
const newComment = ref('')

// ステータスオプション
const statusOptions = [
  { value: 'open', label: '未対応' },
  { value: 'in-progress', label: '対応中' },
  { value: 'closed', label: '完了' }
]

// Markdownレンダリング
function renderMarkdown(text) {
  if (!text) return ''
  const renderedHtml = marked.parse(text)
  return DOMPurify.sanitize(renderedHtml)
}

// 日付フォーマット
function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('ja-JP', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

// ステータスに応じたクラスを返す
function getStatusClass(status) {
  switch (status) {
    case 'open':
      return 'status-open'
    case 'in-progress':
      return 'status-progress'
    case 'closed':
      return 'status-closed'
    default:
      return 'status-open'
  }
}

// ステータスラベルを取得
function getStatusLabel(status) {
  const option = statusOptions.find(opt => opt.value === status)
  return option ? option.label : '未対応'
}

// Issueステータス更新
function updateIssueStatus() {
  // APIリクエスト（実装予定）
  console.log('Issueステータスを更新:', issue.value.status)
}

// コメント投稿
function submitComment() {
  if (!newComment.value.trim()) return
  
  // 新規コメントをAPIに送信（実装予定）
  console.log('新規コメント:', newComment.value)
  
  // 仮実装：コメント追加
  comments.value.push({
    id: Date.now(),
    author: 'Current User',
    body: newComment.value,
    created_at: new Date().toISOString(),
    reactions: {}
  })
  
  // 入力フィールドをクリア
  newComment.value = ''
}

// 担当者選択ダイアログを表示
function showAssigneeSelector() {
  // 実装予定
  console.log('担当者選択ダイアログを表示')
}

// ラベル選択ダイアログを表示
function showLabelSelector() {
  // 実装予定
  console.log('ラベル選択ダイアログを表示')
}

// マイルストーン選択ダイアログを表示
function showMilestoneSelector() {
  // 実装予定
  console.log('マイルストーン選択ダイアログを表示')
}

// リアクション選択パネルを表示
function showReactionPicker(commentId) {
  // 実装予定
  console.log('リアクション選択:', commentId)
}

onMounted(async () => {
  // デモデータ（後でAPIに置き換え）
  issue.value = {
    id: issueId.value,
    title: 'プロジェクト初期構成とDockerセットアップ',
    body: `## 概要
Docker環境を使った開発環境を構築します。

### 要件
- フロントエンド、バックエンド、DBの3コンテナ構成
- 開発環境と本番環境の設定ファイル分離
- ホットリロード対応

## 参考リンク
- [Docker Compose ドキュメント](https://docs.docker.com/compose/)
- [Nuxt.js with Docker](https://nuxtjs.org/docs/deployment/docker)`,
    status: 'in-progress',
    created_by: 'shimauma0312',
    created_at: '2025-05-15T09:30:00.000Z',
    assignees: [
      { id: 1, name: 'shimauma0312' }
    ],
    labels: [
      { id: 1, name: '機能追加', color: '#0366d6' },
      { id: 2, name: '優先度:高', color: '#d73a4a' }
    ],
    milestone: { id: 1, title: 'v1.0 初期リリース' }
  }
  
  comments.value = [
    {
      id: 1,
      author: 'shimauma0312',
      body: 'Docker ComposeのYAMLファイルを作成しました。レビューをお願いします。',
      created_at: '2025-05-16T10:15:00.000Z',
      reactions: {
        '👍': { count: 2 },
        '🎉': { count: 1 }
      }
    },
    {
      id: 2,
      author: 'tanuki456',
      body: 'ホットリロードの設定が必要そうです。frontendのDockerfileを修正しましょう。',
      created_at: '2025-05-17T11:20:00.000Z',
      reactions: {}
    }
  ]
})
</script>

<style scoped>
.issue-detail {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.issue-header {
  margin-bottom: 2rem;
  border-bottom: 1px solid var(--color-border-primary);
  padding-bottom: 1rem;
}

.issue-title-container {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 0.5rem;
}

.issue-title {
  font-size: 1.75rem;
  font-weight: 600;
  margin: 0;
}

.issue-number {
  color: var(--color-text-tertiary);
  font-weight: normal;
}

.issue-status {
  padding: 0.25rem 0.75rem;
  border-radius: 2rem;
  font-size: 0.875rem;
  font-weight: 500;
}

.status-open {
  background-color: #2da44e;
  color: white;
}

.status-progress {
  background-color: #bf8700;
  color: white;
}

.status-closed {
  background-color: #8250df;
  color: white;
}

.issue-meta {
  display: flex;
  gap: 1.5rem;
  font-size: 0.875rem;
  color: var(--color-text-secondary);
}

.issue-content-wrapper {
  display: grid;
  grid-template-columns: 3fr 1fr;
  gap: 2rem;
}

.issue-body {
  background-color: white;
  border: 1px solid var(--color-border-primary);
  border-radius: 0.375rem;
  padding: 1rem;
  margin-bottom: 2rem;
}

.issue-body-empty {
  color: var(--color-text-tertiary);
  font-style: italic;
}

.comments-heading {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
}

.no-comments {
  background-color: var(--color-bg-secondary);
  border-radius: 0.375rem;
  padding: 1rem;
  text-align: center;
  color: var(--color-text-secondary);
}

.comment {
  background-color: white;
  border: 1px solid var(--color-border-primary);
  border-radius: 0.375rem;
  margin-bottom: 1.5rem;
  overflow: hidden;
}

.comment-header {
  background-color: var(--color-bg-secondary);
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--color-border-primary);
}

.comment-date {
  font-size: 0.875rem;
  color: var(--color-text-tertiary);
  margin-left: 0.5rem;
}

.comment-body {
  padding: 1rem;
}

.comment-reactions {
  padding: 0.5rem 1rem;
  border-top: 1px solid var(--color-border-primary);
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.reaction-button {
  background-color: transparent;
  border: 1px solid var(--color-border-primary);
  border-radius: 2rem;
  padding: 0.25rem 0.75rem;
  font-size: 0.875rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.reaction-icon {
  font-size: 1rem;
}

.reaction-badge {
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-border-primary);
  border-radius: 2rem;
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
}

.new-comment {
  background-color: white;
  border: 1px solid var(--color-border-primary);
  border-radius: 0.375rem;
  padding: 1rem;
  margin-top: 2rem;
}

.comment-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 1rem;
}

.issue-sidebar {
  align-self: start;
}

.sidebar-section {
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

.assignees-list, .labels-list {
  margin-bottom: 0.75rem;
}

.assignee-item {
  margin-bottom: 0.25rem;
  font-size: 0.875rem;
}

.no-assignees, .no-labels, .no-milestone {
  color: var(--color-text-tertiary);
  font-size: 0.875rem;
  font-style: italic;
  margin-bottom: 0.75rem;
}

.milestone-info {
  font-size: 0.875rem;
  margin-bottom: 0.75rem;
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .issue-content-wrapper {
    grid-template-columns: 1fr;
  }
  
  .issue-sidebar {
    order: -1;
    margin-bottom: 2rem;
  }
}

/* Markdown表示スタイル */
.markdown-body h1 {
  font-size: 1.5rem;
  margin-top: 1.5rem;
  margin-bottom: 1rem;
  font-weight: 600;
}

.markdown-body h2 {
  font-size: 1.25rem;
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
  font-weight: 600;
}

.markdown-body h3 {
  font-size: 1.125rem;
  margin-top: 1.25rem;
  margin-bottom: 0.75rem;
  font-weight: 600;
}

.markdown-body p {
  margin-bottom: 1rem;
}

.markdown-body ul, .markdown-body ol {
  margin-bottom: 1rem;
  padding-left: 2rem;
}

.markdown-body li {
  margin-bottom: 0.25rem;
}

.markdown-body code {
  background-color: var(--color-bg-secondary);
  padding: 0.2rem 0.4rem;
  border-radius: 0.25rem;
  font-family: monospace;
}

.markdown-body pre {
  background-color: var(--color-bg-secondary);
  padding: 1rem;
  border-radius: 0.375rem;
  overflow-x: auto;
  margin-bottom: 1rem;
}

.markdown-body a {
  color: var(--color-accent-primary);
  text-decoration: none;
}

.markdown-body a:hover {
  text-decoration: underline;
}

.markdown-body blockquote {
  border-left: 4px solid var(--color-border-primary);
  padding-left: 1rem;
  color: var(--color-text-secondary);
  margin: 1rem 0;
}
</style>

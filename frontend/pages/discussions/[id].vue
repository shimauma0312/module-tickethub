<template>
  <div class="discussion-detail">
    <!-- Discussion ヘッダー部分 -->
    <div class="discussion-header">
      <div class="discussion-title-container">
        <div class="discussion-category" :class="`category-${discussion.category.color}`">
          {{ discussion.category.name }}
        </div>
        <h1 class="discussion-title">
          {{ discussion.title }}
        </h1>
      </div>
      
      <div class="discussion-meta">
        <span class="discussion-author">
          <strong>{{ discussion.author }}</strong> さんが {{ formatDate(discussion.created_at) }} に作成
        </span>
        <span class="discussion-views">
          <Icon name="mdi:eye-outline" /> {{ discussion.views_count }} 閲覧
        </span>
      </div>
    </div>
    
    <!-- Discussion 本文とレスポンス -->
    <div class="discussion-content-wrapper">
      <div class="discussion-content">
        <!-- 最初の投稿（本文） -->
        <Card class="post original-post">
          <div class="post-header">
            <div class="author-info">
              <div class="author-name">{{ discussion.author }}</div>
              <div class="post-date">{{ formatDate(discussion.created_at) }}</div>
            </div>
            <div class="post-actions">
              <button class="action-btn">
                <Icon name="mdi:dots-horizontal" />
              </button>
            </div>
          </div>
          
          <div class="post-body markdown-body" v-html="renderMarkdown(discussion.content)"></div>
          
          <div class="post-reactions">
            <button class="reaction-button" @click="showReactionPicker('original')">
              <span class="reaction-icon">😀</span> リアクション
            </button>
            <div 
              v-for="(reaction, emoji) in discussion.reactions" 
              :key="emoji" 
              class="reaction-badge"
              :class="{ 'user-reacted': reaction.userReacted }"
              @click="toggleReaction('original', emoji)"
            >
              {{ emoji }} {{ reaction.count }}
            </div>
          </div>
        </Card>
        
        <!-- 返信セクション -->
        <div class="replies-section">
          <h2 class="replies-heading">
            返信 <span class="replies-count">{{ replies.length }}</span>
          </h2>
          
          <div v-if="!replies.length" class="no-replies">
            まだ返信はありません。会話を始めましょう！
          </div>
          
          <Card v-for="reply in replies" :key="reply.id" class="post reply-post">
            <div class="post-header">
              <div class="author-info">
                <div class="author-name">{{ reply.author }}</div>
                <div class="post-date">{{ formatDate(reply.created_at) }}</div>
              </div>
              <div class="post-actions">
                <button class="action-btn">
                  <Icon name="mdi:dots-horizontal" />
                </button>
              </div>
            </div>
            
            <div class="post-body markdown-body" v-html="renderMarkdown(reply.content)"></div>
            
            <div class="post-reactions">
              <button class="reaction-button" @click="showReactionPicker(reply.id)">
                <span class="reaction-icon">😀</span> リアクション
              </button>
              <div 
                v-for="(reaction, emoji) in reply.reactions" 
                :key="emoji" 
                class="reaction-badge"
                :class="{ 'user-reacted': reaction.userReacted }"
                @click="toggleReaction(reply.id, emoji)"
              >
                {{ emoji }} {{ reaction.count }}
              </div>
            </div>
          </Card>
          
          <!-- 新規返信入力エリア -->
          <Card class="new-reply">
            <h3 class="reply-heading">返信を追加</h3>
            <MarkdownEditor 
              v-model="newReply" 
              :placeholder="'返信を入力してください'" 
              :minHeight="200"
              :showPreview="true"
            />
            <div class="reply-actions">
              <Button 
                @click="submitReply" 
                :disabled="!newReply.trim()" 
                variant="primary"
              >
                返信する
              </Button>
            </div>
          </Card>
        </div>
      </div>
      
      <!-- サイドバー（関連情報） -->
      <div class="discussion-sidebar">
        <Card class="sidebar-section">
          <h3 class="sidebar-heading">ディスカッション情報</h3>
          
          <div class="info-item">
            <div class="info-label">作成者</div>
            <div class="info-value">{{ discussion.author }}</div>
          </div>
          
          <div class="info-item">
            <div class="info-label">作成日</div>
            <div class="info-value">{{ formatDate(discussion.created_at) }}</div>
          </div>
          
          <div class="info-item">
            <div class="info-label">最終更新</div>
            <div class="info-value">{{ formatDate(discussion.last_activity_at) }}</div>
          </div>
          
          <div class="info-item">
            <div class="info-label">閲覧数</div>
            <div class="info-value">{{ discussion.views_count }} 回</div>
          </div>
          
          <div class="info-item">
            <div class="info-label">返信数</div>
            <div class="info-value">{{ replies.length }} 件</div>
          </div>
        </Card>
        
        <Card class="sidebar-section">
          <h3 class="sidebar-heading">関連ディスカッション</h3>
          
          <div v-if="relatedDiscussions.length">
            <div v-for="related in relatedDiscussions" :key="related.id" class="related-discussion">
              <router-link :to="`/discussions/${related.id}`" class="related-link">
                {{ related.title }}
              </router-link>
              <div class="related-meta">
                {{ formatDate(related.created_at) }} • {{ related.replies_count }} 返信
              </div>
            </div>
          </div>
          <div v-else class="no-related">
            関連するディスカッションはありません
          </div>
        </Card>
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
const discussionId = computed(() => route.params.id)

// Discussion データ
const discussion = ref({
  id: null,
  title: '',
  content: '',
  author: '',
  created_at: '',
  last_activity_at: '',
  views_count: 0,
  category: {
    id: null,
    name: '',
    color: 'gray'
  },
  reactions: {}
})

// 返信データ
const replies = ref([])
const newReply = ref('')

// 関連ディスカッション
const relatedDiscussions = ref([])

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

// リアクション選択パネルを表示
function showReactionPicker(postId) {
  // 実装予定
  console.log('リアクション選択:', postId)
}

// リアクションの切り替え
function toggleReaction(postId, emoji) {
  console.log('リアクション切り替え:', postId, emoji)
  
  // 元の投稿へのリアクションの場合
  if (postId === 'original') {
    if (discussion.value.reactions[emoji]) {
      const userReacted = !discussion.value.reactions[emoji].userReacted
      
      // ユーザーが既にリアクションしている場合は削除、していない場合は追加
      discussion.value.reactions[emoji].count += userReacted ? 1 : -1
      discussion.value.reactions[emoji].userReacted = userReacted
      
      // カウントが0になったら削除
      if (discussion.value.reactions[emoji].count <= 0) {
        delete discussion.value.reactions[emoji]
      }
    } else {
      // 新しいリアクションを追加
      discussion.value.reactions[emoji] = {
        count: 1,
        userReacted: true
      }
    }
    return
  }
  
  // 返信へのリアクションの場合
  const replyIndex = replies.value.findIndex(reply => reply.id === postId)
  if (replyIndex >= 0) {
    const reply = replies.value[replyIndex]
    
    if (reply.reactions[emoji]) {
      const userReacted = !reply.reactions[emoji].userReacted
      
      reply.reactions[emoji].count += userReacted ? 1 : -1
      reply.reactions[emoji].userReacted = userReacted
      
      if (reply.reactions[emoji].count <= 0) {
        delete reply.reactions[emoji]
      }
    } else {
      reply.reactions[emoji] = {
        count: 1,
        userReacted: true
      }
    }
  }
}

// 返信投稿
function submitReply() {
  if (!newReply.value.trim()) return
  
  // 新規返信をAPIに送信（実装予定）
  console.log('新規返信:', newReply.value)
  
  // 仮実装：返信追加
  const newReplyObj = {
    id: Date.now(),
    author: 'Current User',
    content: newReply.value,
    created_at: new Date().toISOString(),
    reactions: {}
  }
  
  replies.value.push(newReplyObj)
  
  // 入力フィールドをクリア
  newReply.value = ''
}

onMounted(async () => {
  // デモデータ（後でAPIに置き換え）
  discussion.value = {
    id: discussionId.value,
    title: '開発環境のセットアップについて',
    content: `## 概要
Dockerを使った開発環境のセットアップ方法についてディスカッションを開始します。

### 現在の状況
現在、ローカル環境での開発にはさまざまなツールのインストールが必要で、環境差異によるトラブルが発生しています。

### 提案
Dockerを使って開発環境を統一することで、以下のメリットが期待できます：

- 環境差異の解消
- セットアップ時間の短縮
- 新メンバーの参加障壁の低減

### 検討事項
- Dockerfileの設計
- docker-compose.yml の構成
- ホットリロードの対応
- ボリュームマウントの最適化

皆さんのご意見をお願いします。`,
    author: 'shimauma0312',
    created_at: '2025-05-05T09:30:00.000Z',
    last_activity_at: '2025-06-01T14:15:00.000Z',
    views_count: 45,
    category: {
      id: 2,
      name: '開発環境',
      color: 'blue'
    },
    reactions: {
      '👍': { count: 5, userReacted: false },
      '🎉': { count: 2, userReacted: false },
      '❤️': { count: 3, userReacted: true }
    }
  }
  
  replies.value = [
    {
      id: 1,
      author: 'tanuki456',
      content: `フロントエンドとバックエンドで別々のコンテナを用意するのがいいと思います。

例えばこんな構成はどうでしょうか：
\`\`\`yaml
services:
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    volumes:
      - ./frontend:/app
  
  backend:
    build: ./backend
    ports:
      - "8000:8000"
    volumes:
      - ./backend:/app
  
  db:
    image: postgres:14
    environment:
      POSTGRES_PASSWORD: password
      POSTGRES_USER: user
      POSTGRES_DB: app
\`\`\`

ホットリロードはそれぞれのDockerfileで対応できますね。`,
      created_at: '2025-05-06T10:15:00.000Z',
      reactions: {
        '👍': { count: 3, userReacted: true }
      }
    },
    {
      id: 2,
      author: 'kitsune789',
      content: `開発と本番環境で設定を分けるために、docker-compose.override.yml を使うのがおすすめです。

基本設定は docker-compose.yml に書いて、開発環境固有の設定（ボリュームマウントやホットリロードなど）は override ファイルに記述すると管理しやすいです。`,
      created_at: '2025-05-07T11:20:00.000Z',
      reactions: {
        '👍': { count: 4, userReacted: false },
        '💡': { count: 2, userReacted: true }
      }
    }
  ]
  
  relatedDiscussions.value = [
    {
      id: 2,
      title: 'チケット管理システムの運用ルールについて',
      created_at: '2025-05-15',
      replies_count: 8
    },
    {
      id: 5,
      title: 'プロジェクトのマイルストーン計画',
      created_at: '2025-06-02',
      replies_count: 7
    }
  ]
  
  // 閲覧数をインクリメント（実際のAPIでは自動的に行われる想定）
  discussion.value.views_count++
})
</script>

<style scoped>
.discussion-detail {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.discussion-header {
  margin-bottom: 2rem;
}

.discussion-title-container {
  margin-bottom: 0.75rem;
}

.discussion-category {
  display: inline-block;
  font-size: 0.75rem;
  font-weight: 500;
  padding: 0.2rem 0.5rem;
  border-radius: 1rem;
  text-transform: uppercase;
  margin-bottom: 0.5rem;
}

.category-gray {
  background-color: #f0f0f0;
  color: #555;
}

.category-blue {
  background-color: #dbeafe;
  color: #1e40af;
}

.category-green {
  background-color: #dcfce7;
  color: #166534;
}

.category-purple {
  background-color: #f3e8ff;
  color: #7e22ce;
}

.category-red {
  background-color: #fee2e2;
  color: #b91c1c;
}

.discussion-title {
  font-size: 1.75rem;
  font-weight: 600;
  margin: 0;
}

.discussion-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 1.5rem;
  font-size: 0.875rem;
  color: var(--color-text-secondary);
}

.discussion-views {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.discussion-content-wrapper {
  display: grid;
  grid-template-columns: 3fr 1fr;
  gap: 2rem;
}

.discussion-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.post {
  display: flex;
  flex-direction: column;
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.author-info {
  display: flex;
  align-items: baseline;
  gap: 0.75rem;
}

.author-name {
  font-weight: 600;
}

.post-date {
  font-size: 0.75rem;
  color: var(--color-text-tertiary);
}

.post-actions {
  display: flex;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 0.25rem;
  color: var(--color-text-secondary);
}

.action-btn:hover {
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
}

.post-body {
  margin-bottom: 1.5rem;
}

.post-reactions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: auto;
  padding-top: 0.75rem;
  border-top: 1px solid var(--color-border-primary);
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
  cursor: pointer;
}

.reaction-badge:hover {
  background-color: var(--color-bg-primary);
}

.user-reacted {
  background-color: var(--color-bg-accent-subtle);
  border-color: var(--color-border-accent);
  color: var(--color-text-accent);
}

.replies-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.replies-heading {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--color-border-primary);
}

.replies-count {
  color: var(--color-text-secondary);
  font-weight: normal;
}

.no-replies {
  background-color: var(--color-bg-secondary);
  border-radius: 0.375rem;
  padding: 1.5rem;
  text-align: center;
  color: var(--color-text-secondary);
}

.new-reply {
  margin-top: 1rem;
}

.reply-heading {
  font-size: 1rem;
  font-weight: 600;
  margin-top: 0;
  margin-bottom: 1rem;
}

.reply-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 1rem;
}

.discussion-sidebar {
  align-self: start;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.sidebar-section {
  display: flex;
  flex-direction: column;
}

.sidebar-heading {
  font-size: 1rem;
  font-weight: 600;
  margin-top: 0;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--color-border-primary);
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.75rem;
  font-size: 0.875rem;
}

.info-label {
  color: var(--color-text-secondary);
}

.info-value {
  font-weight: 500;
}

.related-discussion {
  margin-bottom: 1rem;
}

.related-link {
  display: block;
  font-weight: 500;
  color: var(--color-text-accent);
  text-decoration: none;
  margin-bottom: 0.25rem;
}

.related-link:hover {
  text-decoration: underline;
}

.related-meta {
  font-size: 0.75rem;
  color: var(--color-text-tertiary);
}

.no-related {
  color: var(--color-text-tertiary);
  font-size: 0.875rem;
  font-style: italic;
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

.markdown-body pre code {
  background-color: transparent;
  padding: 0;
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

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .discussion-content-wrapper {
    grid-template-columns: 1fr;
  }
  
  .discussion-sidebar {
    order: -1;
    margin-bottom: 1.5rem;
  }
}
</style>

<template>
  <div class="discussions-page">
    <div class="discussions-header">
      <div class="discussions-title">
        <h1>Discussions</h1>
        <p class="subtitle">チームディスカッション</p>
      </div>
      
      <Button 
        variant="primary" 
        @click="navigateToCreate"
        class="create-discussion-btn"
      >
        <Icon name="mdi:plus" class="mr-1" />
        新規Discussion作成
      </Button>
    </div>

    <Card class="discussions-filters">
      <div class="filter-section">
        <div class="search-container">
          <Input 
            v-model="searchQuery" 
            placeholder="Discussionを検索..." 
            prepend-icon="mdi:magnify"
            @input="handleSearch"
          />
        </div>
        
        <div class="filter-options">
          <Select 
            v-model="categoryFilter" 
            :options="categoryOptions" 
            placeholder="カテゴリ"
            class="filter-select"
            @change="applyFilters"
          />
          
          <Select 
            v-model="authorFilter" 
            :options="authorOptions" 
            placeholder="作成者"
            class="filter-select"
            @change="applyFilters"
          />
        </div>
      </div>
    </Card>

    <div class="discussions-list">
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>読み込み中...</p>
      </div>
      
      <div v-else-if="filteredDiscussions.length === 0" class="empty-state">
        <Icon name="mdi:forum-outline" size="48" />
        <p>条件に一致するディスカッションはありません。</p>
        <Button variant="outline" @click="navigateToCreate">新規作成</Button>
      </div>
      
      <Card v-for="discussion in filteredDiscussions" :key="discussion.id" class="discussion-card" @click="navigateToDiscussion(discussion.id)">
        <div class="discussion-header">
          <div class="discussion-title-container">
            <h2 class="discussion-title">{{ discussion.title }}</h2>
            <div class="discussion-category" :class="`category-${discussion.category.color}`">
              {{ discussion.category.name }}
            </div>
          </div>
          
          <div class="discussion-meta">
            <span class="discussion-author">{{ discussion.author }}</span>
            <span class="discussion-date">{{ formatDate(discussion.created_at) }}</span>
          </div>
        </div>
        
        <div class="discussion-preview">
          {{ getContentPreview(discussion.content) }}
        </div>
        
        <div class="discussion-stats">
          <div class="stat-item">
            <Icon name="mdi:comment-outline" />
            <span>{{ discussion.replies_count }} 返信</span>
          </div>
          
          <div class="stat-item">
            <Icon name="mdi:eye-outline" />
            <span>{{ discussion.views_count }} 閲覧</span>
          </div>
          
          <div class="stat-item">
            <Icon name="mdi:calendar-clock" />
            <span>最終更新: {{ formatDate(discussion.last_activity_at) }}</span>
          </div>
        </div>
      </Card>
      
      <div v-if="pagination.totalPages > 1" class="pagination">
        <Button 
          variant="outline" 
          size="sm" 
          :disabled="pagination.currentPage <= 1"
          @click="handlePageChange(pagination.currentPage - 1)"
        >
          <Icon name="mdi:chevron-left" />
        </Button>
        
        <span class="pagination-info">
          {{ pagination.currentPage }} / {{ pagination.totalPages }} ページ
        </span>
        
        <Button 
          variant="outline" 
          size="sm" 
          :disabled="pagination.currentPage >= pagination.totalPages"
          @click="handlePageChange(pagination.currentPage + 1)"
        >
          <Icon name="mdi:chevron-right" />
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const discussions = ref([])
const searchQuery = ref('')

// フィルター状態
const categoryFilter = ref('')
const authorFilter = ref('')

// ページネーション
const pagination = ref({
  currentPage: 1,
  totalPages: 3,
  totalItems: 25,
  itemsPerPage: 10
})

// カテゴリオプション（デモデータ）
const categoryOptions = [
  { value: '', label: 'すべてのカテゴリ' },
  { value: '1', label: '全般' },
  { value: '2', label: '開発環境' },
  { value: '3', label: 'Q&A' },
  { value: '4', label: 'アイデア' },
  { value: '5', label: 'お知らせ' }
]

// 作成者オプション（デモデータ）
const authorOptions = [
  { value: '', label: 'すべての作成者' },
  { value: 'shimauma0312', label: 'shimauma0312' },
  { value: 'tanuki456', label: 'tanuki456' },
  { value: 'kitsune789', label: 'kitsune789' }
]

// フィルター適用したDiscussionリスト
const filteredDiscussions = computed(() => {
  let result = [...discussions.value]
  
  // 検索クエリでフィルタリング
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(discussion => 
      discussion.title.toLowerCase().includes(query) || 
      discussion.content.toLowerCase().includes(query)
    )
  }
  
  // カテゴリでフィルタリング
  if (categoryFilter.value) {
    result = result.filter(discussion => 
      discussion.category.id.toString() === categoryFilter.value
    )
  }
  
  // 作成者でフィルタリング
  if (authorFilter.value) {
    result = result.filter(discussion => 
      discussion.author === authorFilter.value
    )
  }
  
  return result
})

// 日付フォーマット
function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('ja-JP', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  }).format(date)
}

// 内容のプレビューを取得（最初の100文字）
function getContentPreview(content) {
  if (!content) return ''
  const plainText = content.replace(/<\/?[^>]+(>|$)/g, '') // HTMLタグを削除
  return plainText.length > 100 ? plainText.substring(0, 100) + '...' : plainText
}

// 検索処理
function handleSearch() {
  pagination.value.currentPage = 1
}

// フィルター適用
function applyFilters() {
  pagination.value.currentPage = 1
}

// ページ変更処理
function handlePageChange(page) {
  pagination.value.currentPage = page
  
  // APIから指定ページのデータを取得（将来実装）
  // fetchDiscussions()
}

// Discussion詳細画面へ遷移
function navigateToDiscussion(id) {
  router.push(`/discussions/${id}`)
}

// Discussion作成画面へ遷移
function navigateToCreate() {
  router.push('/discussions/create')
}

onMounted(async () => {
  loading.value = true
  
  try {
    // APIからデータを取得する（将来実装）
    // const config = useRuntimeConfig()
    // const response = await fetch(`${config.public.apiBaseUrl}/api/v1/discussions`)
    // discussions.value = await response.json()
    
    // デモデータ
    discussions.value = [
      {
        id: 1,
        title: '開発環境のセットアップについて',
        content: 'Dockerを使った開発環境のセットアップ方法についてディスカッションを開始します。環境構築の手順や注意点について共有しましょう。',
        author: 'shimauma0312',
        created_at: '2025-05-05',
        last_activity_at: '2025-06-01',
        replies_count: 12,
        views_count: 45,
        category: {
          id: 2,
          name: '開発環境',
          color: 'blue'
        }
      },
      {
        id: 2,
        title: 'チケット管理システムの運用ルールについて',
        content: 'チケット管理システムの効果的な運用方法について意見を募集します。特にIssueとPull Requestの連携について良い方法があれば教えてください。',
        author: 'tanuki456',
        created_at: '2025-05-15',
        last_activity_at: '2025-05-25',
        replies_count: 8,
        views_count: 32,
        category: {
          id: 1,
          name: '全般',
          color: 'gray'
        }
      },
      {
        id: 3,
        title: '新機能のアイデア: カンバンボード',
        content: 'チケット管理にカンバンボード形式のビューを追加してはどうでしょうか？タスクの進捗を視覚的に確認できると便利だと思います。',
        author: 'kitsune789',
        created_at: '2025-05-20',
        last_activity_at: '2025-05-30',
        replies_count: 15,
        views_count: 67,
        category: {
          id: 4,
          name: 'アイデア',
          color: 'green'
        }
      },
      {
        id: 4,
        title: 'APIのエラーハンドリングについて質問',
        content: 'バックエンドAPIのエラーハンドリングについて質問です。現在の実装では...',
        author: 'tanuki456',
        created_at: '2025-06-01',
        last_activity_at: '2025-06-05',
        replies_count: 3,
        views_count: 12,
        category: {
          id: 3,
          name: 'Q&A',
          color: 'purple'
        }
      },
      {
        id: 5,
        title: 'プロジェクトのマイルストーン計画',
        content: '今後のプロジェクト進行について、以下のマイルストーンを提案します...',
        author: 'shimauma0312',
        created_at: '2025-06-02',
        last_activity_at: '2025-06-07',
        replies_count: 7,
        views_count: 28,
        category: {
          id: 5,
          name: 'お知らせ',
          color: 'red'
        }
      }
    ]
  } catch (error) {
    console.error('APIからデータを取得できませんでした:', error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.discussions-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.discussions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.discussions-title h1 {
  font-size: 1.75rem;
  font-weight: 600;
  margin: 0;
}

.subtitle {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.discussions-filters {
  margin-bottom: 1.5rem;
}

.filter-section {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 1rem;
}

.search-container {
  flex: 1;
  min-width: 250px;
}

.filter-options {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.filter-select {
  width: 150px;
}

.discussions-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 2rem;
}

.loading-state, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 0;
  text-align: center;
  color: var(--color-text-secondary);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--color-border-primary);
  border-radius: 50%;
  border-top-color: var(--color-accent-primary);
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.empty-state p {
  margin: 1rem 0;
}

.discussion-card {
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.discussion-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.discussion-header {
  margin-bottom: 0.75rem;
}

.discussion-title-container {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-bottom: 0.5rem;
}

.discussion-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
}

.discussion-category {
  font-size: 0.75rem;
  font-weight: 500;
  padding: 0.2rem 0.5rem;
  border-radius: 1rem;
  text-transform: uppercase;
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

.discussion-meta {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
}

.discussion-author {
  font-weight: 500;
}

.discussion-date {
  margin-left: 0.5rem;
}

.discussion-preview {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin-bottom: 1rem;
  line-height: 1.5;
}

.discussion-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--color-border-primary);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.75rem;
  color: var(--color-text-tertiary);
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  margin-top: 1rem;
}

.pagination-info {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .discussions-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .create-discussion-btn {
    width: 100%;
  }
  
  .filter-section {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-container {
    width: 100%;
  }
  
  .filter-options {
    width: 100%;
  }
  
  .filter-select {
    flex: 1;
  }
  
  .discussion-stats {
    flex-direction: column;
    gap: 0.5rem;
  }
}
</style>

<template>
  <div class="issues-page">
    <div class="issues-header">
      <div class="issues-title">
        <h1>Issues</h1>
        <p class="subtitle">課題管理・追跡</p>
      </div>
      
      <Button 
        variant="primary" 
        @click="navigateToCreate"
        class="create-issue-btn"
      >
        <Icon name="mdi:plus" class="mr-1" />
        新規Issue作成
      </Button>
    </div>

    <Card class="issues-filters">
      <div class="filter-section">
        <div class="search-container">
          <Input 
            v-model="searchQuery" 
            placeholder="Issueを検索..." 
            prepend-icon="mdi:magnify"
            @input="handleSearch"
          />
        </div>
        
        <div class="filter-options">
          <Select 
            v-model="statusFilter" 
            :options="statusOptions" 
            placeholder="ステータス"
            class="filter-select"
            @change="applyFilters"
          />
          
          <Select 
            v-model="labelFilter" 
            :options="labelOptions" 
            placeholder="ラベル"
            class="filter-select"
            @change="applyFilters"
          />
          
          <Select 
            v-model="assigneeFilter" 
            :options="assigneeOptions" 
            placeholder="担当者"
            class="filter-select"
            @change="applyFilters"
          />
        </div>
      </div>
    </Card>

    <div class="issues-table-container">
      <Table
        :columns="columns"
        :items="filteredIssues"
        :loading="loading"
        :sortBy="sortBy"
        :sortDirection="sortDirection"
        :pagination="pagination"
        emptyText="条件に一致するIssueはありません。フィルタを調整するか、新しいIssueを作成してください。"
        @sort="handleSort"
        @row-click="handleIssueClick"
        @page-change="handlePageChange"
      >
        <!-- Issue タイトルとラベル -->
        <template #cell(title)="{ item }">
          <div class="issue-title-cell">
            <div class="issue-icon" :class="`status-${item.status}`">
              <Icon :name="getStatusIcon(item.status)" />
            </div>
            <div class="issue-title-content">
              <div class="issue-title-row">
                <span class="issue-title-text">{{ item.title }}</span>
                <span class="issue-number">#{{ item.id }}</span>
              </div>
              
              <div v-if="item.labels && item.labels.length" class="issue-labels">
                <Label 
                  v-for="label in item.labels" 
                  :key="label.id"
                  :name="label.name"
                  :color="label.color"
                  size="sm"
                />
              </div>
            </div>
          </div>
        </template>
        
        <!-- 作成情報 -->
        <template #cell(created)="{ item }">
          <div class="issue-created">
            {{ formatDate(item.created_at) }} に {{ item.created_by }} が作成
          </div>
        </template>
        
        <!-- 担当者 -->
        <template #cell(assignee)="{ item }">
          <div v-if="item.assignees && item.assignees.length" class="issue-assignees">
            <span v-for="assignee in item.assignees" :key="assignee.id" class="assignee-name">
              {{ assignee.name }}
            </span>
          </div>
          <div v-else class="no-assignee">未割り当て</div>
        </template>
        
        <!-- コメント数 -->
        <template #cell(comments)="{ item }">
          <div v-if="item.comments_count" class="issue-comments">
            <Icon name="mdi:comment-outline" />
            <span>{{ item.comments_count }}</span>
          </div>
        </template>
        
        <!-- 行アクション -->
        <template #rowActions="{ item }">
          <div class="row-actions">
            <button class="action-btn" @click.stop="showMoreActions(item)">
              <Icon name="mdi:dots-horizontal" />
            </button>
          </div>
        </template>
      </Table>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const issues = ref([])
const searchQuery = ref('')

// フィルター状態
const statusFilter = ref('')
const labelFilter = ref('')
const assigneeFilter = ref('')

// ソート状態
const sortBy = ref('updated_at')
const sortDirection = ref('desc')

// ページネーション
const pagination = ref({
  currentPage: 1,
  totalPages: 5,
  totalItems: 47,
  itemsPerPage: 10
})

// テーブルカラム定義
const columns = [
  { 
    key: 'title', 
    label: 'タイトル', 
    sortable: true,
    width: '60%'
  },
  { 
    key: 'created', 
    label: '作成情報',
    width: '20%'
  },
  { 
    key: 'assignee', 
    label: '担当者',
    width: '10%'
  },
  { 
    key: 'comments', 
    label: 'コメント',
    width: '10%',
    cellClass: 'text-center'
  }
]

// ステータスオプション
const statusOptions = [
  { value: '', label: 'すべてのステータス' },
  { value: 'open', label: '未対応' },
  { value: 'in-progress', label: '対応中' },
  { value: 'closed', label: '完了' }
]

// ラベルオプション（デモデータ）
const labelOptions = [
  { value: '', label: 'すべてのラベル' },
  { value: '1', label: 'バグ' },
  { value: '2', label: '機能追加' },
  { value: '3', label: 'ドキュメント' },
  { value: '4', label: '質問' },
  { value: '5', label: '優先度:高' },
  { value: '6', label: '優先度:中' },
  { value: '7', label: '優先度:低' }
]

// 担当者オプション（デモデータ）
const assigneeOptions = [
  { value: '', label: 'すべての担当者' },
  { value: '1', label: 'shimauma0312' },
  { value: '2', label: 'tanuki456' },
  { value: '3', label: 'kitsune789' }
]

// フィルター適用したIssueリスト
const filteredIssues = computed(() => {
  let result = [...issues.value]
  
  // 検索クエリでフィルタリング
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(issue => 
      issue.title.toLowerCase().includes(query) || 
      `#${issue.id}`.includes(query)
    )
  }
  
  // ステータスでフィルタリング
  if (statusFilter.value) {
    result = result.filter(issue => issue.status === statusFilter.value)
  }
  
  // ラベルでフィルタリング
  if (labelFilter.value) {
    result = result.filter(issue => 
      issue.labels && issue.labels.some(label => label.id.toString() === labelFilter.value)
    )
  }
  
  // 担当者でフィルタリング
  if (assigneeFilter.value) {
    result = result.filter(issue => 
      issue.assignees && issue.assignees.some(assignee => assignee.id.toString() === assigneeFilter.value)
    )
  }
  
  return result
})

// ステータスに応じたアイコンを返す
function getStatusIcon(status) {
  switch (status) {
    case 'open':
      return 'mdi:alert-circle-outline'
    case 'in-progress':
      return 'mdi:progress-wrench'
    case 'closed':
      return 'mdi:check-circle'
    default:
      return 'mdi:alert-circle-outline'
  }
}

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

// 検索処理
function handleSearch() {
  pagination.value.currentPage = 1
}

// フィルター適用
function applyFilters() {
  pagination.value.currentPage = 1
}

// ソート処理
function handleSort(sort) {
  sortBy.value = sort.key
  sortDirection.value = sort.direction
  
  // APIからソート済みデータを取得（将来実装）
  // fetchIssues()
}

// ページ変更処理
function handlePageChange(page) {
  pagination.value.currentPage = page
  
  // APIから指定ページのデータを取得（将来実装）
  // fetchIssues()
}

// Issue詳細画面へ遷移
function handleIssueClick({ item }) {
  router.push(`/issues/${item.id}`)
}

// Issue作成画面へ遷移
function navigateToCreate() {
  router.push('/issues/create')
}

// その他アクションメニュー表示
function showMoreActions(issue) {
  console.log('アクションメニュー表示:', issue)
  // 実装予定: ドロップダウンメニュー表示
}

onMounted(async () => {
  loading.value = true
  
  try {
    // APIからデータを取得する（将来実装）
    // const config = useRuntimeConfig()
    // const response = await fetch(`${config.public.apiBaseUrl}/api/v1/issues`)
    // issues.value = await response.json()
    
    // デモデータ
    issues.value = [
      {
        id: 1,
        title: 'アーキテクチャ図の改善',
        status: 'open',
        created_by: 'shimauma0312',
        created_at: '2025-05-01',
        updated_at: '2025-05-02',
        comments_count: 3,
        labels: [
          { id: 3, name: 'ドキュメント', color: '#0075ca' }
        ],
        assignees: [
          { id: 1, name: 'shimauma0312' }
        ]
      },
      {
        id: 2,
        title: 'プロジェクト初期構成とDockerセットアップ',
        status: 'in-progress',
        created_by: 'shimauma0312',
        created_at: '2025-05-15',
        updated_at: '2025-06-01',
        comments_count: 2,
        labels: [
          { id: 2, name: '機能追加', color: '#0366d6' },
          { id: 5, name: '優先度:高', color: '#e99695' }
        ],
        assignees: [
          { id: 1, name: 'shimauma0312' },
          { id: 2, name: 'tanuki456' }
        ]
      },
      {
        id: 3,
        title: 'フロントエンドのテスト環境構築',
        status: 'open',
        created_by: 'tanuki456',
        created_at: '2025-05-20',
        updated_at: '2025-05-21',
        comments_count: 1,
        labels: [
          { id: 2, name: '機能追加', color: '#0366d6' },
          { id: 6, name: '優先度:中', color: '#fbca04' }
        ],
        assignees: []
      },
      {
        id: 4,
        title: 'APIドキュメントの更新',
        status: 'closed',
        created_by: 'kitsune789',
        created_at: '2025-05-10',
        updated_at: '2025-05-18',
        comments_count: 4,
        labels: [
          { id: 3, name: 'ドキュメント', color: '#0075ca' },
          { id: 7, name: '優先度:低', color: '#c2e0c6' }
        ],
        assignees: [
          { id: 3, name: 'kitsune789' }
        ]
      },
      {
        id: 5,
        title: 'ナビゲーションメニューのUI改善',
        status: 'open',
        created_by: 'tanuki456',
        created_at: '2025-06-01',
        updated_at: '2025-06-02',
        comments_count: 0,
        labels: [
          { id: 2, name: '機能追加', color: '#0366d6' }
        ],
        assignees: [
          { id: 2, name: 'tanuki456' }
        ]
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
.issues-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.issues-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.issues-title h1 {
  font-size: 1.75rem;
  font-weight: 600;
  margin: 0;
}

.subtitle {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.issues-filters {
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

.issues-table-container {
  margin-bottom: 2rem;
}

.issue-title-cell {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
}

.issue-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.125rem;
  margin-top: 2px;
}

.status-open {
  color: #2da44e;
}

.status-in-progress {
  color: #bf8700;
}

.status-closed {
  color: #8250df;
}

.issue-title-content {
  flex: 1;
  min-width: 0;
}

.issue-title-row {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
  margin-bottom: 0.25rem;
}

.issue-title-text {
  font-weight: 500;
  color: var(--color-text-primary);
}

.issue-number {
  color: var(--color-text-tertiary);
  font-size: 0.875rem;
}

.issue-labels {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.issue-created, .no-assignee {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
}

.issue-assignees {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.assignee-name {
  font-size: 0.875rem;
}

.issue-comments {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.25rem;
  color: var(--color-text-secondary);
}

.row-actions {
  display: flex;
  justify-content: flex-end;
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

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .issues-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .create-issue-btn {
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
}
</style>

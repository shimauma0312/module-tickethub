<template>
  <div class="search-page">
    <div class="search-header">
      <h1>検索結果</h1>
      <div class="search-input-container">
        <Input
          v-model="searchQuery"
          placeholder="検索キーワードを入力..."
          prepend-icon="mdi:magnify"
          :clearable="true"
          @keyup.enter="performSearch"
        />
        <Button @click="performSearch" :loading="loading">検索</Button>
      </div>
    </div>

    <div class="search-content">
      <div class="search-filters">
        <Card>
          <div class="filters-container">
            <h3 class="filters-title">フィルター</h3>
            
            <div class="filter-section">
              <h4 class="filter-section-title">タイプ</h4>
              <div class="filter-options">
                <label class="filter-option">
                  <input type="checkbox" v-model="filters.types.issue" />
                  <span>Issue</span>
                  <span class="filter-count">{{ counts.issue || 0 }}</span>
                </label>
                <label class="filter-option">
                  <input type="checkbox" v-model="filters.types.discussion" />
                  <span>Discussion</span>
                  <span class="filter-count">{{ counts.discussion || 0 }}</span>
                </label>
                <label class="filter-option">
                  <input type="checkbox" v-model="filters.types.comment" />
                  <span>コメント</span>
                  <span class="filter-count">{{ counts.comment || 0 }}</span>
                </label>
                <label class="filter-option">
                  <input type="checkbox" v-model="filters.types.user" />
                  <span>ユーザー</span>
                  <span class="filter-count">{{ counts.user || 0 }}</span>
                </label>
              </div>
            </div>
            
            <div class="filter-section">
              <h4 class="filter-section-title">ステータス</h4>
              <div class="filter-options">
                <label class="filter-option">
                  <input type="checkbox" v-model="filters.status.open" />
                  <span>未対応</span>
                </label>
                <label class="filter-option">
                  <input type="checkbox" v-model="filters.status.inProgress" />
                  <span>対応中</span>
                </label>
                <label class="filter-option">
                  <input type="checkbox" v-model="filters.status.closed" />
                  <span>完了</span>
                </label>
              </div>
            </div>
            
            <div class="filter-section">
              <h4 class="filter-section-title">期間</h4>
              <div class="filter-options">
                <select v-model="filters.timeRange" class="time-range-select">
                  <option value="any">指定なし</option>
                  <option value="day">24時間以内</option>
                  <option value="week">1週間以内</option>
                  <option value="month">1ヶ月以内</option>
                  <option value="year">1年以内</option>
                </select>
              </div>
            </div>
            
            <div class="filter-section">
              <h4 class="filter-section-title">ラベル</h4>
              <div class="filter-options">
                <div class="labels-filter">
                  <div v-if="selectedLabels.length" class="selected-labels">
                    <span 
                      v-for="label in selectedLabels" 
                      :key="label.id"
                      class="selected-label"
                      :style="{ backgroundColor: label.color + '33', color: label.color }"
                    >
                      {{ label.name }}
                      <button class="remove-label" @click="removeLabel(label.id)">
                        <Icon name="mdi:close" />
                      </button>
                    </span>
                  </div>
                  
                  <select v-model="labelFilter" class="label-select" @change="addLabel">
                    <option value="">ラベルを選択...</option>
                    <option 
                      v-for="label in availableLabels" 
                      :key="label.id" 
                      :value="label.id"
                    >
                      {{ label.name }}
                    </option>
                  </select>
                </div>
              </div>
            </div>
            
            <div class="filter-actions">
              <Button variant="outline" size="sm" @click="resetFilters">
                フィルターをリセット
              </Button>
              <Button variant="primary" size="sm" @click="applyFilters">
                適用
              </Button>
            </div>
          </div>
        </Card>
      </div>

      <div class="search-results">
        <!-- 読み込み中 -->
        <div v-if="loading" class="search-loading">
          <div class="loading-spinner"></div>
          <p>検索結果を読み込み中...</p>
        </div>

        <!-- 検索結果なし -->
        <Card v-else-if="!results.length && searchPerformed" class="empty-results">
          <div class="empty-state">
            <Icon name="mdi:file-search-outline" size="48" />
            <h3>検索結果はありません</h3>
            <p>別のキーワードで検索するか、フィルターを調整してみてください</p>
          </div>
        </Card>

        <!-- 検索前 -->
        <Card v-else-if="!searchPerformed" class="empty-results">
          <div class="empty-state">
            <Icon name="mdi:magnify" size="48" />
            <h3>検索を開始</h3>
            <p>キーワードを入力して検索してください</p>
          </div>
        </Card>

        <!-- 検索結果リスト -->
        <div v-else>
          <div class="results-header">
            <div class="results-count">
              {{ totalResults }} 件の検索結果
            </div>
            <div class="sort-options">
              <label>
                並び替え:
                <select v-model="sortOption" @change="applySort">
                  <option value="relevance">関連性</option>
                  <option value="newest">最新順</option>
                  <option value="oldest">古い順</option>
                </select>
              </label>
            </div>
          </div>

          <div class="results-list">
            <Card 
              v-for="result in results" 
              :key="result.id"
              class="result-card"
              hoverable
              @click="navigateToResult(result)"
            >
              <div class="result-content">
                <div class="result-header">
                  <div class="result-type">
                    <Icon :name="getResultTypeIcon(result.type)" />
                    <span>{{ getResultTypeName(result.type) }}</span>
                    <span v-if="result.status" :class="['result-status', `status-${result.status}`]">
                      {{ getStatusName(result.status) }}
                    </span>
                  </div>
                  <div class="result-meta">
                    {{ formatDate(result.updated_at || result.created_at) }}
                  </div>
                </div>
                
                <div class="result-title">
                  <span v-html="highlightText(result.title, searchQuery)"></span>
                </div>
                
                <div class="result-excerpt" v-if="result.excerpt">
                  <p v-html="highlightText(result.excerpt, searchQuery)"></p>
                </div>
                
                <div class="result-footer">
                  <div class="result-tags" v-if="result.labels && result.labels.length">
                    <Label 
                      v-for="label in result.labels.slice(0, 3)" 
                      :key="label.id"
                      :name="label.name"
                      :color="label.color"
                      size="sm"
                    />
                    <span v-if="result.labels.length > 3" class="more-labels">
                      +{{ result.labels.length - 3 }}
                    </span>
                  </div>
                  
                  <div class="result-user" v-if="result.user">
                    <span class="user-avatar">{{ getUserInitials(result.user.name) }}</span>
                    <span class="user-name">{{ result.user.name }}</span>
                  </div>
                </div>
              </div>
            </Card>
          </div>

          <!-- ページネーション -->
          <div v-if="totalPages > 1" class="pagination">
            <button 
              class="pagination-button"
              :disabled="currentPage <= 1"
              @click="goToPage(currentPage - 1)"
            >
              <Icon name="mdi:chevron-left" />
            </button>
            
            <div class="pagination-pages">
              <button 
                v-for="page in displayedPages" 
                :key="page"
                class="page-button"
                :class="{ active: page === currentPage }"
                @click="goToPage(page)"
              >
                {{ page }}
              </button>
            </div>
            
            <button 
              class="pagination-button"
              :disabled="currentPage >= totalPages"
              @click="goToPage(currentPage + 1)"
            >
              <Icon name="mdi:chevron-right" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();

// 検索状態
const searchQuery = ref('');
const searchPerformed = ref(false);
const loading = ref(false);
const error = ref(null);

// 検索結果
const results = ref([]);
const totalResults = ref(0);

// フィルター
const filters = ref({
  types: {
    issue: true,
    discussion: true,
    comment: true,
    user: true
  },
  status: {
    open: true,
    inProgress: true,
    closed: true
  },
  timeRange: 'any',
  labels: []
});

// ラベル管理
const labelFilter = ref('');
const selectedLabels = ref([]);
const availableLabels = ref([
  { id: '1', name: 'バグ', color: '#d73a49' },
  { id: '2', name: '機能追加', color: '#0366d6' },
  { id: '3', name: 'ドキュメント', color: '#0075ca' },
  { id: '4', name: '質問', color: '#d876e3' },
  { id: '5', name: '優先度:高', color: '#e99695' },
  { id: '6', name: '優先度:中', color: '#fbca04' },
  { id: '7', name: '優先度:低', color: '#c2e0c6' }
]);

// 結果タイプの集計
const counts = ref({
  issue: 24,
  discussion: 8,
  comment: 42,
  user: 3
});

// ページネーション
const currentPage = ref(1);
const itemsPerPage = ref(10);
const totalPages = computed(() => Math.ceil(totalResults.value / itemsPerPage.value));

// 表示するページ番号
const displayedPages = computed(() => {
  const range = [];
  const delta = 2;
  const left = currentPage.value - delta;
  const right = currentPage.value + delta + 1;
  
  for (let i = 1; i <= totalPages.value; i++) {
    if (i === 1 || i === totalPages.value || (i >= left && i < right)) {
      range.push(i);
    }
  }
  
  return range;
});

// 並び替え
const sortOption = ref('relevance');

// URLからクエリとフィルター情報を取得
onMounted(() => {
  if (route.query.q) {
    searchQuery.value = route.query.q;
    performSearch();
  }
  
  // その他のフィルターもURLから復元する（実装省略）
});

// クエリ変更時のハンドラ
watch(searchQuery, () => {
  // 検索ボックスが空になったら結果をクリア
  if (!searchQuery.value.trim()) {
    results.value = [];
    searchPerformed.value = false;
  }
});

// ラベル追加
function addLabel() {
  if (!labelFilter.value) return;
  
  const label = availableLabels.value.find(l => l.id === labelFilter.value);
  if (label && !selectedLabels.value.some(l => l.id === label.id)) {
    selectedLabels.value.push(label);
    filters.value.labels.push(label.id);
  }
  
  labelFilter.value = '';
}

// ラベル削除
function removeLabel(labelId) {
  selectedLabels.value = selectedLabels.value.filter(l => l.id !== labelId);
  filters.value.labels = filters.value.labels.filter(id => id !== labelId);
}

// フィルターリセット
function resetFilters() {
  filters.value = {
    types: {
      issue: true,
      discussion: true,
      comment: true,
      user: true
    },
    status: {
      open: true,
      inProgress: true,
      closed: true
    },
    timeRange: 'any',
    labels: []
  };
  
  selectedLabels.value = [];
  labelFilter.value = '';
}

// フィルター適用
function applyFilters() {
  performSearch();
}

// 並び替え適用
function applySort() {
  performSearch();
}

// ページ移動
function goToPage(page) {
  currentPage.value = page;
  performSearch();
}

// 検索実行
async function performSearch() {
  if (!searchQuery.value.trim()) return;
  
  loading.value = true;
  error.value = null;
  searchPerformed.value = true;
  
  try {
    // URLのクエリパラメータを更新
    router.push({
      query: { 
        ...route.query,
        q: searchQuery.value,
        page: currentPage.value
      }
    });
    
    // 実際のAPI呼び出しはバックエンドの準備ができてから実装
    // 現在はモックデータで対応
    await new Promise(resolve => setTimeout(resolve, 1500)); // API呼び出しの遅延をシミュレート
    
    // モック検索結果データ
    const mockResults = [
      {
        id: 1,
        type: 'issue',
        title: 'アーキテクチャ図の改善',
        excerpt: 'アーキテクチャ図が古くなっているので、最新の構成を反映した図を作成する必要があります。特にWebSocketの部分が...',
        status: 'open',
        created_at: '2025-05-01T09:23:45Z',
        updated_at: '2025-05-02T14:10:20Z',
        user: {
          id: 1,
          name: 'shimauma0312'
        },
        labels: [
          { id: '3', name: 'ドキュメント', color: '#0075ca' }
        ]
      },
      {
        id: 2,
        type: 'issue',
        title: 'プロジェクト初期構成とDockerセットアップ',
        excerpt: 'プロジェクトの初期構成を作成し、Docker環境をセットアップします。これには、フロントエンド（Nuxt）とバックエンド（Go）の開発環境が含まれます。',
        status: 'in-progress',
        created_at: '2025-05-15T11:42:18Z',
        updated_at: '2025-06-01T08:15:33Z',
        user: {
          id: 1,
          name: 'shimauma0312'
        },
        labels: [
          { id: '2', name: '機能追加', color: '#0366d6' },
          { id: '5', name: '優先度:高', color: '#e99695' }
        ]
      },
      {
        id: 1,
        type: 'discussion',
        title: 'TicketHubの開発方針について',
        excerpt: 'TicketHubの今後の開発方針について議論したいと思います。現在の機能のうち、特に重視すべき点と、追加機能の優先順位について...',
        created_at: '2025-04-28T10:15:00Z',
        updated_at: '2025-05-05T16:42:30Z',
        user: {
          id: 1,
          name: 'shimauma0312'
        },
        labels: [
          { id: '4', name: '質問', color: '#d876e3' }
        ]
      },
      {
        id: 12,
        type: 'comment',
        title: 'アーキテクチャ図の改善 に対するコメント',
        excerpt: 'アーキテクチャ図の改善案を作成しました。特にWebSocketの部分を詳細化し、スケーリング時の構成も追加しています。レビューをお願いします。',
        parent_id: 1,
        parent_type: 'issue',
        created_at: '2025-05-02T13:45:22Z',
        user: {
          id: 2,
          name: 'tanuki456'
        }
      },
      {
        id: 3,
        type: 'user',
        title: 'tanuki456',
        excerpt: 'フロントエンド開発者。TicketHubプロジェクトでUI/UXの設計と実装を担当。',
        created_at: '2024-12-01T09:00:00Z'
      }
    ];
    
    results.value = mockResults;
    totalResults.value = 77; // 仮の総件数
    
  } catch (error) {
    console.error('検索中にエラーが発生しました:', error);
    error.value = '検索処理に失敗しました。後でもう一度お試しください。';
  } finally {
    loading.value = false;
  }
}

// 検索結果に遷移
function navigateToResult(result) {
  let path;
  
  switch (result.type) {
    case 'issue':
      path = `/issues/${result.id}`;
      break;
    case 'discussion':
      path = `/discussions/${result.id}`;
      break;
    case 'comment':
      // コメントの場合は親アイテムに移動し、コメント位置にスクロール
      path = `/${result.parent_type}s/${result.parent_id}#comment-${result.id}`;
      break;
    case 'user':
      path = `/users/${result.title}`;
      break;
    default:
      path = '/';
  }
  
  router.push(path);
}

// タイプに応じたアイコンを取得
function getResultTypeIcon(type) {
  switch (type) {
    case 'issue':
      return 'mdi:alert-circle-outline';
    case 'discussion':
      return 'mdi:forum-outline';
    case 'comment':
      return 'mdi:comment-text-outline';
    case 'user':
      return 'mdi:account';
    default:
      return 'mdi:file-document-outline';
  }
}

// タイプ名を取得
function getResultTypeName(type) {
  switch (type) {
    case 'issue':
      return 'Issue';
    case 'discussion':
      return 'Discussion';
    case 'comment':
      return 'コメント';
    case 'user':
      return 'ユーザー';
    default:
      return 'ドキュメント';
  }
}

// ステータス名を取得
function getStatusName(status) {
  switch (status) {
    case 'open':
      return '未対応';
    case 'in-progress':
      return '対応中';
    case 'closed':
      return '完了';
    default:
      return status;
  }
}

// ユーザーのイニシャルを取得
function getUserInitials(name) {
  if (!name) return '';
  return name.charAt(0).toUpperCase();
}

// 日付フォーマット
function formatDate(dateStr) {
  if (!dateStr) return '';
  
  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now - date;
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
  
  if (diffDays > 30) {
    return new Intl.DateTimeFormat('ja-JP', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    }).format(date);
  } else if (diffDays > 0) {
    return `${diffDays}日前`;
  } else {
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    if (diffHours > 0) {
      return `${diffHours}時間前`;
    } else {
      const diffMinutes = Math.floor(diffMs / (1000 * 60));
      if (diffMinutes > 0) {
        return `${diffMinutes}分前`;
      } else {
        return '数秒前';
      }
    }
  }
}

// テキスト内のキーワードをハイライト
function highlightText(text, query) {
  if (!text || !query) return text;
  
  const words = query.split(/\s+/).filter(word => word.length > 0);
  let highlightedText = text;
  
  words.forEach(word => {
    const regex = new RegExp(`(${word})`, 'gi');
    highlightedText = highlightedText.replace(regex, '<mark>$1</mark>');
  });
  
  return highlightedText;
}
</script>

<style scoped>
.search-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.search-header {
  margin-bottom: 2rem;
}

.search-header h1 {
  font-size: 1.75rem;
  font-weight: 600;
  margin: 0 0 1rem 0;
}

.search-input-container {
  display: flex;
  gap: 0.75rem;
}

.search-input-container :deep(input) {
  font-size: 1rem;
}

.search-content {
  display: flex;
  gap: 2rem;
}

.search-filters {
  width: 280px;
  flex-shrink: 0;
}

.search-results {
  flex: 1;
  min-width: 0;
}

.filters-container {
  padding: 0.5rem;
}

.filters-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0 0 1rem 0;
}

.filter-section {
  margin-bottom: 1.5rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid var(--color-border-primary);
}

.filter-section:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.filter-section-title {
  font-size: 0.875rem;
  font-weight: 600;
  margin: 0 0 0.75rem 0;
  color: var(--color-text-secondary);
}

.filter-options {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.filter-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
}

.filter-count {
  margin-left: auto;
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  background-color: var(--color-bg-secondary);
  padding: 0.125rem 0.375rem;
  border-radius: 10px;
}

.time-range-select, .label-select {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  background-color: var(--color-bg-canvas);
  font-size: 0.875rem;
  color: var(--color-text-primary);
}

.selected-labels {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.selected-label {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  border-radius: 10px;
  font-size: 0.75rem;
  font-weight: 500;
}

.remove-label {
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  font-size: 0.75rem;
}

.filter-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 1.5rem;
  gap: 0.5rem;
}

.search-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 0;
  gap: 1rem;
}

.loading-spinner {
  width: 2rem;
  height: 2rem;
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-top-color: var(--color-accent-primary);
  border-radius: 50%;
  animation: spin 1s ease-in-out infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.empty-results {
  padding: 3rem 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 2rem 1rem;
}

.empty-state h3 {
  margin: 1rem 0 0.5rem;
  font-size: 1.125rem;
  font-weight: 500;
}

.empty-state p {
  margin: 0;
  color: var(--color-text-secondary);
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.results-count {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
}

.sort-options select {
  padding: 0.375rem 0.5rem;
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  background-color: var(--color-bg-canvas);
  font-size: 0.875rem;
  color: var(--color-text-primary);
}

.results-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 2rem;
}

.result-card {
  padding: 0;
  overflow: hidden;
}

.result-content {
  padding: 1rem;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.result-type {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.result-status {
  padding: 0.125rem 0.375rem;
  border-radius: 10px;
  font-size: 0.75rem;
  font-weight: 500;
}

.status-open {
  background-color: rgba(45, 164, 78, 0.1);
  color: #2da44e;
}

.status-in-progress {
  background-color: rgba(191, 135, 0, 0.1);
  color: #bf8700;
}

.status-closed {
  background-color: rgba(130, 80, 223, 0.1);
  color: #8250df;
}

.result-meta {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.result-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: var(--color-text-primary);
}

.result-excerpt {
  margin-bottom: 1rem;
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.5;
}

.result-excerpt p {
  margin: 0;
}

.result-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.result-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.more-labels {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.result-user {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.user-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
  background-color: var(--color-accent-primary);
  color: white;
  border-radius: 50%;
  font-size: 0.75rem;
  font-weight: 500;
}

.user-name {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  margin-top: 2rem;
}

.pagination-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border-radius: var(--border-radius-md);
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
  border: none;
  cursor: pointer;
  transition: var(--transition-base);
}

.pagination-button:hover:not(:disabled) {
  background-color: var(--color-border-primary);
}

.pagination-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-pages {
  display: flex;
  gap: 0.25rem;
}

.page-button {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 2rem;
  height: 2rem;
  padding: 0 0.5rem;
  border-radius: var(--border-radius-md);
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
  border: none;
  cursor: pointer;
  transition: var(--transition-base);
  font-size: 0.875rem;
}

.page-button:hover:not(.active) {
  background-color: var(--color-border-primary);
}

.page-button.active {
  background-color: var(--color-accent-primary);
  color: white;
  font-weight: 500;
}

:deep(mark) {
  background-color: rgba(255, 234, 127, 0.5);
  padding: 0.125rem 0;
  border-radius: 2px;
}

/* レスポンシブデザイン */
@media (max-width: 992px) {
  .search-content {
    flex-direction: column;
  }
  
  .search-filters {
    width: 100%;
    margin-bottom: 1.5rem;
  }
}

@media (max-width: 768px) {
  .search-input-container {
    flex-direction: column;
  }
  
  .results-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .result-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }
}
</style>

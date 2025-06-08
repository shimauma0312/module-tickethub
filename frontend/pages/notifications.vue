<template>
  <div class="notifications-page">
    <div class="notifications-header">
      <h1>通知センター</h1>
      <div class="header-actions">
        <Button 
          v-if="hasUnreadNotifications" 
          variant="outline" 
          size="sm"
          @click="markAllAsRead" 
          :loading="markingAllAsRead"
        >
          すべて既読にする
        </Button>
        <Button 
          variant="outline" 
          size="sm"
          @click="refreshNotifications" 
          :loading="loading"
        >
          <Icon name="mdi:refresh" />
          更新
        </Button>
      </div>
    </div>

    <div class="notifications-filters">
      <Card>
        <div class="filters-container">
          <div class="filter-tabs">
            <button 
              class="filter-tab" 
              :class="{ 'active': filter === 'all' }"
              @click="filter = 'all'"
            >
              すべて
            </button>
            <button 
              class="filter-tab" 
              :class="{ 'active': filter === 'unread' }"
              @click="filter = 'unread'"
            >
              未読
              <span v-if="unreadCount > 0" class="filter-count">{{ unreadCount }}</span>
            </button>
            <button 
              class="filter-tab" 
              :class="{ 'active': filter === 'mentions' }"
              @click="filter = 'mentions'"
            >
              メンション
            </button>
            <button 
              class="filter-tab" 
              :class="{ 'active': filter === 'assignments' }"
              @click="filter = 'assignments'"
            >
              割り当て
            </button>
            <button 
              class="filter-tab" 
              :class="{ 'active': filter === 'system' }"
              @click="filter = 'system'"
            >
              システム
            </button>
          </div>
          
          <div class="search-filter">
            <Input 
              v-model="searchQuery" 
              placeholder="通知を検索..." 
              prepend-icon="mdi:magnify"
              size="sm"
            />
          </div>
        </div>
      </Card>
    </div>

    <div class="notifications-content">
      <!-- 読み込み中 -->
      <div v-if="loading && !notifications.length" class="notifications-loading">
        <div class="loading-spinner"></div>
        <p>通知を読み込み中...</p>
      </div>

      <!-- 通知なし -->
      <Card v-else-if="!filteredNotifications.length" class="empty-notifications">
        <div class="empty-state">
          <Icon name="mdi:bell-off-outline" size="48" />
          <h3>{{ getEmptyStateMessage() }}</h3>
          <p>新しい通知が届くとここに表示されます</p>
        </div>
      </Card>

      <!-- 通知リスト -->
      <div v-else class="notification-list">
        <Card 
          v-for="notification in filteredNotifications" 
          :key="notification.id"
          class="notification-card"
          :class="{ 'unread': !notification.read_at }"
          hoverable
        >
          <div class="notification-content">
            <div class="notification-icon-container">
              <Icon :name="getNotificationIcon(notification.type)" class="notification-icon" />
            </div>
            
            <div class="notification-details">
              <div class="notification-header">
                <span class="notification-sender">{{ notification.sender }}</span>
                <span class="notification-time">{{ formatTime(notification.created_at) }}</span>
              </div>
              
              <div class="notification-message" v-html="notification.message"></div>
              
              <div v-if="notification.target" class="notification-target">
                <NuxtLink :to="getTargetLink(notification)" class="target-link">
                  <span class="target-icon">
                    <Icon :name="getTargetIcon(notification.target_type)" />
                  </span>
                  <span class="target-title">{{ notification.target.title }}</span>
                </NuxtLink>
              </div>
            </div>
          </div>
          
          <div class="notification-actions">
            <button 
              v-if="!notification.read_at" 
              class="action-button read-button"
              @click="markAsRead(notification.id)"
              :disabled="markingAsRead === notification.id"
            >
              <Icon name="mdi:check" />
              <span class="action-text">既読</span>
            </button>
            <button 
              class="action-button"
              @click="deleteNotification(notification.id)"
              :disabled="deleting === notification.id"
            >
              <Icon name="mdi:delete-outline" />
              <span class="action-text">削除</span>
            </button>
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
        
        <span class="pagination-info">{{ currentPage }} / {{ totalPages }}</span>
        
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
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useAuthStore } from '~/stores/auth';

const authStore = useAuthStore();

// 通知データ
const notifications = ref([]);
const loading = ref(false);
const error = ref(null);

// フィルターと検索
const filter = ref('all');
const searchQuery = ref('');

// ページネーション
const currentPage = ref(1);
const itemsPerPage = 10;
const totalItems = ref(0);

// アクション状態
const markingAsRead = ref(null);
const markingAllAsRead = ref(false);
const deleting = ref(null);

// 計算されたプロパティ
const unreadCount = computed(() => {
  return notifications.value.filter(n => !n.read_at).length;
});

const hasUnreadNotifications = computed(() => unreadCount.value > 0);

const totalPages = computed(() => {
  return Math.ceil(totalItems.value / itemsPerPage);
});

const filteredNotifications = computed(() => {
  let result = [...notifications.value];
  
  // フィルター適用
  if (filter.value === 'unread') {
    result = result.filter(n => !n.read_at);
  } else if (filter.value === 'mentions') {
    result = result.filter(n => n.type === 'mention');
  } else if (filter.value === 'assignments') {
    result = result.filter(n => n.type === 'assignment');
  } else if (filter.value === 'system') {
    result = result.filter(n => n.type === 'system');
  }
  
  // 検索クエリ適用
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(n => 
      n.message.toLowerCase().includes(query) || 
      (n.sender && n.sender.toLowerCase().includes(query)) ||
      (n.target && n.target.title && n.target.title.toLowerCase().includes(query))
    );
  }
  
  return result;
});

// フィルター変更時にページをリセット
watch(filter, () => {
  currentPage.value = 1;
});

// 検索クエリ変更時にページをリセット
watch(searchQuery, () => {
  currentPage.value = 1;
});

// 通知を取得
onMounted(async () => {
  if (authStore.isAuthenticated) {
    await fetchNotifications();
  }
});

async function fetchNotifications() {
  loading.value = true;
  error.value = null;
  
  try {
    // 実際のAPI呼び出しはバックエンドの準備ができてから実装
    // 現在はモックデータで対応
    await new Promise(resolve => setTimeout(resolve, 1000)); // API呼び出しの遅延をシミュレート
    
    // モック通知データ
    const mockNotifications = [
      {
        id: 1,
        type: 'mention',
        sender: 'tanuki456',
        message: '<strong>@tanuki456</strong>があなたをメンションしました',
        created_at: new Date(Date.now() - 1 * 60 * 60 * 1000), // 1時間前
        read_at: null,
        target: {
          type: 'issue',
          id: 2,
          title: 'プロジェクト初期構成とDockerセットアップ'
        }
      },
      {
        id: 2,
        type: 'assignment',
        sender: 'shimauma0312',
        message: '<strong>@shimauma0312</strong>があなたをIssueに割り当てました',
        created_at: new Date(Date.now() - 3 * 60 * 60 * 1000), // 3時間前
        read_at: null,
        target: {
          type: 'issue',
          id: 3,
          title: 'フロントエンドのテスト環境構築'
        }
      },
      {
        id: 3,
        type: 'comment',
        sender: 'kitsune789',
        message: '<strong>@kitsune789</strong>があなたのIssueにコメントしました',
        created_at: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000), // 1日前
        read_at: new Date(Date.now() - 12 * 60 * 60 * 1000), // 12時間前
        target: {
          type: 'issue',
          id: 5,
          title: 'ナビゲーションメニューのUI改善'
        }
      },
      {
        id: 4,
        type: 'status',
        sender: 'shimauma0312',
        message: '<strong>@shimauma0312</strong>がIssueのステータスを「未対応」から「対応中」に変更しました',
        created_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2日前
        read_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2日前
        target: {
          type: 'issue',
          id: 2,
          title: 'プロジェクト初期構成とDockerセットアップ'
        }
      },
      {
        id: 5,
        type: 'system',
        sender: 'System',
        message: 'システムメンテナンスのお知らせ：2025年7月15日 23:00 - 2025年7月16日 01:00',
        created_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000), // 5日前
        read_at: null,
        target: null
      }
    ];
    
    notifications.value = mockNotifications;
    totalItems.value = mockNotifications.length;
    
  } catch (error) {
    console.error('通知の取得中にエラーが発生しました:', error);
    error.value = '通知の読み込みに失敗しました。後でもう一度お試しください。';
  } finally {
    loading.value = false;
  }
}

// 通知を更新
async function refreshNotifications() {
  await fetchNotifications();
}

// 通知を既読にする
async function markAsRead(notificationId) {
  markingAsRead.value = notificationId;
  
  try {
    // 実際のAPI呼び出しはバックエンドの準備ができてから実装
    // 現在はモックデータで対応
    await new Promise(resolve => setTimeout(resolve, 500)); // API呼び出しの遅延をシミュレート
    
    // 成功後の処理
    const notification = notifications.value.find(n => n.id === notificationId);
    if (notification) {
      notification.read_at = new Date();
    }
    
  } catch (error) {
    console.error('通知の既読処理中にエラーが発生しました:', error);
  } finally {
    markingAsRead.value = null;
  }
}

// すべての通知を既読にする
async function markAllAsRead() {
  markingAllAsRead.value = true;
  
  try {
    // 実際のAPI呼び出しはバックエンドの準備ができてから実装
    // 現在はモックデータで対応
    await new Promise(resolve => setTimeout(resolve, 1000)); // API呼び出しの遅延をシミュレート
    
    // 成功後の処理
    notifications.value.forEach(notification => {
      if (!notification.read_at) {
        notification.read_at = new Date();
      }
    });
    
  } catch (error) {
    console.error('すべての通知の既読処理中にエラーが発生しました:', error);
  } finally {
    markingAllAsRead.value = false;
  }
}

// 通知を削除する
async function deleteNotification(notificationId) {
  deleting.value = notificationId;
  
  try {
    // 実際のAPI呼び出しはバックエンドの準備ができてから実装
    // 現在はモックデータで対応
    await new Promise(resolve => setTimeout(resolve, 500)); // API呼び出しの遅延をシミュレート
    
    // 成功後の処理
    notifications.value = notifications.value.filter(n => n.id !== notificationId);
    totalItems.value = notifications.value.length;
    
  } catch (error) {
    console.error('通知の削除中にエラーが発生しました:', error);
  } finally {
    deleting.value = null;
  }
}

// ページ移動
function goToPage(page) {
  currentPage.value = page;
}

// 通知タイプに応じたアイコンを取得
function getNotificationIcon(type) {
  switch (type) {
    case 'mention':
      return 'mdi:at';
    case 'assignment':
      return 'mdi:account-check';
    case 'comment':
      return 'mdi:comment-text-outline';
    case 'status':
      return 'mdi:sync';
    case 'system':
      return 'mdi:bell-ring-outline';
    default:
      return 'mdi:bell-outline';
  }
}

// ターゲットタイプに応じたアイコンを取得
function getTargetIcon(type) {
  switch (type) {
    case 'issue':
      return 'mdi:alert-circle-outline';
    case 'discussion':
      return 'mdi:forum-outline';
    case 'pull-request':
      return 'mdi:source-pull';
    default:
      return 'mdi:link-variant';
  }
}

// ターゲットへのリンクを取得
function getTargetLink(notification) {
  if (!notification.target) return '#';
  
  switch (notification.target.type) {
    case 'issue':
      return `/issues/${notification.target.id}`;
    case 'discussion':
      return `/discussions/${notification.target.id}`;
    case 'pull-request':
      return `/pull-requests/${notification.target.id}`;
    default:
      return '#';
  }
}

// 空の状態メッセージを取得
function getEmptyStateMessage() {
  switch (filter.value) {
    case 'unread':
      return '未読の通知はありません';
    case 'mentions':
      return 'メンション通知はありません';
    case 'assignments':
      return '割り当て通知はありません';
    case 'system':
      return 'システム通知はありません';
    default:
      return '通知はありません';
  }
}

// 時間フォーマット
function formatTime(dateStr) {
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
</script>

<style scoped>
.notifications-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.notifications-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.notifications-header h1 {
  font-size: 1.75rem;
  font-weight: 600;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
}

.notifications-filters {
  margin-bottom: 1.5rem;
}

.filters-container {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 1rem;
}

.filter-tabs {
  display: flex;
  gap: 0.25rem;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  padding-bottom: 0.25rem;
}

.filter-tab {
  background: none;
  border: none;
  padding: 0.5rem 0.75rem;
  border-radius: var(--border-radius-md);
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  cursor: pointer;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: var(--transition-base);
}

.filter-tab.active {
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
  font-weight: 500;
}

.filter-tab:hover:not(.active) {
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
}

.filter-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 1.25rem;
  height: 1.25rem;
  padding: 0 0.25rem;
  background-color: var(--color-accent-primary);
  color: white;
  border-radius: 1rem;
  font-size: 0.75rem;
  font-weight: 500;
}

.search-filter {
  flex: 1;
  min-width: 200px;
  max-width: 300px;
}

.notifications-loading {
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

.empty-notifications {
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

.notification-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 2rem;
}

.notification-card {
  padding: 0;
  overflow: hidden;
}

.notification-card.unread {
  border-left: 3px solid var(--color-accent-primary);
}

.notification-content {
  display: flex;
  padding: 1rem;
}

.notification-icon-container {
  padding-right: 1rem;
}

.notification-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  background-color: var(--color-bg-secondary);
  border-radius: 50%;
  color: var(--color-accent-primary);
  font-size: 1.25rem;
}

.notification-details {
  flex: 1;
  min-width: 0;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.notification-sender {
  font-weight: 500;
}

.notification-time {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.notification-message {
  margin-bottom: 0.75rem;
  line-height: 1.4;
}

.notification-message :deep(strong) {
  font-weight: 600;
}

.notification-target {
  margin-top: 0.5rem;
}

.target-link {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.375rem 0.75rem;
  background-color: var(--color-bg-secondary);
  border-radius: var(--border-radius-md);
  color: var(--color-text-primary);
  text-decoration: none;
  font-size: 0.875rem;
  transition: var(--transition-base);
}

.target-link:hover {
  background-color: var(--color-border-primary);
}

.target-icon {
  color: var(--color-text-secondary);
}

.notification-actions {
  display: flex;
  border-top: 1px solid var(--color-border-primary);
}

.action-button {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.5rem;
  background: none;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: var(--transition-base);
}

.action-button:hover {
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
}

.action-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.read-button {
  color: var(--color-accent-primary);
}

.read-button:hover {
  color: var(--color-accent-primary);
}

.action-text {
  font-size: 0.875rem;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  margin-top: 1.5rem;
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

.pagination-info {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .notifications-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .header-actions {
    width: 100%;
  }
  
  .header-actions button {
    flex: 1;
  }
  
  .filters-container {
    flex-direction: column;
    gap: 1rem;
  }
  
  .search-filter {
    max-width: 100%;
  }
  
  .notification-content {
    flex-direction: column;
    gap: 1rem;
  }
  
  .notification-icon-container {
    padding-right: 0;
  }
  
  .notification-icon {
    width: 2rem;
    height: 2rem;
    font-size: 1rem;
  }
  
  .notification-actions {
    flex-direction: column;
  }
  
  .action-button {
    padding: 0.75rem;
  }
}
</style>

<template>
  <div class="settings-page">
    <div class="settings-header">
      <h1>アカウント管理</h1>
      <p class="subtitle">アカウント設定と管理</p>
    </div>

    <div class="settings-content">
      <Card>
        <template #header>
          <h2 class="section-title">アカウント情報</h2>
        </template>

        <div class="account-info">
          <div class="info-row">
            <div class="info-label">ユーザーID:</div>
            <div class="info-value">{{ user?.id || '-' }}</div>
          </div>
          <div class="info-row">
            <div class="info-label">アカウント作成日:</div>
            <div class="info-value">{{ formatDate(user?.created_at) }}</div>
          </div>
          <div class="info-row">
            <div class="info-label">最終ログイン:</div>
            <div class="info-value">{{ formatDate(user?.last_login) }}</div>
          </div>
          <div class="info-row">
            <div class="info-label">アカウント状態:</div>
            <div class="info-value">
              <span class="account-status" :class="{ 'active': user?.active }">
                {{ user?.active ? 'アクティブ' : '停止中' }}
              </span>
            </div>
          </div>
        </div>
      </Card>

      <Card class="mt-6">
        <template #header>
          <h2 class="section-title">セッション管理</h2>
        </template>

        <div class="sessions-container" v-if="!sessionsLoading">
          <div v-if="sessions.length === 0" class="no-sessions">
            アクティブなセッションはありません
          </div>
          <div v-else>
            <div class="session-list">
              <div v-for="(session, index) in sessions" :key="index" class="session-item">
                <div class="session-info">
                  <div class="session-device">
                    <Icon :name="getDeviceIcon(session.device_type)" class="device-icon" />
                    <div class="device-details">
                      <div class="device-name">{{ session.device_name }}</div>
                      <div class="device-meta">
                        {{ session.browser }} · {{ session.os }}
                      </div>
                    </div>
                  </div>
                  <div class="session-meta">
                    <div class="session-location">{{ session.location }}</div>
                    <div class="session-time">
                      <span v-if="session.current" class="current-session">現在のセッション</span>
                      <span v-else>最終アクセス: {{ formatTime(session.last_active) }}</span>
                    </div>
                  </div>
                </div>
                <div class="session-actions">
                  <Button 
                    v-if="!session.current"
                    variant="danger" 
                    size="sm" 
                    @click="revokeSession(session.id)"
                    :loading="revoking === session.id"
                  >
                    終了する
                  </Button>
                </div>
              </div>
            </div>

            <div class="sessions-actions">
              <Button 
                variant="danger" 
                @click="revokeAllSessions"
                :loading="revokingAll"
                :disabled="sessions.length <= 1"
              >
                他のすべてのセッションを終了
              </Button>
            </div>
          </div>
        </div>
        <div v-else class="sessions-loading">
          <div class="loading-spinner"></div>
          <p>セッション情報を読み込み中...</p>
        </div>
      </Card>

      <Card class="mt-6 danger-zone">
        <template #header>
          <h2 class="section-title">危険ゾーン</h2>
        </template>

        <div class="danger-actions">
          <div class="danger-action">
            <div class="danger-info">
              <h3>アカウントを削除</h3>
              <p>アカウントとすべての関連データは完全に削除されます。この操作は元に戻せません。</p>
            </div>
            <Button 
              variant="danger" 
              @click="showDeleteConfirm = true"
            >
              アカウントを削除
            </Button>
          </div>
        </div>
      </Card>
    </div>

    <!-- アカウント削除確認モーダル -->
    <div v-if="showDeleteConfirm" class="delete-modal">
      <div class="delete-modal-overlay" @click="showDeleteConfirm = false"></div>
      <div class="delete-modal-content">
        <h3 class="delete-modal-title">アカウントを削除しますか？</h3>
        <p class="delete-modal-message">
          この操作は元に戻せません。すべてのデータが完全に削除されます。
        </p>
        
        <div class="delete-confirmation">
          <Input
            v-model="deleteConfirmText"
            placeholder="削除するには「delete my account」と入力してください"
            :error="deleteError"
          />
        </div>
        
        <div class="delete-modal-actions">
          <Button 
            variant="text" 
            @click="showDeleteConfirm = false"
            :disabled="deleting"
          >
            キャンセル
          </Button>
          <Button 
            variant="danger" 
            @click="deleteAccount"
            :loading="deleting"
            :disabled="deleteConfirmText !== 'delete my account'"
          >
            アカウントを完全に削除
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '~/stores/auth';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();

// ユーザー情報
const user = computed(() => authStore.user);

// セッション管理
const sessions = ref([]);
const sessionsLoading = ref(false);
const revoking = ref(null);
const revokingAll = ref(false);

// アカウント削除
const showDeleteConfirm = ref(false);
const deleteConfirmText = ref('');
const deleteError = ref('');
const deleting = ref(false);

// セッション情報を取得
onMounted(async () => {
  if (authStore.isAuthenticated) {
    await fetchSessions();
  }
});

async function fetchSessions() {
  sessionsLoading.value = true;
  
  try {
    // 実際のAPI呼び出しはバックエンドの準備ができてから実装
    // 現在はモックデータで対応
    await new Promise(resolve => setTimeout(resolve, 1000)); // API呼び出しの遅延をシミュレート
    
    // モックセッションデータ
    sessions.value = [
      {
        id: 'session1',
        device_type: 'desktop',
        device_name: 'Windows PC',
        browser: 'Chrome 118',
        os: 'Windows 11',
        location: '東京, 日本',
        ip_address: '192.168.1.1',
        last_active: new Date(),
        current: true
      },
      {
        id: 'session2',
        device_type: 'mobile',
        device_name: 'iPhone 15',
        browser: 'Safari Mobile',
        os: 'iOS 17',
        location: '大阪, 日本',
        ip_address: '192.168.1.2',
        last_active: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2日前
        current: false
      },
      {
        id: 'session3',
        device_type: 'tablet',
        device_name: 'iPad Pro',
        browser: 'Firefox 117',
        os: 'iPadOS 16',
        location: '名古屋, 日本',
        ip_address: '192.168.1.3',
        last_active: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000), // 5日前
        current: false
      }
    ];
    
  } catch (error) {
    console.error('セッション情報の取得中にエラーが発生しました:', error);
  } finally {
    sessionsLoading.value = false;
  }
}

// セッションを取り消す
async function revokeSession(sessionId) {
  revoking.value = sessionId;
  
  try {
    // 実際のAPI呼び出しはバックエンドの準備ができてから実装
    // 現在はモックデータで対応
    await new Promise(resolve => setTimeout(resolve, 1000)); // API呼び出しの遅延をシミュレート
    
    // 成功後の処理
    sessions.value = sessions.value.filter(session => session.id !== sessionId);
    
  } catch (error) {
    console.error('セッションの取り消し中にエラーが発生しました:', error);
  } finally {
    revoking.value = null;
  }
}

// 他のすべてのセッションを取り消す
async function revokeAllSessions() {
  revokingAll.value = true;
  
  try {
    // 実際のAPI呼び出しはバックエンドの準備ができてから実装
    // 現在はモックデータで対応
    await new Promise(resolve => setTimeout(resolve, 1500)); // API呼び出しの遅延をシミュレート
    
    // 成功後の処理
    sessions.value = sessions.value.filter(session => session.current);
    
  } catch (error) {
    console.error('すべてのセッションの取り消し中にエラーが発生しました:', error);
  } finally {
    revokingAll.value = false;
  }
}

// アカウントを削除
async function deleteAccount() {
  if (deleteConfirmText.value !== 'delete my account') {
    deleteError.value = '確認テキストが一致しません';
    return;
  }
  
  deleting.value = true;
  deleteError.value = '';
  
  try {
    // 実際のAPI呼び出しはバックエンドの準備ができてから実装
    // 現在はモックデータで対応
    await new Promise(resolve => setTimeout(resolve, 2000)); // API呼び出しの遅延をシミュレート
    
    // 成功後の処理
    await authStore.logout();
    router.push('/');
    
  } catch (error) {
    console.error('アカウントの削除中にエラーが発生しました:', error);
    deleteError.value = 'アカウントの削除に失敗しました。後でもう一度お試しください。';
  } finally {
    deleting.value = false;
  }
}

// デバイスタイプに応じたアイコンを取得
function getDeviceIcon(deviceType) {
  switch (deviceType) {
    case 'desktop':
      return 'mdi:desktop-mac';
    case 'mobile':
      return 'mdi:cellphone';
    case 'tablet':
      return 'mdi:tablet';
    default:
      return 'mdi:devices';
  }
}

// 日付フォーマット
function formatDate(dateStr) {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return new Intl.DateTimeFormat('ja-JP', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  }).format(date);
}

// 時間フォーマット
function formatTime(dateStr) {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now - date;
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
  
  if (diffDays > 30) {
    return formatDate(dateStr);
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
.settings-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.settings-header {
  margin-bottom: 2rem;
}

.settings-header h1 {
  font-size: 1.75rem;
  font-weight: 600;
  margin: 0;
}

.subtitle {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
}

.mt-6 {
  margin-top: 1.5rem;
}

/* アカウント情報 */
.account-info {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.info-row {
  display: flex;
  align-items: center;
}

.info-label {
  width: 150px;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.info-value {
  flex: 1;
}

.account-status {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 2rem;
  font-size: 0.75rem;
  font-weight: 500;
}

.account-status.active {
  background-color: var(--color-success);
  color: white;
}

/* セッション管理 */
.sessions-container {
  margin-bottom: 1rem;
}

.no-sessions {
  text-align: center;
  padding: 2rem 0;
  color: var(--color-text-secondary);
}

.session-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.session-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  background-color: var(--color-bg-secondary);
}

.session-info {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.session-device {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.device-icon {
  font-size: 1.5rem;
  color: var(--color-text-secondary);
}

.device-name {
  font-weight: 500;
}

.device-meta {
  font-size: 0.8rem;
  color: var(--color-text-secondary);
}

.session-meta {
  margin-left: 2.25rem;
  font-size: 0.8rem;
  color: var(--color-text-secondary);
}

.current-session {
  color: var(--color-accent-primary);
  font-weight: 500;
}

.sessions-actions {
  margin-top: 1.5rem;
  display: flex;
  justify-content: flex-end;
}

.sessions-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem 0;
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

/* 危険ゾーン */
.danger-zone :deep(.gh-card__header) {
  background-color: rgba(215, 58, 73, 0.1);
  border-color: rgba(215, 58, 73, 0.2);
}

.danger-zone .section-title {
  color: var(--color-danger);
}

.danger-actions {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.danger-action {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
  border-bottom: 1px solid var(--color-border-primary);
}

.danger-action:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.danger-info h3 {
  font-size: 1rem;
  font-weight: 600;
  margin: 0 0 0.5rem 0;
}

.danger-info p {
  margin: 0;
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  max-width: 500px;
}

/* 削除確認モーダル */
.delete-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.delete-modal-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
}

.delete-modal-content {
  background-color: var(--color-bg-canvas);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-lg);
  width: 500px;
  max-width: 90%;
  padding: 1.5rem;
  position: relative;
  z-index: 1001;
}

.delete-modal-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0 0 1rem 0;
  color: var(--color-danger);
}

.delete-modal-message {
  margin: 0 0 1.5rem 0;
  color: var(--color-text-secondary);
}

.delete-confirmation {
  margin-bottom: 1.5rem;
}

.delete-modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .session-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .session-actions {
    align-self: flex-end;
  }
  
  .danger-action {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .danger-action button {
    align-self: flex-end;
  }
  
  .delete-modal-actions {
    flex-direction: column;
  }
  
  .delete-modal-actions button {
    width: 100%;
  }
}
</style>

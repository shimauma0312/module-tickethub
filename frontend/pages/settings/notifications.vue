<template>
  <div class="settings-page">
    <div class="settings-header">
      <h1>通知設定</h1>
      <p class="subtitle">通知の受け取り方法と頻度を設定</p>
    </div>

    <div class="settings-content">
      <Card>
        <template #header>
          <h2 class="section-title">通知設定</h2>
        </template>

        <form @submit.prevent="saveNotificationSettings" class="notification-form">
          <div class="notification-section">
            <h3 class="notification-section-title">通知チャネル</h3>
            <div class="notification-option">
              <div class="option-label">
                <div class="label-content">
                  <h4>アプリケーション内通知</h4>
                  <p class="option-description">Web画面内で通知を表示します</p>
                </div>
                <div class="toggle-switch">
                  <input
                    type="checkbox"
                    id="inAppNotifications"
                    v-model="settings.channels.inApp"
                    :disabled="loading"
                  />
                  <label for="inAppNotifications"></label>
                </div>
              </div>
            </div>

            <div class="notification-option">
              <div class="option-label">
                <div class="label-content">
                  <h4>メール通知</h4>
                  <p class="option-description">登録メールアドレスに通知を送信します</p>
                </div>
                <div class="toggle-switch">
                  <input
                    type="checkbox"
                    id="emailNotifications"
                    v-model="settings.channels.email"
                    :disabled="loading"
                  />
                  <label for="emailNotifications"></label>
                </div>
              </div>
            </div>
          </div>

          <div class="notification-section">
            <h3 class="notification-section-title">通知イベント</h3>
            
            <div class="notification-option">
              <div class="option-label">
                <div class="label-content">
                  <h4>新しいIssueがあなたに割り当てられたとき</h4>
                </div>
                <div class="toggle-switch">
                  <input
                    type="checkbox"
                    id="assignedIssue"
                    v-model="settings.events.assignedIssue"
                    :disabled="loading"
                  />
                  <label for="assignedIssue"></label>
                </div>
              </div>
            </div>

            <div class="notification-option">
              <div class="option-label">
                <div class="label-content">
                  <h4>担当Issueにコメントが追加されたとき</h4>
                </div>
                <div class="toggle-switch">
                  <input
                    type="checkbox"
                    id="issueComment"
                    v-model="settings.events.issueComment"
                    :disabled="loading"
                  />
                  <label for="issueComment"></label>
                </div>
              </div>
            </div>

            <div class="notification-option">
              <div class="option-label">
                <div class="label-content">
                  <h4>あなたがメンションされたとき</h4>
                </div>
                <div class="toggle-switch">
                  <input
                    type="checkbox"
                    id="mentioned"
                    v-model="settings.events.mentioned"
                    :disabled="loading"
                  />
                  <label for="mentioned"></label>
                </div>
              </div>
            </div>

            <div class="notification-option">
              <div class="option-label">
                <div class="label-content">
                  <h4>Issueのステータスが変更されたとき</h4>
                </div>
                <div class="toggle-switch">
                  <input
                    type="checkbox"
                    id="statusChange"
                    v-model="settings.events.statusChange"
                    :disabled="loading"
                  />
                  <label for="statusChange"></label>
                </div>
              </div>
            </div>

            <div class="notification-option">
              <div class="option-label">
                <div class="label-content">
                  <h4>システム通知と更新情報</h4>
                </div>
                <div class="toggle-switch">
                  <input
                    type="checkbox"
                    id="systemUpdates"
                    v-model="settings.events.systemUpdates"
                    :disabled="loading"
                  />
                  <label for="systemUpdates"></label>
                </div>
              </div>
            </div>
          </div>

          <div class="notification-section">
            <h3 class="notification-section-title">メール通知頻度</h3>
            <div class="email-frequency">
              <div class="radio-option">
                <input
                  type="radio"
                  id="immediate"
                  value="immediate"
                  v-model="settings.emailFrequency"
                  :disabled="loading || !settings.channels.email"
                />
                <label for="immediate">即時通知（イベント発生時）</label>
              </div>
              
              <div class="radio-option">
                <input
                  type="radio"
                  id="daily"
                  value="daily"
                  v-model="settings.emailFrequency"
                  :disabled="loading || !settings.channels.email"
                />
                <label for="daily">デイリーダイジェスト（1日1回）</label>
              </div>
              
              <div class="radio-option">
                <input
                  type="radio"
                  id="weekly"
                  value="weekly"
                  v-model="settings.emailFrequency"
                  :disabled="loading || !settings.channels.email"
                />
                <label for="weekly">ウィークリーダイジェスト（週1回）</label>
              </div>
            </div>
          </div>

          <div class="form-actions">
            <Button
              type="submit"
              variant="primary"
              :loading="loading"
              :disabled="!settingsChanged"
            >
              設定を保存
            </Button>
            <Button
              variant="text"
              @click="resetSettings"
              :disabled="loading || !settingsChanged"
            >
              変更を取り消し
            </Button>
          </div>
        </form>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '~/stores/auth';

const authStore = useAuthStore();
const loading = ref(false);
const successMessage = ref('');

// 設定の初期値
const defaultSettings = {
  channels: {
    inApp: true,
    email: false
  },
  events: {
    assignedIssue: true,
    issueComment: true,
    mentioned: true,
    statusChange: false,
    systemUpdates: true
  },
  emailFrequency: 'daily'
};

// 現在の設定とフォームの状態
const originalSettings = ref({ ...defaultSettings });
const settings = ref({ ...defaultSettings });

// 設定が変更されたかどうか
const settingsChanged = computed(() => {
  return JSON.stringify(settings.value) !== JSON.stringify(originalSettings.value);
});

// 設定のロード（実際のAPIが実装されるまではモック）
onMounted(async () => {
  if (authStore.isAuthenticated) {
    loading.value = true;
    
    try {
      // 実際のAPI呼び出しはバックエンドの準備ができてから実装
      // 現在はモックデータで対応
      await new Promise(resolve => setTimeout(resolve, 500)); // API呼び出しの遅延をシミュレート
      
      // ここでは例としてデフォルト設定を少し変更したものを「APIから取得した」と仮定
      const mockResponse = {
        channels: {
          inApp: true,
          email: true
        },
        events: {
          assignedIssue: true,
          issueComment: true,
          mentioned: true,
          statusChange: false,
          systemUpdates: false
        },
        emailFrequency: 'immediate'
      };
      
      originalSettings.value = { ...mockResponse };
      settings.value = { ...mockResponse };
      
    } catch (error) {
      console.error('通知設定の取得中にエラーが発生しました:', error);
    } finally {
      loading.value = false;
    }
  }
});

// 設定のリセット
function resetSettings() {
  settings.value = JSON.parse(JSON.stringify(originalSettings.value));
  successMessage.value = '';
}

// 設定の保存
async function saveNotificationSettings() {
  if (!settingsChanged.value) return;
  
  loading.value = true;
  successMessage.value = '';
  
  try {
    // 実際のAPI呼び出しはバックエンドの準備ができてから実装
    // 現在はモックデータで対応
    await new Promise(resolve => setTimeout(resolve, 1000)); // API呼び出しの遅延をシミュレート
    
    // 成功後の処理
    originalSettings.value = JSON.parse(JSON.stringify(settings.value));
    successMessage.value = '通知設定が更新されました';
    
  } catch (error) {
    console.error('通知設定の保存中にエラーが発生しました:', error);
  } finally {
    loading.value = false;
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

.notification-form {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.notification-section {
  border-bottom: 1px solid var(--color-border-primary);
  padding-bottom: 1.5rem;
}

.notification-section:last-of-type {
  border-bottom: none;
}

.notification-section-title {
  font-size: 1rem;
  font-weight: 600;
  margin: 0 0 1rem 0;
  color: var(--color-text-primary);
}

.notification-option {
  margin-bottom: 1rem;
}

.option-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.label-content h4 {
  font-size: 0.9rem;
  font-weight: 500;
  margin: 0;
  color: var(--color-text-primary);
}

.option-description {
  font-size: 0.8rem;
  color: var(--color-text-secondary);
  margin: 0.25rem 0 0 0;
}

/* トグルスイッチのスタイル */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 40px;
  height: 22px;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-switch label {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: .4s;
  border-radius: 34px;
}

.toggle-switch label:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 2px;
  bottom: 2px;
  background-color: white;
  transition: .4s;
  border-radius: 50%;
}

.toggle-switch input:checked + label {
  background-color: var(--color-accent-primary);
}

.toggle-switch input:focus + label {
  box-shadow: 0 0 1px var(--color-accent-primary);
}

.toggle-switch input:checked + label:before {
  transform: translateX(18px);
}

.toggle-switch input:disabled + label {
  opacity: 0.5;
  cursor: not-allowed;
}

/* ラジオボタンのスタイル */
.email-frequency {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.radio-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.radio-option input[type="radio"] {
  margin: 0;
}

.radio-option label {
  font-size: 0.9rem;
  color: var(--color-text-primary);
  cursor: pointer;
}

.radio-option input[type="radio"]:disabled + label {
  opacity: 0.5;
  cursor: not-allowed;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .form-actions {
    flex-direction: column;
  }
  
  .form-actions button {
    width: 100%;
  }
}
</style>

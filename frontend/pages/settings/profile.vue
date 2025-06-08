<template>
  <div class="settings-page">
    <div class="settings-header">
      <h1>プロフィール設定</h1>
      <p class="subtitle">アカウントプロフィール情報の編集</p>
    </div>

    <div class="settings-content">
      <Card>
        <template #header>
          <h2 class="section-title">個人情報</h2>
        </template>

        <form @submit.prevent="updateProfile" class="profile-form">
          <div class="form-group">
            <Input
              v-model="profileForm.username"
              label="ユーザーネーム"
              id="username"
              required
              :disabled="loading"
              :error="errors.username"
            />
          </div>

          <div class="form-group">
            <Input
              v-model="profileForm.displayName"
              label="表示名"
              id="displayName"
              :disabled="loading"
              :error="errors.displayName"
            />
          </div>

          <div class="form-group">
            <Input
              v-model="profileForm.email"
              label="メールアドレス"
              id="email"
              type="email"
              required
              :disabled="loading"
              :error="errors.email"
            />
          </div>

          <div class="form-group">
            <label class="textarea-label">自己紹介</label>
            <textarea
              v-model="profileForm.bio"
              class="bio-textarea"
              rows="4"
              :disabled="loading"
            ></textarea>
          </div>

          <div class="form-actions">
            <Button
              type="submit"
              variant="primary"
              :loading="loading"
              :disabled="!formChanged"
            >
              変更を保存
            </Button>
            <Button
              variant="text"
              @click="resetForm"
              :disabled="loading || !formChanged"
            >
              変更を取り消し
            </Button>
          </div>
        </form>
      </Card>

      <Card class="mt-6">
        <template #header>
          <h2 class="section-title">パスワード変更</h2>
        </template>

        <form @submit.prevent="updatePassword" class="password-form">
          <div class="form-group">
            <Input
              v-model="passwordForm.currentPassword"
              label="現在のパスワード"
              id="currentPassword"
              type="password"
              required
              :disabled="passwordLoading"
              :error="passwordErrors.currentPassword"
            />
          </div>

          <div class="form-group">
            <Input
              v-model="passwordForm.newPassword"
              label="新しいパスワード"
              id="newPassword"
              type="password"
              required
              :disabled="passwordLoading"
              :error="passwordErrors.newPassword"
              hint="8文字以上で、英字・数字を含める必要があります"
            />
          </div>

          <div class="form-group">
            <Input
              v-model="passwordForm.confirmPassword"
              label="新しいパスワード (確認)"
              id="confirmPassword"
              type="password"
              required
              :disabled="passwordLoading"
              :error="passwordErrors.confirmPassword"
            />
          </div>

          <div class="form-actions">
            <Button
              type="submit"
              variant="primary"
              :loading="passwordLoading"
              :disabled="!passwordFormValid"
            >
              パスワードを変更
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
const passwordLoading = ref(false);
const successMessage = ref('');
const errors = ref({});
const passwordErrors = ref({});

// プロフィールフォームの初期化
const originalProfile = ref({});
const profileForm = ref({
  username: '',
  displayName: '',
  email: '',
  bio: ''
});

// パスワード変更フォーム
const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
});

// フォームのバリデーション
const formChanged = computed(() => {
  return profileForm.value.username !== originalProfile.value.username ||
    profileForm.value.displayName !== originalProfile.value.displayName ||
    profileForm.value.email !== originalProfile.value.email ||
    profileForm.value.bio !== originalProfile.value.bio;
});

const passwordFormValid = computed(() => {
  return passwordForm.value.currentPassword &&
    passwordForm.value.newPassword &&
    passwordForm.value.confirmPassword &&
    passwordForm.value.newPassword === passwordForm.value.confirmPassword &&
    passwordForm.value.newPassword.length >= 8;
});

// プロフィール情報の読み込み
onMounted(() => {
  if (authStore.user) {
    originalProfile.value = {
      username: authStore.user.username || '',
      displayName: authStore.user.display_name || '',
      email: authStore.user.email || '',
      bio: authStore.user.bio || ''
    };
    
    // フォームを初期値でセット
    profileForm.value = { ...originalProfile.value };
  }
});

// フォームのリセット
function resetForm() {
  profileForm.value = { ...originalProfile.value };
  errors.value = {};
  successMessage.value = '';
}

// プロフィール更新
async function updateProfile() {
  loading.value = true;
  errors.value = {};
  successMessage.value = '';
  
  try {
    // 実際のAPI呼び出しはバックエンドの準備ができてから実装
    // 現在はモックデータで対応
    await new Promise(resolve => setTimeout(resolve, 1000)); // API呼び出しの遅延をシミュレート
    
    // 成功後の処理
    originalProfile.value = { ...profileForm.value };
    successMessage.value = 'プロフィールが更新されました';
    
    // ユーザー情報の更新（仮実装）
    const updatedUser = {
      ...authStore.user,
      username: profileForm.value.username,
      display_name: profileForm.value.displayName,
      email: profileForm.value.email,
      bio: profileForm.value.bio
    };
    
    // ストアとローカルストレージを更新
    authStore.user = updatedUser;
    if (process.client) {
      localStorage.setItem('user', JSON.stringify(updatedUser));
    }
    
  } catch (error) {
    // エラーハンドリング
    console.error('プロフィール更新エラー:', error);
    errors.value = {
      form: 'プロフィールの更新中にエラーが発生しました。'
    };
  } finally {
    loading.value = false;
  }
}

// パスワード更新
async function updatePassword() {
  if (!passwordFormValid.value) return;
  
  passwordLoading.value = true;
  passwordErrors.value = {};
  
  try {
    // 実際のAPI呼び出しはバックエンドの準備ができてから実装
    // 現在はモックデータで対応
    await new Promise(resolve => setTimeout(resolve, 1000)); // API呼び出しの遅延をシミュレート
    
    // 成功後の処理
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    };
    
    successMessage.value = 'パスワードが更新されました';
    
  } catch (error) {
    // エラーハンドリング
    console.error('パスワード更新エラー:', error);
    passwordErrors.value = {
      form: 'パスワードの更新中にエラーが発生しました。'
    };
  } finally {
    passwordLoading.value = false;
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

.profile-form, .password-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.textarea-label {
  font-weight: 500;
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
}

.bio-textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  background-color: var(--color-bg-canvas);
  font-size: var(--font-size-sm);
  transition: var(--transition-base);
}

.bio-textarea:focus {
  outline: none;
  border-color: var(--color-accent-primary);
  box-shadow: 0 0 0 3px rgba(3, 102, 214, 0.3);
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.mt-6 {
  margin-top: 1.5rem;
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

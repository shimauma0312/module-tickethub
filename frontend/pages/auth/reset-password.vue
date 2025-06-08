<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <h1 class="auth-title">パスワードリセット</h1>
        <p class="auth-subtitle">
          {{ step === 'request' ? 'アカウントに登録したメールアドレスを入力してください' : 
             step === 'verify' ? '確認コードを入力してください' : 
             '新しいパスワードを設定してください' }}
        </p>
      </div>
      
      <div v-if="error" class="auth-error">
        {{ error }}
      </div>
      
      <div v-if="success" class="auth-success">
        {{ success }}
      </div>
      
      <!-- ステップ1: リセットリクエスト -->
      <form v-if="step === 'request'" @submit.prevent="handleRequestReset" class="auth-form">
        <div class="form-group">
          <Input
            v-model="resetData.email"
            label="メールアドレス"
            type="email"
            placeholder="example@example.com"
            :error="!!errors.email"
            required
          />
          <p v-if="errors.email" class="form-error">{{ errors.email }}</p>
        </div>
        
        <div class="form-actions">
          <Button 
            type="submit" 
            variant="primary" 
            full 
            :loading="loading"
          >
            リセットコードを送信
          </Button>
        </div>
      </form>
      
      <!-- ステップ2: コード確認 -->
      <form v-if="step === 'verify'" @submit.prevent="handleVerifyCode" class="auth-form">
        <div class="form-group">
          <Input
            v-model="resetData.token"
            label="確認コード"
            placeholder="メールに送信されたコードを入力"
            :error="!!errors.token"
            required
          />
          <p v-if="errors.token" class="form-error">{{ errors.token }}</p>
        </div>
        
        <div class="form-actions">
          <Button 
            type="submit" 
            variant="primary" 
            full 
            :loading="loading"
          >
            コードを確認
          </Button>
        </div>
      </form>
      
      <!-- ステップ3: 新パスワード設定 -->
      <form v-if="step === 'reset'" @submit.prevent="handleCompleteReset" class="auth-form">
        <div class="form-group">
          <Input
            v-model="resetData.newPassword"
            label="新しいパスワード"
            type="password"
            placeholder="8文字以上で入力"
            :error="!!errors.newPassword"
            required
          />
          <p v-if="errors.newPassword" class="form-error">{{ errors.newPassword }}</p>
        </div>
        
        <div class="form-group">
          <Input
            v-model="resetData.confirmPassword"
            label="新しいパスワード（確認）"
            type="password"
            placeholder="パスワードを再入力"
            :error="!!errors.confirmPassword"
            required
          />
          <p v-if="errors.confirmPassword" class="form-error">{{ errors.confirmPassword }}</p>
        </div>
        
        <div class="form-actions">
          <Button 
            type="submit" 
            variant="primary" 
            full 
            :loading="loading"
          >
            パスワードをリセット
          </Button>
        </div>
      </form>
      
      <div class="auth-footer">
        <p>
          <NuxtLink to="/auth/login" class="auth-link">
            ログイン画面に戻る
          </NuxtLink>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const loading = ref(false);
const error = ref('');
const success = ref('');
const step = ref('request'); // 'request', 'verify', 'reset'

const resetData = reactive({
  email: '',
  token: '',
  newPassword: '',
  confirmPassword: ''
});

const errors = reactive({
  email: '',
  token: '',
  newPassword: '',
  confirmPassword: ''
});

// ステップ1: リセットリクエスト
async function handleRequestReset() {
  // メールアドレスのバリデーション
  errors.email = '';
  if (!resetData.email) {
    errors.email = 'メールアドレスを入力してください';
    return;
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(resetData.email)) {
    errors.email = '正しいメールアドレス形式で入力してください';
    return;
  }
  
  loading.value = true;
  error.value = '';
  success.value = '';
  
  try {
    // APIエンドポイントを使用してリセットリクエスト
    const config = useRuntimeConfig();
    const response = await fetch(`${config.public.apiBaseUrl}/auth/password-reset`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: resetData.email
      }),
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || 'パスワードリセットリクエストに失敗しました');
    }
    
    // 成功メッセージ
    success.value = 'メールアドレスに確認コードを送信しました。メールをご確認ください。';
    // 次のステップへ
    step.value = 'verify';
  } catch (err) {
    console.error('Password reset request error:', err);
    error.value = err.message || 'リクエスト中にエラーが発生しました。もう一度お試しください。';
  } finally {
    loading.value = false;
  }
}

// ステップ2: コード確認
async function handleVerifyCode() {
  // トークンのバリデーション
  errors.token = '';
  if (!resetData.token) {
    errors.token = '確認コードを入力してください';
    return;
  }
  
  loading.value = true;
  error.value = '';
  success.value = '';
  
  try {
    // APIエンドポイントを使用してコード確認
    const config = useRuntimeConfig();
    const response = await fetch(`${config.public.apiBaseUrl}/auth/password-reset/validate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: resetData.email,
        token: resetData.token
      }),
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || '確認コードが無効です');
    }
    
    // 成功メッセージ
    success.value = '確認コードが検証されました。新しいパスワードを設定してください。';
    // 次のステップへ
    step.value = 'reset';
  } catch (err) {
    console.error('Code verification error:', err);
    error.value = err.message || 'コード確認中にエラーが発生しました。もう一度お試しください。';
  } finally {
    loading.value = false;
  }
}

// ステップ3: パスワードリセット完了
async function handleCompleteReset() {
  // パスワードのバリデーション
  errors.newPassword = '';
  errors.confirmPassword = '';
  
  if (!resetData.newPassword) {
    errors.newPassword = '新しいパスワードを入力してください';
    return;
  } else if (resetData.newPassword.length < 8) {
    errors.newPassword = 'パスワードは8文字以上で入力してください';
    return;
  }
  
  if (!resetData.confirmPassword) {
    errors.confirmPassword = '確認用パスワードを入力してください';
    return;
  } else if (resetData.newPassword !== resetData.confirmPassword) {
    errors.confirmPassword = 'パスワードが一致しません';
    return;
  }
  
  loading.value = true;
  error.value = '';
  success.value = '';
  
  try {
    // APIエンドポイントを使用してパスワードリセット完了
    const config = useRuntimeConfig();
    const response = await fetch(`${config.public.apiBaseUrl}/auth/password-reset/complete`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: resetData.email,
        token: resetData.token,
        new_password: resetData.newPassword
      }),
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || 'パスワードリセットに失敗しました');
    }
    
    // 成功メッセージ
    success.value = 'パスワードが正常にリセットされました。新しいパスワードでログインしてください。';
    
    // 3秒後にログイン画面へリダイレクト
    setTimeout(() => {
      router.push('/auth/login');
    }, 3000);
  } catch (err) {
    console.error('Password reset error:', err);
    error.value = err.message || 'パスワードリセット中にエラーが発生しました。もう一度お試しください。';
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 200px);
  padding: 2rem 1rem;
}

.auth-card {
  background-color: white;
  border-radius: 8px;
  box-shadow: var(--shadow-md);
  width: 100%;
  max-width: 450px;
  padding: 2rem;
}

.auth-header {
  text-align: center;
  margin-bottom: 2rem;
}

.auth-title {
  font-size: 1.75rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: var(--color-text-primary);
}

.auth-subtitle {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
}

.auth-error {
  background-color: #ffebe9;
  border: 1px solid #ff8182;
  border-radius: 4px;
  color: #cf222e;
  padding: 0.75rem;
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
}

.auth-success {
  background-color: #dafbe1;
  border: 1px solid #56d364;
  border-radius: 4px;
  color: #116329;
  padding: 0.75rem;
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
}

.auth-form {
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-error {
  margin-top: 0.25rem;
  color: var(--color-danger);
  font-size: 0.75rem;
}

.form-actions {
  margin-top: 1.5rem;
}

.auth-footer {
  text-align: center;
  margin-top: 1.5rem;
  color: var(--color-text-secondary);
  font-size: 0.875rem;
}

.auth-link {
  color: var(--color-accent-primary);
  text-decoration: none;
}

.auth-link:hover {
  text-decoration: underline;
}
</style>

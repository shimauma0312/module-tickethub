<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <h1 class="auth-title">ログイン</h1>
        <p class="auth-subtitle">アカウント情報を入力してログインしてください</p>
      </div>
      
      <div v-if="error" class="auth-error">
        {{ error }}
      </div>
      
      <form @submit.prevent="handleLogin" class="auth-form">
        <div class="form-group">
          <Input
            v-model="credentials.usernameOrEmail"
            label="ユーザー名またはメールアドレス"
            placeholder="username or example@example.com"
            :error="!!errors.usernameOrEmail"
            required
          />
          <p v-if="errors.usernameOrEmail" class="form-error">{{ errors.usernameOrEmail }}</p>
        </div>
        
        <div class="form-group">
          <Input
            v-model="credentials.password"
            label="パスワード"
            type="password"
            placeholder="パスワードを入力"
            :error="!!errors.password"
            required
          />
          <p v-if="errors.password" class="form-error">{{ errors.password }}</p>
          <div class="forgot-password">
            <NuxtLink to="/auth/reset-password" class="forgot-password-link">
              パスワードをお忘れですか？
            </NuxtLink>
          </div>
        </div>
        
        <div class="form-actions">
          <Button 
            type="submit" 
            variant="primary" 
            full 
            :loading="loading"
          >
            ログイン
          </Button>
        </div>
      </form>
      
      <div class="auth-footer">
        <p>
          アカウントをお持ちでない方は
          <NuxtLink to="/auth/register" class="auth-link">
            新規登録
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

const credentials = reactive({
  usernameOrEmail: '',
  password: ''
});

const errors = reactive({
  usernameOrEmail: '',
  password: ''
});

// バリデーション
function validate() {
  let isValid = true;
  errors.usernameOrEmail = '';
  errors.password = '';
  
  if (!credentials.usernameOrEmail) {
    errors.usernameOrEmail = 'ユーザー名またはメールアドレスを入力してください';
    isValid = false;
  }
  
  if (!credentials.password) {
    errors.password = 'パスワードを入力してください';
    isValid = false;
  }
  
  return isValid;
}

// ログイン処理
async function handleLogin() {
  if (!validate()) return;
  
  loading.value = true;
  error.value = '';
  
  try {
    // APIエンドポイントを使用してログイン
    const config = useRuntimeConfig();
    const response = await fetch(`${config.public.apiBaseUrl}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username_or_email: credentials.usernameOrEmail,
        password: credentials.password,
      }),
      credentials: 'include', // クッキーを含める
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || 'ログインに失敗しました');
    }
    
    // トークンの保存（Pinia ストアなどで管理予定）
    // ローカルストレージは一時的な措置
    localStorage.setItem('access_token', data.token.access_token);
    localStorage.setItem('refresh_token', data.token.refresh_token);
    
    // ユーザー情報の保存（仮実装）
    localStorage.setItem('user', JSON.stringify(data.user));
    
    // ダッシュボードにリダイレクト
    router.push('/');
  } catch (err) {
    console.error('Login error:', err);
    error.value = err.message || 'ログイン中にエラーが発生しました。もう一度お試しください。';
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

.forgot-password {
  margin-top: 0.5rem;
  text-align: right;
}

.forgot-password-link {
  color: var(--color-accent-primary);
  font-size: 0.875rem;
  text-decoration: none;
}

.forgot-password-link:hover {
  text-decoration: underline;
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

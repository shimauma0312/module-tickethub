<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <h1 class="auth-title">アカウント登録</h1>
        <p class="auth-subtitle">新しいアカウントを作成してTicketHubを始めましょう</p>
      </div>
      
      <div v-if="error" class="auth-error">
        {{ error }}
      </div>
      
      <form @submit.prevent="handleRegister" class="auth-form">
        <div class="form-group">
          <Input
            v-model="user.username"
            label="ユーザー名"
            placeholder="username"
            :error="!!errors.username"
            required
          />
          <p v-if="errors.username" class="form-error">{{ errors.username }}</p>
        </div>
        
        <div class="form-group">
          <Input
            v-model="user.email"
            label="メールアドレス"
            type="email"
            placeholder="example@example.com"
            :error="!!errors.email"
            required
          />
          <p v-if="errors.email" class="form-error">{{ errors.email }}</p>
        </div>
        
        <div class="form-group">
          <Input
            v-model="user.fullName"
            label="氏名"
            placeholder="山田 太郎"
            :error="!!errors.fullName"
          />
          <p v-if="errors.fullName" class="form-error">{{ errors.fullName }}</p>
        </div>
        
        <div class="form-group">
          <Input
            v-model="user.password"
            label="パスワード"
            type="password"
            placeholder="8文字以上で入力"
            :error="!!errors.password"
            required
          />
          <p v-if="errors.password" class="form-error">{{ errors.password }}</p>
        </div>
        
        <div class="form-group">
          <Input
            v-model="user.confirmPassword"
            label="パスワード（確認）"
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
            アカウント作成
          </Button>
        </div>
      </form>
      
      <div class="auth-footer">
        <p>
          すでにアカウントをお持ちの方は
          <NuxtLink to="/auth/login" class="auth-link">
            ログイン
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

const user = reactive({
  username: '',
  email: '',
  fullName: '',
  password: '',
  confirmPassword: ''
});

const errors = reactive({
  username: '',
  email: '',
  fullName: '',
  password: '',
  confirmPassword: ''
});

// バリデーション
function validate() {
  let isValid = true;
  errors.username = '';
  errors.email = '';
  errors.fullName = '';
  errors.password = '';
  errors.confirmPassword = '';
  
  // ユーザー名チェック
  if (!user.username) {
    errors.username = 'ユーザー名を入力してください';
    isValid = false;
  } else if (user.username.length < 3) {
    errors.username = 'ユーザー名は3文字以上で入力してください';
    isValid = false;
  }
  
  // メールアドレスチェック
  if (!user.email) {
    errors.email = 'メールアドレスを入力してください';
    isValid = false;
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(user.email)) {
    errors.email = '正しいメールアドレス形式で入力してください';
    isValid = false;
  }
  
  // 氏名チェック（任意）
  if (user.fullName && user.fullName.length > 100) {
    errors.fullName = '氏名は100文字以内で入力してください';
    isValid = false;
  }
  
  // パスワードチェック
  if (!user.password) {
    errors.password = 'パスワードを入力してください';
    isValid = false;
  } else if (user.password.length < 8) {
    errors.password = 'パスワードは8文字以上で入力してください';
    isValid = false;
  }
  
  // パスワード（確認）チェック
  if (!user.confirmPassword) {
    errors.confirmPassword = '確認用パスワードを入力してください';
    isValid = false;
  } else if (user.password !== user.confirmPassword) {
    errors.confirmPassword = 'パスワードが一致しません';
    isValid = false;
  }
  
  return isValid;
}

// 登録処理
async function handleRegister() {
  if (!validate()) return;
  
  loading.value = true;
  error.value = '';
  
  try {
    // APIエンドポイントを使用してユーザー登録
    const config = useRuntimeConfig();
    const response = await fetch(`${config.public.apiBaseUrl}/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: user.username,
        email: user.email,
        full_name: user.fullName,
        password: user.password,
      }),
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || 'ユーザー登録に失敗しました');
    }
    
    // 登録成功メッセージ
    alert('ユーザー登録が完了しました。ログインしてください。');
    
    // ログインページへリダイレクト
    router.push('/auth/login');
  } catch (err) {
    console.error('Registration error:', err);
    
    // エラーメッセージの処理
    if (err.message.includes('username already exists')) {
      error.value = 'このユーザー名は既に使用されています';
    } else if (err.message.includes('email already exists')) {
      error.value = 'このメールアドレスは既に登録されています';
    } else {
      error.value = err.message || 'ユーザー登録中にエラーが発生しました。もう一度お試しください。';
    }
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
  max-width: 500px;
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

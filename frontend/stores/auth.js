import { defineStore } from 'pinia';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    accessToken: null,
    refreshToken: null,
    loading: false,
    error: null,
  }),
  
  getters: {
    isAuthenticated: (state) => !!state.accessToken && !!state.user,
    isAdmin: (state) => state.user && state.user.is_admin,
    userProfile: (state) => state.user,
  },
  
  actions: {
    // ローカルストレージからユーザー情報をロード
    initFromLocalStorage() {
      if (process.client) {
        const accessToken = localStorage.getItem('access_token');
        const refreshToken = localStorage.getItem('refresh_token');
        const user = localStorage.getItem('user');
        
        if (accessToken) this.accessToken = accessToken;
        if (refreshToken) this.refreshToken = refreshToken;
        if (user) {
          try {
            this.user = JSON.parse(user);
          } catch (e) {
            console.error('Failed to parse user from localStorage', e);
            this.user = null;
          }
        }
      }
    },
    
    // ログイン
    async login(usernameOrEmail, password) {
      this.loading = true;
      this.error = null;
      
      try {
        const config = useRuntimeConfig();
        const response = await fetch(`${config.public.apiBaseUrl}/auth/login`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            username_or_email: usernameOrEmail,
            password: password,
          }),
          credentials: 'include', // クッキーを含める
        });
        
        const data = await response.json();
        
        if (!response.ok) {
          throw new Error(data.error || 'ログインに失敗しました');
        }
        
        // ステート更新
        this.accessToken = data.token.access_token;
        this.refreshToken = data.token.refresh_token;
        this.user = data.user;
        
        // ローカルストレージに保存
        if (process.client) {
          localStorage.setItem('access_token', data.token.access_token);
          localStorage.setItem('refresh_token', data.token.refresh_token);
          localStorage.setItem('user', JSON.stringify(data.user));
        }
        
        return true;
      } catch (err) {
        this.error = err.message;
        return false;
      } finally {
        this.loading = false;
      }
    },
    
    // ログアウト
    async logout() {
      this.loading = true;
      this.error = null;
      
      try {
        // すでにトークンがない場合は何もしない
        if (!this.accessToken) {
          this.clearAuth();
          return true;
        }
        
        const config = useRuntimeConfig();
        const response = await fetch(`${config.public.apiBaseUrl}/auth/logout`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${this.accessToken}`,
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        });
        
        // ステータスコードに関わらず認証情報をクリア
        this.clearAuth();
        
        return response.ok;
      } catch (err) {
        this.error = err.message;
        // エラーが発生してもログアウトは行う
        this.clearAuth();
        return false;
      } finally {
        this.loading = false;
      }
    },
    
    // 全デバイスでログアウト
    async logoutAll() {
      this.loading = true;
      this.error = null;
      
      try {
        if (!this.accessToken) {
          this.clearAuth();
          return true;
        }
        
        const config = useRuntimeConfig();
        const response = await fetch(`${config.public.apiBaseUrl}/auth/logout-all`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${this.accessToken}`,
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        });
        
        // ステータスコードに関わらず認証情報をクリア
        this.clearAuth();
        
        return response.ok;
      } catch (err) {
        this.error = err.message;
        // エラーが発生してもログアウトは行う
        this.clearAuth();
        return false;
      } finally {
        this.loading = false;
      }
    },
    
    // トークンのリフレッシュ
    async refreshAccessToken() {
      if (!this.refreshToken) return false;
      
      try {
        const config = useRuntimeConfig();
        const response = await fetch(`${config.public.apiBaseUrl}/auth/refresh-token`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            refresh_token: this.refreshToken,
          }),
          credentials: 'include',
        });
        
        const data = await response.json();
        
        if (!response.ok) {
          throw new Error(data.error || 'トークンの更新に失敗しました');
        }
        
        // 新しいトークンで更新
        this.accessToken = data.access_token;
        this.refreshToken = data.refresh_token || this.refreshToken;
        
        // ローカルストレージに保存
        if (process.client) {
          localStorage.setItem('access_token', this.accessToken);
          if (data.refresh_token) {
            localStorage.setItem('refresh_token', data.refresh_token);
          }
        }
        
        return true;
      } catch (err) {
        console.error('Token refresh failed:', err);
        // リフレッシュに失敗した場合は認証情報をクリア
        this.clearAuth();
        return false;
      }
    },
    
    // ユーザープロフィール取得
    async fetchUserProfile() {
      if (!this.accessToken) return false;
      
      this.loading = true;
      this.error = null;
      
      try {
        const config = useRuntimeConfig();
        const response = await fetch(`${config.public.apiBaseUrl}/api/v1/users/me`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${this.accessToken}`,
          },
        });
        
        // トークンが無効な場合はリフレッシュを試みる
        if (response.status === 401) {
          const refreshed = await this.refreshAccessToken();
          if (refreshed) {
            return await this.fetchUserProfile(); // 再帰的に再試行
          } else {
            throw new Error('認証の有効期限が切れました。再度ログインしてください。');
          }
        }
        
        const data = await response.json();
        
        if (!response.ok) {
          throw new Error(data.error || 'プロフィール取得に失敗しました');
        }
        
        // ユーザー情報を更新
        this.user = data;
        
        // ローカルストレージに保存
        if (process.client) {
          localStorage.setItem('user', JSON.stringify(data));
        }
        
        return true;
      } catch (err) {
        this.error = err.message;
        return false;
      } finally {
        this.loading = false;
      }
    },
    
    // 認証情報をクリア
    clearAuth() {
      this.user = null;
      this.accessToken = null;
      this.refreshToken = null;
      
      // ローカルストレージからも削除
      if (process.client) {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
      }
    },
  },
});

<template>
  <div class="app-container">
    <header class="header">
      <div class="header-container">
        <div class="header-content">
          <div class="logo-area">
            <h1 class="logo">TicketHub</h1>
          </div>
          <nav class="main-nav">
            <NuxtLink to="/" class="nav-link">
              Dashboard
            </NuxtLink>
            <NuxtLink to="/issues" class="nav-link">
              Issues
            </NuxtLink>
            <NuxtLink to="/discussions" class="nav-link">
              Discussions
            </NuxtLink>
          </nav>
          
          <div class="user-nav">
            <!-- 検索バー -->
            <div class="search-trigger" @click="showSearchModal = true">
              <Icon name="mdi:magnify" class="icon" />
              <span class="search-text">検索...</span>
              <span class="search-shortcut">/</span>
            </div>
            
            <!-- 通知アイコン -->
            <div v-if="isAuthenticated" class="notification-icon" @click="showNotifications = !showNotifications">
              <div class="icon-wrapper">
                <Icon name="mdi:bell-outline" class="icon" />
                <span v-if="unreadNotifications > 0" class="notification-badge">{{ unreadNotifications }}</span>
              </div>
              
              <!-- 通知ドロップダウン -->
              <div v-if="showNotifications" class="notification-dropdown">
                <div class="notification-header">
                  <h3>通知</h3>
                  <button @click="markAllAsRead" class="read-all-btn">すべて既読</button>
                </div>
                
                <div v-if="loading" class="notification-loading">
                  <div class="loading-spinner"></div>
                  <p>読み込み中...</p>
                </div>
                
                <div v-else-if="notifications.length === 0" class="no-notifications">
                  通知はありません
                </div>
                
                <ul v-else class="notification-list">
                  <li v-for="notification in notifications" :key="notification.id" 
                      class="notification-item" 
                      :class="{ 'unread': !notification.read_at }">
                    <div class="notification-content">
                      <Icon :name="getNotificationIcon(notification.type)" class="notification-type-icon" />
                      <div>
                        <div class="notification-text" v-html="notification.message"></div>
                        <div class="notification-time">{{ formatTime(notification.created_at) }}</div>
                      </div>
                    </div>
                    <button @click="markAsRead(notification.id)" class="mark-read-btn" v-if="!notification.read_at">
                      <Icon name="mdi:check" />
                    </button>
                  </li>
                </ul>
                
                <div class="notification-footer">
                  <NuxtLink to="/notifications" class="view-all-link">すべての通知を見る</NuxtLink>
                </div>
              </div>
            </div>
            
            <!-- ユーザーメニュー -->
            <div v-if="isAuthenticated" class="user-menu" @click="showUserDropdown = !showUserDropdown">
              <div class="user-avatar">
                {{ userInitials }}
              </div>
              
              <!-- ユーザードロップダウン -->
              <div v-if="showUserDropdown" class="user-dropdown">
                <div class="user-info">
                  <div class="user-name">{{ user.username }}</div>
                  <div class="user-email">{{ user.email }}</div>
                </div>
                
                <div class="dropdown-divider"></div>
                
                <ul class="dropdown-menu">
                  <li>
                    <NuxtLink to="/settings/profile" class="dropdown-item">
                      <Icon name="mdi:account-outline" />
                      プロフィール設定
                    </NuxtLink>
                  </li>
                  <li>
                    <NuxtLink to="/settings/notifications" class="dropdown-item">
                      <Icon name="mdi:bell-outline" />
                      通知設定
                    </NuxtLink>
                  </li>
                  <li v-if="isAdmin">
                    <NuxtLink to="/admin" class="dropdown-item">
                      <Icon name="mdi:shield-account-outline" />
                      管理者ダッシュボード
                    </NuxtLink>
                  </li>
                </ul>
                
                <div class="dropdown-divider"></div>
                
                <button @click="logout" class="logout-button">
                  <Icon name="mdi:logout" />
                  ログアウト
                </button>
              </div>
            </div>
            
            <!-- 未ログイン時のリンク -->
            <div v-else class="auth-links">
              <NuxtLink to="/auth/login" class="auth-link">ログイン</NuxtLink>
              <NuxtLink to="/auth/register" class="auth-link register-link">アカウント登録</NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </header>

    <main class="main-content">
      <div class="page-container">
        <NuxtPage />
      </div>
    </main>

    <!-- 検索モーダル -->
    <div v-if="showSearchModal" class="search-modal">
      <div class="search-modal-overlay" @click="showSearchModal = false"></div>
      <div class="search-modal-content">
        <div class="search-input-container">
          <Icon name="mdi:magnify" class="search-icon" />
          <input
            v-model="searchQuery"
            type="text"
            class="search-input"
            placeholder="検索キーワードを入力..."
            @keyup.enter="performSearch"
            @keyup.esc="showSearchModal = false"
            ref="searchInput"
          />
          <button v-if="searchQuery" @click="searchQuery = ''" class="clear-search-btn">
            <Icon name="mdi:close" />
          </button>
        </div>
        
        <div class="search-filters">
          <span class="filter-label">フィルター:</span>
          <div class="filter-tags">
            <button 
              v-for="filter in searchFilters" 
              :key="filter.id"
              class="filter-tag"
              :class="{ 'active': filter.active }"
              @click="toggleFilter(filter.id)"
            >
              {{ filter.label }}
            </button>
          </div>
        </div>
        
        <div class="keyboard-shortcuts">
          <span><kbd>↑</kbd> <kbd>↓</kbd> 移動</span>
          <span><kbd>Enter</kbd> 選択</span>
          <span><kbd>Esc</kbd> 閉じる</span>
        </div>
      </div>
    </div>

    <footer class="footer">
      <div class="footer-container">
        <p class="copyright">
          &copy; 2025 TicketHub -
        </p>
      </div>
    </footer>
  </div>
</template>

<style>
:root {
  /* GitHub風のカラーパレット */
  --color-bg-canvas: #ffffff;
  --color-bg-secondary: #f6f8fa;
  --color-border-primary: #e1e4e8;
  --color-text-primary: #24292e;
  --color-text-secondary: #586069;
  --color-accent-primary: #0366d6;
  --color-accent-secondary: #2ea44f;
  --color-danger: #d73a49;
  --color-warning: #ffea7f;
  --color-success: #2ea44f;
  
  /* スペーシング変数 */
  --spacing-xxs: 4px;
  --spacing-xs: 8px;
  --spacing-sm: 12px;
  --spacing-md: 16px;
  --spacing-lg: 24px;
  --spacing-xl: 32px;
  --spacing-xxl: 48px;
  
  /* シャドウ変数 */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 3px 6px rgba(0, 0, 0, 0.1);
  --shadow-lg: 0 8px 16px rgba(0, 0, 0, 0.1);
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
  font-size: 14px;
  line-height: 1.5;
  color: var(--color-text-primary);
  background-color: var(--color-bg-secondary);
}

.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background-color: var(--color-bg-canvas);
  border-bottom: 1px solid var(--color-border-primary);
  box-shadow: var(--shadow-sm);
}

.header-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 var(--spacing-md);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 64px;
}

.logo-area {
  display: flex;
  align-items: center;
}

.logo {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.main-nav {
  display: flex;
  gap: var(--spacing-md);
}

.nav-link {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: 6px;
  color: var(--color-text-secondary);
  text-decoration: none;
  font-weight: 500;
  font-size: 14px;
  transition: background-color 0.2s;
}

.nav-link:hover {
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
}

.main-content {
  flex: 1;
}

.page-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: var(--spacing-xl) var(--spacing-md);
}

.footer {
  background-color: var(--color-bg-canvas);
  border-top: 1px solid var(--color-border-primary);
  padding: var(--spacing-lg) 0;
}

.footer-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 var(--spacing-md);
}

.copyright {
  text-align: center;
  font-size: 12px;
  color: var(--color-text-secondary);
}

/* レスポンシブデザイン対応 */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    height: auto;
    padding: var(--spacing-md) 0;
    gap: var(--spacing-md);
  }
  
  .main-nav {
    width: 100%;
    justify-content: center;
  }
}
</style>

import { useAuthStore } from '~/stores/auth';

export default defineNuxtRouteMiddleware(async (to, from) => {
  // SPA モードで動作中かつクライアントサイドでのみ実行
  if (process.client) {
    const authStore = useAuthStore();
    
    // ストアの初期化（ローカルストレージからのロード）
    authStore.initFromLocalStorage();
    
    // 認証が必要なルートかつ、未認証の場合
    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
      // 元々のURLをクエリパラメータに含めてログインページにリダイレクト
      return navigateTo({
        path: '/auth/login',
        query: { redirect: to.fullPath }
      });
    }
    
    // 管理者権限が必要なルートかつ、管理者でない場合
    if (to.meta.requiresAdmin && !authStore.isAdmin) {
      // 権限がなければダッシュボードにリダイレクト
      return navigateTo('/');
    }
    
    // すでに認証済みの状態で認証ページにアクセスした場合、ダッシュボードへリダイレクト
    if (authStore.isAuthenticated && to.meta.isAuthRoute) {
      return navigateTo('/');
    }
  }
});

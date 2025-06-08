import { defineStore } from 'pinia';

export const useNotificationsStore = defineStore('notifications', {
  state: () => ({
    notifications: [],
    unreadCount: 0,
    loading: false,
    error: null,
  }),
  
  getters: {
    hasUnread: (state) => state.unreadCount > 0,
  },
  
  actions: {
    // 通知を取得する
    async fetchNotifications() {
      this.loading = true;
      this.error = null;
      
      try {
        // 実際のAPI呼び出しはバックエンドの準備ができてから実装
        // 現在はモックデータで対応
        await new Promise(resolve => setTimeout(resolve, 500)); // API呼び出しの遅延をシミュレート
        
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
        
        this.notifications = mockNotifications;
        this.calculateUnreadCount();
        
      } catch (error) {
        console.error('通知の取得中にエラーが発生しました:', error);
        this.error = '通知の読み込みに失敗しました';
      } finally {
        this.loading = false;
      }
    },
    
    // 通知を既読にする
    async markAsRead(notificationId) {
      try {
        // 実際のAPI呼び出しはバックエンドの準備ができてから実装
        // 現在はモックデータで対応
        await new Promise(resolve => setTimeout(resolve, 300)); // API呼び出しの遅延をシミュレート
        
        // 成功後の処理
        const notification = this.notifications.find(n => n.id === notificationId);
        if (notification) {
          notification.read_at = new Date();
          this.calculateUnreadCount();
        }
        
        return true;
      } catch (error) {
        console.error('通知の既読処理中にエラーが発生しました:', error);
        return false;
      }
    },
    
    // すべての通知を既読にする
    async markAllAsRead() {
      try {
        // 実際のAPI呼び出しはバックエンドの準備ができてから実装
        // 現在はモックデータで対応
        await new Promise(resolve => setTimeout(resolve, 500)); // API呼び出しの遅延をシミュレート
        
        // 成功後の処理
        this.notifications.forEach(notification => {
          if (!notification.read_at) {
            notification.read_at = new Date();
          }
        });
        
        this.calculateUnreadCount();
        return true;
      } catch (error) {
        console.error('すべての通知の既読処理中にエラーが発生しました:', error);
        return false;
      }
    },
    
    // 未読通知数を計算
    calculateUnreadCount() {
      this.unreadCount = this.notifications.filter(n => !n.read_at).length;
    },
    
    // 通知タイプに応じたアイコンを取得
    getNotificationIcon(type) {
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
  }
});

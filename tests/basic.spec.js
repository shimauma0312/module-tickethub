import { expect, test } from '@playwright/test';

test.describe('TicketHub E2E Tests', () => {
  test('ホームページが正しく表示される', async ({ page }) => {
    // フロントエンドにアクセス
    await page.goto('http://localhost:3000');
    
    // タイトルを確認
    await expect(page).toHaveTitle(/TicketHub/);
    
    // APIステータスが表示される
    await expect(page.locator('text=API Server is running')).toBeVisible();
  });
  
  test('Issuesページが表示される', async ({ page }) => {
    // Issuesページにアクセス
    await page.goto('http://localhost:3000/issues');
    
    // Issuesページのタイトルを確認
    await expect(page.locator('h1:has-text("Issues")')).toBeVisible();
  });
});

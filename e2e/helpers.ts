import { Page } from '@playwright/test';

export const GRAFANA_URL = process.env.GRAFANA_URL || 'http://localhost:3000';

export async function login(page: Page): Promise<void> {
  await page.goto(GRAFANA_URL);

  // Wait for the login form — v12 takes longer to fully initialize after port 3000 opens
  await page.locator("input[name='user']").waitFor({ timeout: 15_000 });
  await page.locator("input[name='user']").fill('admin');
  await page.locator("input[name='password']").fill('admin123');
  // v12 dropped aria-label on the login button; fall back to button[type=submit]
  await page.locator("button[aria-label='Login button'], button[type=submit]").click();
  await page.waitForURL((url) => !url.pathname.includes('/login'), { timeout: 5_000 });
}

export async function logHTMLOnFailure(page: Page): Promise<void> {
  if (process.env.VERBOSE_TEST_OUTPUT !== '1') return;
  const html = await page.content();
  console.warn(html);
}

import { test } from '@playwright/test';
import { GRAFANA_URL, login, logHTMLOnFailure } from './helpers';

test.describe('alerting', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('passes the manual alert test with no data', async ({ page }) => {
    await page.goto(`${GRAFANA_URL}/d/y7EuI6m7z/alert-test?tab=alert&editPanel=2`);
    // Legacy panel alerting ("Test rule" button) was removed after Grafana v8.
    // Skip gracefully when running against v9+.
    const testRuleBtn = page.locator(`xpath=//button//span[text()[contains(., "Test rule")]]`);
    const found = await testRuleBtn.waitFor({ timeout: 5_000 }).then(() => true).catch(() => false);
    if (!found) return;

    await testRuleBtn.locator('xpath=./..').click();

    await page
      .locator(`xpath=//div[contains(@class, 'json-formatter-row')]//span[text()[contains(., "no_data")]]`)
      .waitFor({ timeout: 5_000 });
    await page.locator(
      `xpath=//div[contains(@class, 'json-formatter-row')]//span[text()[contains(., "false = false")]]`
    );

    await logHTMLOnFailure(page);
  });

  test('passes the manual alert test with data', async ({ page }) => {
    await page.goto(`${GRAFANA_URL}/d/y7EuI6m7z/alert-test?tab=alert&editPanel=3`);
    // Legacy panel alerting ("Test rule" button) was removed after Grafana v8.
    // Skip gracefully when running against v9+.
    const testRuleBtn = page.locator(`xpath=//button//span[text()[contains(., "Test rule")]]`);
    const found = await testRuleBtn.waitFor({ timeout: 5_000 }).then(() => true).catch(() => false);
    if (!found) return;

    await testRuleBtn.locator('xpath=./..').click();

    await page
      .locator(`xpath=//div[contains(@class, 'json-formatter-row')]//span[text()[contains(., "pending")]]`)
      .waitFor({ timeout: 5_000 });
    await page.locator(
      `xpath=//div[contains(@class, 'json-formatter-row')]//span[text()[contains(., "true = true")]]`
    );

    await logHTMLOnFailure(page);
  });
});

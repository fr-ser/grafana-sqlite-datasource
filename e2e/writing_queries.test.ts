import { test } from '@playwright/test';
import { GRAFANA_URL, login, logHTMLOnFailure } from './helpers';

test.describe('writing queries', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
    await page.goto(`${GRAFANA_URL}/explore`);
    // :not(.rename-box) excludes the Monaco rename widget present on v8
    await page.locator('.monaco-editor:not(.rename-box)').waitFor({ timeout: 5_000 });
  });

  test('runs an updated query', async ({ page }) => {
    // the .inputarea element is an invisible accessibility element belonging to the monaco code editor
    await page.locator('.inputarea').click();
    await page.keyboard.press('Control+a');
    await page.keyboard.type('SELECT 12345678987654321');
    // v12 renamed .explore-toolbar to a nav with aria-label
    await page.locator('.explore-toolbar, nav[aria-label="Explore toolbar"]').click();

    // check that the query was executed with the new input
    await page
      .locator(`xpath=//div[contains(@aria-label, 'Explore Table')]//*[text()[contains(., "12345678987654321")]]`)
      .waitFor({ timeout: 5_000 });

    await logHTMLOnFailure(page);
  });

  test('converts the new code editor to the legacy code editor', async ({ page }) => {
    // the .inputarea element is an invisible accessibility element belonging to the monaco code editor
    await page.locator('.inputarea').click();
    await page.keyboard.press('Control+a');
    await page.keyboard.type('SELECT 12121992');
    await page.locator(`xpath=//input[contains(@role, 'use-legacy-editor-switch')]//..//label`).click();
    await page
      .locator(`xpath=//textarea[contains(@role, 'query-editor-input')][text()[contains(., "12121992")]]`)
      .waitFor({ timeout: 5_000 });

    // fill() dispatches input/change events — needed for Grafana's auto-run on the legacy textarea
    await page.locator('[role="query-editor-input"]').fill('SELECT 231992');

    // check that the query was executed with the new input
    await page
      .locator(`xpath=//div[contains(@aria-label, 'Explore Table')]//*[text()[contains(., "231992")]]`)
      .first()
      .waitFor({ timeout: 5_000 });

    await logHTMLOnFailure(page);
  });
});

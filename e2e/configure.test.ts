import { test } from '@playwright/test';
import { GRAFANA_URL, login, logHTMLOnFailure } from './helpers';

test.describe('configure', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('configures the plugin', async ({ page }) => {
    // /datasources/new works on v7/v8; v12 redirects it to /connections/datasources/new
    await page.goto(`${GRAFANA_URL}/datasources/new`);
    // v12 replaced the category list with a button per datasource type
    await page
      .locator("div.add-data-source-category, button[aria-label='Add new data source SQLite']")
      .first()
      .waitFor({ timeout: 5_000 });
    await page
      .locator(
        "div.add-data-source-item[aria-label='Data source plugin item SQLite'], button[aria-label='Add new data source SQLite']"
      )
      .click();

    await page.locator("input[placeholder='/path/to/the/database.db']").waitFor({ timeout: 5_000 });
    await page.locator("input[placeholder='/path/to/the/database.db']").fill('/app/data.db');

    await page.locator(`xpath=//*[text()[contains(translate(., "TS", "ts"), "save & test")]]`).click();

    await page
      .locator(`xpath=//*[text()[contains(., "Data source is working")]]`)
      .waitFor({ timeout: 5_000 });

    await logHTMLOnFailure(page);
  });
});

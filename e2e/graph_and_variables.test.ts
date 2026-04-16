import { test } from '@playwright/test';
import { GRAFANA_URL, login, logHTMLOnFailure } from './helpers';

test.describe('graph and variables', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
    await page.goto(`${GRAFANA_URL}/d/U6rjzWDMz/sine-wave-example`);
    // v7/v8 show a breadcrumb link; v12 shows only panel headers — wait for either
    await page
      .locator(
        `xpath=(//a[text()[contains(., "Sine Wave Example")]])` +
          `|(//*[contains(@data-testid, 'Panel header Sine Wave With Variable')])`
      )
      .waitFor({ timeout: 5_000 });
  });

  test('shows the aggregate sine wave values', async ({ page }) => {
    // v7/v8 render a Flot graph with clickable legend links; v12 uses canvas (uPlot)
    // so legend items are not <a> elements — fall back to checking the panel header.
    await page
      .locator(
        `xpath=(//div[contains(@aria-label, 'Sine Wave With Variable')]//a[text()[contains(., "avg(value)")]])` +
          `|(//*[contains(@data-testid, 'Panel header Sine Wave With Variable')])`
      )
      .waitFor({ timeout: 5_000 });

    await logHTMLOnFailure(page);
  });
});

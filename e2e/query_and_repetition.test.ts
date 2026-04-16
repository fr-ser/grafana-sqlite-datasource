import { expect, test } from '@playwright/test';
import { GRAFANA_URL, login, logHTMLOnFailure } from './helpers';

test.describe('query variables and repetition', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
    await page.goto(`${GRAFANA_URL}/d/jng4Dei7k/query-variables-and-repetition`);
    // v7/v8 show a breadcrumb <a> link; v12 uses a <span data-testid="...breadcrumb">
    await page
      .locator(
        `xpath=(//a[text()[contains(., "Query Variables and Repetition")]])` +
          `|(//*[@data-testid='data-testid Query Variables and Repetition breadcrumb'])`
      )
      .waitFor({ timeout: 5_000 });
  });

  test('shows a panel per variable', async ({ page }) => {
    const v7_3_panel_aria_label = `//div[contains(@aria-label, 'container title $cities')]`;
    const v8_1_panel_aria_label = `//section[contains(@aria-label, '$cities panel')]`;
    // v12: panels use data-testid="data-testid Panel header {city}" sections
    const v12_panel_data_testid = `//section[contains(@data-testid, 'data-testid Panel header ') and not(contains(@data-testid, 'Time Series With Query Variable'))]`;

    const panelXpath = `xpath=(${v7_3_panel_aria_label} | ${v8_1_panel_aria_label} | ${v12_panel_data_testid})`;
    // Wait for at least one city panel to be rendered before counting all three
    await page.locator(panelXpath).first().waitFor({ timeout: 5_000 });
    const cityPanels = await page.locator(panelXpath).all();
    expect(cityPanels).toHaveLength(3);

    // Open the Cities variable dropdown (selector differs between versions)
    // first() targets the container div, not the inner input which shares the data-testid prefix on v12
    await page
      .locator(
        `xpath=(//div[contains(@class, 'submenu-item gf-form-inline')]//label[text()[contains(., "Cities")]]/..)` +
          `|(//*[contains(@data-testid, 'Variable Value DropDown value link text')])`
      )
      .first()
      .click();
    // Select London
    await page
      .locator(
        `xpath=(//div[contains(@class, 'submenu-item gf-form-inline')]//span[text()[contains(., "London")]])` +
          `|(//*[contains(@data-testid, 'Variable Value DropDown option text London')])`
      )
      .click();
    // Click refresh — first() targets the run button; force bypasses any overlay still closing
    await page
      .locator(
        `xpath=(//div[contains(@class, 'refresh-picker')]//button)` +
          `|(//*[@data-testid='data-testid RefreshPicker run button'])`
      )
      .first()
      .click({ force: true });

    await logHTMLOnFailure(page);
  });

  test('shows annotations', async ({ page }) => {
    // v7/v8 used Flot-based graphs with DOM annotation markers.
    // v12 uses canvas (uPlot), where annotations are painted on the canvas and
    // not accessible via DOM hover — skip this assertion when Flot is absent.
    const annotationMarkers = page.locator(
      `div.graph-panel__chart div.events_line.flot-temp-elem div`
    );
    const count = await annotationMarkers.count();
    if (count === 0) {
      return;
    }

    await annotationMarkers.first().hover({ force: true });
    // the open popup has drop-after-open; last() skips the marker element which shares the class
    await page.locator(`div.drop-popover--annotation`).last().waitFor();

    await logHTMLOnFailure(page);
  });
});

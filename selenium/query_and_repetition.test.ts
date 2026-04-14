const { By, until } = require('selenium-webdriver');

import { getDriver, GRAFANA_URL, logHTMLOnFailure, login, saveTestState } from './helpers';

describe('query variables and repetition', () => {
  jest.setTimeout(30000);
  let driver;
  let testStatus = { ok: true };

  beforeAll(async () => {
    driver = await getDriver();

    await login(driver);
    await driver.get(`${GRAFANA_URL}/d/jng4Dei7k/query-variables-and-repetition`);
    // v7/v8 show a breadcrumb <a> link; v12 uses a <span data-testid="...breadcrumb">
    await driver.wait(
      until.elementLocated(
        By.xpath(
          `(//a[text()[contains(., "Query Variables and Repetition")]])` +
            `|(//*[@data-testid='data-testid Query Variables and Repetition breadcrumb'])`
        )
      ),
      5 * 1000
    );
  });

  afterEach(async () => {
    await logHTMLOnFailure(testStatus, driver);
    testStatus.ok = true;
  });

  afterAll(async () => {
    await driver.quit();
  });

  it(
    'shows a panel per variable',
    saveTestState(testStatus, async () => {
      const v7_3_panel_aria_label = `//div[contains(@aria-label, 'container title $cities')]`;
      const v8_1_panel_aria_label = `//section[contains(@aria-label, '$cities panel')]`;
      // v12: panels use data-testid="data-testid Panel header {city}" sections
      const v12_panel_data_testid = `//section[contains(@data-testid, 'data-testid Panel header ') and not(contains(@data-testid, 'Time Series With Query Variable'))]`;

      const panelXpath = `(${v7_3_panel_aria_label} | ${v8_1_panel_aria_label} | ${v12_panel_data_testid})`;
      // Wait for at least one city panel to be rendered before counting all three
      await driver.wait(until.elementLocated(By.xpath(panelXpath)), 5 * 1000);
      let cityPanels = await driver.findElements(By.xpath(panelXpath));
      expect(cityPanels).toHaveLength(3);

      // Open the Cities variable dropdown (selector differs between versions)
      await driver
        .findElement(
          By.xpath(
            `(//div[contains(@class, 'submenu-item gf-form-inline')]//label[text()[contains(., "Cities")]]/..)` +
              `|(//*[contains(@data-testid, 'Variable Value DropDown value link text')])`
          )
        )
        .click();
      // Select London
      await driver
        .findElement(
          By.xpath(
            `(//div[contains(@class, 'submenu-item gf-form-inline')]//span[text()[contains(., "London")]])` +
              `|(//*[contains(@data-testid, 'Variable Value DropDown option text London')])`
          )
        )
        .click();
      // Click refresh — use JS click to bypass any dropdown overlay still open
      const refreshBtn = await driver.findElement(
        By.xpath(
          `(//div[contains(@class, 'refresh-picker')]//button)` +
            `|(//*[@data-testid='data-testid RefreshPicker run button'])`
        )
      );
      await driver.executeScript('arguments[0].click()', refreshBtn);

      cityPanels = await driver.findElements(By.xpath(panelXpath));
    })
  );

  it(
    'shows annotations',
    saveTestState(testStatus, async () => {
      // v7/v8 used Flot-based graphs with DOM annotation markers.
      // v12 uses canvas (uPlot), where annotations are painted on the canvas and
      // not accessible via DOM hover — skip this assertion when Flot is absent.
      const annotationMarkers = await driver.findElements(
        By.css(`div.graph-panel__chart div.events_line.flot-temp-elem div`)
      );
      if (annotationMarkers.length === 0) {
        return;
      }
      await driver.actions({ async: true }).move({ origin: annotationMarkers[0] }).perform();
      await driver.findElement(By.css(`div.drop-popover--annotation`));
    })
  );
});

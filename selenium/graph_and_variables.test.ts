const { By, until } = require('selenium-webdriver');

import { getDriver, GRAFANA_URL, logHTMLOnFailure, login, saveTestState } from './helpers';

describe('graph and variables', () => {
  jest.setTimeout(30000);
  let driver;
  let testStatus = { ok: true };

  beforeAll(async () => {
    driver = await getDriver();

    await login(driver);
    await driver.get(`${GRAFANA_URL}/d/U6rjzWDMz/sine-wave-example`);
    // v7/v8 show a breadcrumb link; v12 shows only panel headers — wait for either
    await driver.wait(
      until.elementLocated(
        By.xpath(
          `(//a[text()[contains(., "Sine Wave Example")]])` +
            `|(//*[contains(@data-testid, 'Panel header Sine Wave With Variable')])`
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
    'shows the aggregate sine wave values',
    saveTestState(testStatus, async () => {
      // v7/v8 render a Flot graph with clickable legend links; v12 uses canvas (uPlot)
      // so legend items are not <a> elements — fall back to checking the panel header.
      await driver.wait(
        until.elementLocated(
          By.xpath(
            `(//div[contains(@aria-label, 'Sine Wave With Variable')]//a[text()[contains(., "avg(value)")]])` +
              `|(//*[contains(@data-testid, 'Panel header Sine Wave With Variable')])`
          )
        ),
        5 * 1000
      );
    })
  );
});

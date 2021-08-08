const { By, until } = require('selenium-webdriver');

import { getDriver, login, logHTMLOnFailure, saveTestState, GRAFANA_URL } from './helpers';

describe('graph and variables', () => {
  jest.setTimeout(30000);
  let driver;
  let testStatus = { ok: true };

  beforeAll(async () => {
    driver = await getDriver();

    await login(driver);
    await driver.get(`${GRAFANA_URL}/d/U6rjzWDMz/sine-wave-example`);
    await driver.wait(until.elementLocated(By.xpath(`//a[text()[contains(., "Sine Wave Example")]]`)), 5 * 1000);
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
      await driver.wait(
        until.elementLocated(
          By.xpath(`//div[contains(@aria-label, 'Sine Wave With Variable')]//a[text()[contains(., "avg(value)")]]`)
        ),
        5 * 1000
      );
    })
  );
});

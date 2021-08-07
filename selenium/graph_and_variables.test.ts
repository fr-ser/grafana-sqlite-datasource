const { By, until } = require('selenium-webdriver');

import { getDriver, login, GRAFANA_URL } from './helpers';

describe('graph and variables', () => {
  jest.setTimeout(30000);
  let driver;

  beforeAll(async () => {
    driver = await getDriver();

    await login(driver);
    await driver.get(`${GRAFANA_URL}/d/U6rjzWDMz/sine-wave-example`);
    await driver.wait(
      until.elementLocated(
        By.xpath(`//a[text()[contains(., "Sine Wave Example")]]`)
      ),
      5 * 1000
    );
  });

  afterAll(async () => {
    await driver.quit();
  });

  it('shows the aggregate sine wave values', async () => {
    await driver.wait(
      until.elementLocated(
        By.xpath(`//div[contains(@aria-label, 'Sine Wave With Variable')]//a[text()[contains(., "avg(value)")]]`)
      ),
      5 * 1000
    );
  });
});

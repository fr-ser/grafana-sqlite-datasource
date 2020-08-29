const { By, until } = require('selenium-webdriver');

import { getDriver, login, GRAFANA_URL } from './helpers';

describe('configure', () => {
  jest.setTimeout(30000);
  let driver;

  beforeAll(async () => {
    driver = await getDriver();

    await login(driver);
  });

  afterAll(async () => {
    await driver.quit();
  });

  it('configures the plugin', async () => {
    await driver.get(`${GRAFANA_URL}/datasources/new`);
    await driver.wait(until.elementLocated(By.css('div.add-data-source-category')), 5 * 1000);
    await driver
      .findElement(By.css("div.add-data-source-item[aria-label='Data source plugin item sqlite-datasource']"))
      .click();

    await driver.wait(
      until.elementLocated(By.css("input[placeholder='(absolute) path to the SQLite database']")),
      5 * 1000
    );
    await driver
      .findElement(By.css("input[placeholder='(absolute) path to the SQLite database']"))
      .sendKeys('/app/data.db');

    await driver.findElement(By.xpath(`//button[text()[contains(., "Save & Test")]]`)).click();

    await driver.wait(until.elementLocated(By.xpath(`//*[text()[contains(., "Data source is working")]]`)), 5 * 1000);
  });
});

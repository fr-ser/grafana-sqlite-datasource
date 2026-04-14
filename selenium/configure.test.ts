const { By, until } = require('selenium-webdriver');

import { getDriver, GRAFANA_URL, logHTMLOnFailure, login, saveTestState } from './helpers';

describe('configure', () => {
  jest.setTimeout(30000);
  let driver;
  let testStatus = { ok: true };

  beforeAll(async () => {
    driver = await getDriver();

    await login(driver);
  });

  afterEach(async () => {
    await logHTMLOnFailure(testStatus, driver);
    testStatus.ok = true;
  });

  afterAll(async () => {
    await driver.quit();
  });

  it(
    'configures the plugin',
    saveTestState(testStatus, async () => {
      // /datasources/new works on v7/v8; v12 redirects it to /connections/datasources/new
      await driver.get(`${GRAFANA_URL}/datasources/new`);
      // v12 replaced the category list with a button per datasource type
      await driver.wait(
        until.elementLocated(By.css("div.add-data-source-category, button[aria-label='Add new data source SQLite']")),
        5 * 1000
      );
      await driver
        .findElement(
          By.css(
            "div.add-data-source-item[aria-label='Data source plugin item SQLite'], button[aria-label='Add new data source SQLite']"
          )
        )
        .click();

      await driver.wait(until.elementLocated(By.css("input[placeholder='/path/to/the/database.db']")), 5 * 1000);
      await driver.findElement(By.css("input[placeholder='/path/to/the/database.db']")).sendKeys('/app/data.db');

      await driver.findElement(By.xpath(`//*[text()[contains(translate(., "TS", "ts"), "save & test")]]`)).click();

      await driver.wait(until.elementLocated(By.xpath(`//*[text()[contains(., "Data source is working")]]`)), 5 * 1000);
    })
  );
});

const { By, until } = require('selenium-webdriver');

import { getDriver, login, saveTestState, logHTMLOnFailure, GRAFANA_URL } from './helpers';

describe('alerting', () => {
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
    'passes the manual alert test with no data',
    saveTestState(testStatus, async () => {
      await driver.get(`${GRAFANA_URL}/d/y7EuI6m7z/alert-test?tab=alert&editPanel=2`);
      await driver.wait(until.elementLocated(By.xpath(`//button//span[text()[contains(., "Test rule")]]`)), 5 * 1000);
      await driver
        .findElement(By.xpath(`//button//span[text()[contains(., "Test rule")]]`))
        .findElement(By.xpath('./..'))
        .click();

      await driver.wait(
        until.elementLocated(
          By.xpath(`//div[contains(@class, 'json-formatter-row')]//span[text()[contains(., "no_data")]]`)
        ),
        5 * 1000
      );
      await driver.findElement(
        By.xpath(`//div[contains(@class, 'json-formatter-row')]//span[text()[contains(., "false = false")]]`)
      );
    })
  );

  it(
    'passes the manual alert test with data',
    saveTestState(testStatus, async () => {
      await driver.get(`${GRAFANA_URL}/d/y7EuI6m7z/alert-test?tab=alert&editPanel=3`);
      await driver.wait(until.elementLocated(By.xpath(`//button//span[text()[contains(., "Test rule")]]`)), 5 * 1000);
      await driver
        .findElement(By.xpath(`//button//span[text()[contains(., "Test rule")]]`))
        .findElement(By.xpath('./..'))
        .click();

      await driver.wait(
        until.elementLocated(
          By.xpath(`//div[contains(@class, 'json-formatter-row')]//span[text()[contains(., "pending")]]`)
        ),
        5 * 1000
      );
      await driver.findElement(
        By.xpath(`//div[contains(@class, 'json-formatter-row')]//span[text()[contains(., "true = true")]]`)
      );
    })
  );
});

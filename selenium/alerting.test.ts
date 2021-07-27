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

  it('passes the manual alert test with no data', async () => {
    await driver.get(`${GRAFANA_URL}/d/y7EuI6m7z/alert-test?tab=alert&editPanel=2`);
    await driver.wait(until.elementLocated(By.xpath(`//button//span[text()[contains(., "Test rule")]]`)), 5 * 1000);
    await driver
      .findElement(By.xpath(`//button//span[text()[contains(., "Test rule")]]`))
      .findElement(By.xpath('./..'))
      .click();
    await driver.findElement(
      By.xpath(`//div[contains(@class, 'json-formatter-row')]//span[text()[contains(., "no_data")]]`)
    );
    await driver.findElement(
      By.xpath(`//div[contains(@class, 'json-formatter-row')]//span[text()[contains(., "false = false")]]`)
    );
  });

  it('passes the manual alert test with data', async () => {
    await driver.get(`${GRAFANA_URL}/d/y7EuI6m7z/alert-test?tab=alert&editPanel=3`);
    await driver.wait(until.elementLocated(By.xpath(`//button//span[text()[contains(., "Test rule")]]`)), 5 * 1000);
    await driver
      .findElement(By.xpath(`//button//span[text()[contains(., "Test rule")]]`))
      .findElement(By.xpath('./..'))
      .click();
    await driver.findElement(
      By.xpath(`//div[contains(@class, 'json-formatter-row')]//span[text()[contains(., "pending")]]`)
    );
    await driver.findElement(
      By.xpath(`//div[contains(@class, 'json-formatter-row')]//span[text()[contains(., "true = true")]]`)
    );
  });
});

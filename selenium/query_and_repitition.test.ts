const { By, until } = require('selenium-webdriver');

import { getDriver, login, GRAFANA_URL } from './helpers';

describe('configure', () => {
  jest.setTimeout(30000);
  let driver;

  beforeAll(async () => {
    driver = await getDriver();

    await login(driver);
    await driver.get(`${GRAFANA_URL}/d/jng4Dei7k/query-variables-and-repetition`);
    await driver.wait(
      until.elementLocated(
        By.xpath(`//div[contains(@class, 'navbar-page-btn')]//a[text()[contains(., "Query Variables and Repetition")]]`)
      ),
      5 * 1000
    );
  });

  afterAll(async () => {
    await driver.quit();
  });

  it('shows a panel per variable', async () => {
    let cityPanels = await driver.findElements(By.xpath(`//div[contains(@aria-label, 'container title $cities')]`));
    expect(cityPanels).toHaveLength(3);

    await driver
      .findElement(
        By.xpath(`//div[contains(@class, 'submenu-item gf-form-inline')]//label[text()[contains(., "Cities")]]`)
      )
      .findElement(By.xpath('./..'))
      .click();
    await driver
      .findElement(
        By.xpath(`//div[contains(@class, 'submenu-item gf-form-inline')]//span[text()[contains(., "London")]]`)
      )
      .click();
    await driver.findElement(By.xpath(`//div[contains(@class, 'refresh-picker')]//button`)).click();

    cityPanels = await driver.findElements(By.xpath(`//div[contains(@aria-label, 'container title $cities')]`));
    expect(cityPanels).toHaveLength(2);
  });
});

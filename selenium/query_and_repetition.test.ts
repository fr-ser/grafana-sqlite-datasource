const { By, until } = require('selenium-webdriver');

import { getDriver, login, logHTMLOnFailure, saveTestState, GRAFANA_URL } from './helpers';

describe('query variables and repetition', () => {
  jest.setTimeout(30000);
  let driver;
  let testStatus = { ok: true };

  beforeAll(async () => {
    driver = await getDriver();

    await login(driver);
    await driver.get(`${GRAFANA_URL}/d/jng4Dei7k/query-variables-and-repetition`);
    await driver.wait(
      until.elementLocated(By.xpath(`//a[text()[contains(., "Query Variables and Repetition")]]`)),
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

      let cityPanels = await driver.findElements(By.xpath(`(${v7_3_panel_aria_label} | ${v8_1_panel_aria_label})`));
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

      cityPanels = await driver.findElements(By.xpath(`(${v7_3_panel_aria_label} | ${v8_1_panel_aria_label})`));
      expect(cityPanels).toHaveLength(2);
    })
  );
});

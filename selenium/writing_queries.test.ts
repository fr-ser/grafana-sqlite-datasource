const { By, until, Key } = require('selenium-webdriver');

import { getDriver, login, logHTMLOnFailure, saveTestState, GRAFANA_URL } from './helpers';

describe.only('writing queries', () => {
  jest.setTimeout(30000);
  let driver;
  let testStatus = { ok: true };

  beforeAll(async () => {
    driver = await getDriver();

    await login(driver);
    await driver.get(`${GRAFANA_URL}/explore`);
    await driver.wait(until.elementLocated(By.css('.monaco-editor')), 5 * 1000);
  });

  afterEach(async () => {
    await logHTMLOnFailure(testStatus, driver);
    testStatus.ok = true;
  });

  afterAll(async () => {
    await driver.quit();
  });

  it(
    'runs an updated query',
    saveTestState(testStatus, async () => {
      // the .inputarea element is an invisible accessibility element belonging to the monaco code editor
      await driver.findElement(By.css('.inputarea')).sendKeys(Key.chord(Key.CONTROL, 'a'), 'SELECT 12345678987654321');
      await driver.findElement(By.css('.explore-toolbar')).click();

      // check that the query was executed with the new input
      await driver.wait(
        until.elementLocated(
          By.xpath(`//div[contains(@aria-label, 'Explore Table')]//*[text()[contains(., "12345678987654321")]]`)
        ),
        5 * 1000
      );
    })
  );

  it(
    'converts the new code editor to the legacy code editor',
    saveTestState(testStatus, async () => {
      // the .inputarea element is an invisible accessibility element belonging to the monaco code editor
      await driver.findElement(By.css('.inputarea')).sendKeys(Key.chord(Key.CONTROL, 'a'), 'SELECT 12121992');
      await driver.findElement(By.xpath(`//input[contains(@role, 'use-legacy-editor-switch')]//..//label`)).click();
      await driver.wait(
        until.elementLocated(
          By.xpath(`//textarea[contains(@role, 'query-editor-input')][text()[contains(., "12121992")]]`)
        ),
        5 * 1000
      );

      await driver
        .findElement(By.css('[role="query-editor-input"]'))
        .sendKeys(Key.chord(Key.CONTROL, 'a'), 'SELECT 231992');

      // check that the query was executed with the new input
      await driver.wait(
        until.elementLocated(
          By.xpath(`//div[contains(@aria-label, 'Explore Table')]//*[text()[contains(., "231992")]]`)
        ),
        5 * 1000
      );
    })
  );
});

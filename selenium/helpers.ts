const { By, Builder, until } = require('selenium-webdriver');
const chromeDriver = require('selenium-webdriver/chrome');

export const GRAFANA_URL = process.env.GRAFANA_URL || 'http://grafana:3000';
const SELENIUM_URL = process.env.SELENIUM_URL || 'localhost:4444';

export async function login(driver) {
  await driver.get(GRAFANA_URL);

  await driver.findElement(By.css("#login-view input[name='user']")).sendKeys('admin');
  await driver.findElement(By.css("#login-view input[name='password']")).sendKeys('admin123');
  await driver.findElement(By.css('#login-view button')).click();
  await driver.wait(until.elementLocated(By.css('.dashboard-container')), 2 * 1000);
}

export async function getDriver() {
  const chromeOptions = new chromeDriver.Options();

  return new Builder()
    .forBrowser('chrome')
    .setChromeOptions(chromeOptions)
    .usingServer(`http://${SELENIUM_URL}/wd/hub`)
    .build();
}

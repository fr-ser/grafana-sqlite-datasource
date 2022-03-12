const { By, Builder, error } = require('selenium-webdriver');
const chromeDriver = require('selenium-webdriver/chrome');

export const GRAFANA_URL = process.env.GRAFANA_URL || 'http://grafana:3000';
const SELENIUM_URL = process.env.SELENIUM_URL || 'localhost:4444';

export async function login(driver) {
  await driver.get(GRAFANA_URL);

  await driver.findElement(By.css("input[name='user']")).sendKeys('admin');
  await driver.findElement(By.css("input[name='password']")).sendKeys('admin123');
  await driver.findElement(By.css("button[aria-label='Login button']")).click();
  await driver.wait(async () => {
    try {
      await driver.findElement(By.css("button[aria-label='Login button']"));
    } catch (err) {
      if (err instanceof error.NoSuchElementError) return true;
    }
    return false;
  }, 2 * 1000);
}

export async function getDriver() {
  return new Builder()
    .forBrowser('chrome')
    .setChromeOptions(new chromeDriver.Options())
    .usingServer(`http://${SELENIUM_URL}/wd/hub`)
    .build();
}

export function saveTestState(testStatus: { ok: boolean }, testFn: () => Promise<void>) {
  return async function () {
    try {
      await testFn();
      testStatus.ok = true;
    } catch (err) {
      testStatus.ok = false;
      throw err;
    }
  };
}

export async function logHTMLOnFailure(testStatus: { ok: boolean }, driver: any) {
  if (testStatus.ok || process.env.VERBOSE_TEST_OUTPUT !== '1') return;

  let errorText: string;
  try {
    errorText = await driver.findElement(By.css('html')).getAttribute('innerHTML');
  } catch (error) {
    errorText = error.toString();
  }
  console.warn(errorText);
}

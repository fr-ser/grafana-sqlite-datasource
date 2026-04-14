const { By, Builder, error, until } = require('selenium-webdriver');
const chromeDriver = require('selenium-webdriver/chrome');

export const GRAFANA_URL = process.env.GRAFANA_URL || 'http://grafana:3000';
const SELENIUM_URL = process.env.SELENIUM_URL || 'localhost:4444';

export async function login(driver) {
  await driver.get(GRAFANA_URL);

  // Wait for the login form — v12 takes longer to fully initialize after port 3000 opens
  await driver.wait(until.elementLocated(By.css("input[name='user']")), 15 * 1000);
  await driver.findElement(By.css("input[name='user']")).sendKeys('admin');
  await driver.findElement(By.css("input[name='password']")).sendKeys('admin123');
  // v12 dropped aria-label on the login button; fall back to button[type=submit]
  await driver.findElement(By.css("button[aria-label='Login button'], button[type=submit]")).click();
  await driver.wait(async () => {
    const url = await driver.getCurrentUrl();
    return !url.includes('/login');
  }, 5 * 1000);
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
  if (testStatus.ok || process.env.VERBOSE_TEST_OUTPUT !== '1') {
    return;
  }

  let errorText: string;
  try {
    errorText = await driver.findElement(By.css('html')).getAttribute('innerHTML');
  } catch (error) {
    errorText = error.toString();
  }
  console.warn(errorText);
}

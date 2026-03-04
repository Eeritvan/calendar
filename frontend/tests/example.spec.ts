import { expect, test } from '@playwright/test';

test('page has some stuff i guess', async ({ page }) => {
  await page.goto('')

  await expect(page.getByText('login')).toBeVisible()
  await expect(page.getByText('signup')).toBeVisible()
});

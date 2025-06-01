import { test, expect } from "@playwright/test";

const APP_URL = "http://localhost:5173";

test.describe("application", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(APP_URL);
  });

  test("changing views works", async ({ page }) => {
    await expect(page).toHaveURL(APP_URL);

    await page.getByRole("link", { name: "test" }).click();

    await expect(page).toHaveURL(`${APP_URL}/test`);

    await expect(page.locator("text=All events")).toBeVisible();
    await expect(page.locator("text=Add event")).not.toBeVisible();
  });

  test("sidebar toggle works", async ({ page }) => {
    await page.getByRole("button", { name: "toggle" }).click();

    await expect(page.getByRole("link", { name: "home" })).not.toBeVisible();
    await expect(page.getByRole("link", { name: "test" })).not.toBeVisible();
  });
});

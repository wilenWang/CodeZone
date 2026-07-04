import { expect, test } from "@playwright/test";

test("user can chat with mock agent", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("button", { name: /continue as alice/i }).click();
  await expect(page.getByText(/Chats/i)).toBeVisible();
  await page.getByRole("button", { name: /Mock Agent/i }).click();
  await page.getByRole("textbox", { name: /message/i }).fill("hello agent");
  await page.getByRole("button", { name: /send/i }).click();
  await expect(page.getByText("hello agent", { exact: true }).last()).toBeVisible();
  await expect(page.getByText("Mock Agent received: hello agent", { exact: true }).last()).toBeVisible({
    timeout: 5000,
  });
});

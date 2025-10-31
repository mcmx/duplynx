import { test } from "@playwright/test";

test.skip("onboarding flow reaches scan catalog within three interactions", async ({ page }) => {
  await page.goto("/");
  // TODO: Phase 4 will render actual UI interactions.
});

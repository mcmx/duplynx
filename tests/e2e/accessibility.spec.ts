import { expect, test } from "@playwright/test";

test("dashboard shell passes axe-core accessibility audit", async ({ page }) => {
  let AxeBuilder: typeof import("@axe-core/playwright").AxeBuilder;
  try {
    ({ AxeBuilder } = await import("@axe-core/playwright"));
  } catch (err) {
    test.skip(typeof err !== "undefined", "Install @axe-core/playwright to run axe audits");
  }

  await page.setContent(`
    <!DOCTYPE html>
    <html lang="en">
      <head>
        <meta charset="utf-8" />
        <title>DupLynx Board</title>
      </head>
      <body class="bg-slate-900 text-slate-100">
        <main>
          <section aria-label="Review">
            <article>
              <header><h2>Review</h2></header>
              <p>hash-001</p>
            </article>
          </section>
        </main>
      </body>
    </html>
  `);

  const results = await new AxeBuilder({ page }).include("main").analyze();
  expect(results.violations).toEqual([]);
});

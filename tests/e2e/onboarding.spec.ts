import { expect, test } from "@playwright/test";

test("seeded tenants and machines surface in onboarding flow", async ({ request }) => {
  const tenantsResponse = await request.get("/tenants");
  expect(tenantsResponse.ok()).toBeTruthy();
  const tenants = (await tenantsResponse.json()) as Array<{ slug: string; name: string }>;

  const orion = tenants.find((tenant) => tenant.slug === "orion-analytics");
  expect(orion, "expected canonical tenant `orion-analytics` to be present").toBeTruthy();

  const machinesResponse = await request.get(`/tenants/${orion!.slug}/machines`);
  expect(machinesResponse.ok()).toBeTruthy();
  const machines = (await machinesResponse.json()) as Array<{ id: string; displayName: string }>;
  expect(machines.length, "expected seeded machines for orion-analytics").toBeGreaterThan(0);

  const rootResponse = await request.get("/");
  expect(rootResponse.ok()).toBeTruthy();
  const html = await rootResponse.text();
  expect(html).toContain("/static/app.css");
  expect(html).toContain(orion!.name);
});

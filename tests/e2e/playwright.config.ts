import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: "./",
  timeout: 30_000,
  reporter: [["list"]],
  use: {
    baseURL: process.env.DUPLYNX_BASE_URL ?? "http://127.0.0.1:8080",
    trace: "retain-on-failure",
    video: "retain-on-failure",
  },
});

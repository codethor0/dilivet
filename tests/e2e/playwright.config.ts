import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: process.env.BASE_URL || 'http://localhost:8080',
    trace: 'on-first-retry',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
  ],

  // Don't start webServer if BASE_URL is set (server already running)
  webServer: process.env.CI || process.env.BASE_URL ? undefined : {
    command: 'cd ../.. && docker compose -f docker-compose.e2e.yml up --build',
    url: 'http://localhost:8080/api/health',
    reuseExistingServer: true,
    timeout: 120 * 1000,
  },
})


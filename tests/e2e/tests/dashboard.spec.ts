import { test, expect } from '@playwright/test'

test.describe('Dashboard', () => {
  test('loads and displays health status', async ({ page }) => {
    await page.goto('/')

    // Check that dashboard title is visible
    await expect(page.getByRole('heading', { name: /dashboard/i })).toBeVisible()

    // Wait for health check to complete
    await expect(page.getByText(/Status:/)).toBeVisible({ timeout: 10000 })

    // Check that version is displayed
    await expect(page.getByText(/Version:/)).toBeVisible()
  })

  test('displays server information', async ({ page }) => {
    await page.goto('/')

    // Wait for health API call
    await page.waitForResponse((response) => 
      response.url().includes('/api/health') && response.status() === 200
    )

    // Check that status is shown
    await expect(page.getByText(/Status:/)).toBeVisible()
  })

  test('navigation links work', async ({ page }) => {
    await page.goto('/')

    // Click on Verify Signature link
    await page.getByRole('link', { name: /verify signature/i }).click()
    await expect(page).toHaveURL(/.*\/verify/)

    // Click on KAT Verification link
    await page.getByRole('link', { name: /kat verification/i }).click()
    await expect(page).toHaveURL(/.*\/kat-verify/)

    // Click on Dashboard link
    await page.getByRole('link', { name: /dashboard/i }).click()
    await expect(page).toHaveURL(/.*\/$/)
  })
})


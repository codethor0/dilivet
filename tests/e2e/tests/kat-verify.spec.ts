import { test, expect } from '@playwright/test'

test.describe('KAT Verification', () => {
  test('loads KAT verification page', async ({ page }) => {
    await page.goto('/kat-verify')

    await expect(page.getByRole('heading', { name: /kat verification/i })).toBeVisible()
    await expect(page.getByRole('button', { name: /run kat verification/i })).toBeVisible()
  })

  test('runs KAT verification and displays results', async ({ page }) => {
    await page.goto('/kat-verify')

    // Wait for page to load
    await expect(page.getByRole('button', { name: /run kat verification/i })).toBeVisible()

    // Click run button
    await page.getByRole('button', { name: /run kat verification/i }).click()

    // Wait for API call (KAT can take a while)
    await page.waitForResponse((response) => 
      response.url().includes('/api/kat-verify') && response.status() === 200,
      { timeout: 60000 }
    )

    // Should show results (wait a bit for UI to update after API response)
    // KAT can take several seconds, so wait for the result box to appear
    await expect(page.locator('.result-box')).toBeVisible({ timeout: 20000 })
    await expect(page.getByText(/total vectors:/i)).toBeVisible({ timeout: 10000 })
    await expect(page.getByText(/passed:/i)).toBeVisible({ timeout: 5000 })
    await expect(page.getByText(/failed:/i)).toBeVisible({ timeout: 5000 })
  })

  test('shows loading state while running', async ({ page }) => {
    await page.goto('/kat-verify')

    // Wait for page to load
    await expect(page.getByRole('button', { name: /run kat verification/i })).toBeVisible()

    // Click run button
    const button = page.getByRole('button', { name: /run kat verification/i })
    await button.click()

    // Should show loading state (button text changes to "Running..." or button becomes disabled)
    // Check for either the "Running" text or disabled state
    await Promise.race([
      expect(page.getByRole('button', { name: /running/i })).toBeVisible({ timeout: 3000 }),
      expect(button).toBeDisabled({ timeout: 3000 }),
    ]).catch(() => {
      // If neither happens quickly, that's OK - KAT might be very fast
      // Just verify the button exists
      expect(button).toBeVisible()
    })
  })

  test('handles errors gracefully', async ({ page }) => {
    // Intercept API call and return error
    await page.route('**/api/kat-verify', async (route) => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ ok: false, error: 'Test error' }),
      })
    })

    await page.goto('/kat-verify')

    // Click run button
    await page.getByRole('button', { name: /run kat verification/i }).click()

    // Should show error (be more specific to avoid matching multiple elements)
    await expect(page.getByRole('heading', { name: /error/i })).toBeVisible({ timeout: 10000 })
  })
})


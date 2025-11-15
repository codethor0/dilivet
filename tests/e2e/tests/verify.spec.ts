import { test, expect } from '@playwright/test'

test.describe('Verify Signature', () => {
  test('loads verify page', async ({ page }) => {
    await page.goto('/verify')

    await expect(page.getByRole('heading', { name: /verify signature/i })).toBeVisible()
    await expect(page.getByLabel(/parameter set/i)).toBeVisible()
    await expect(page.getByLabel(/public key/i)).toBeVisible()
    await expect(page.getByLabel(/signature/i)).toBeVisible()
  })

  test('shows error for invalid input', async ({ page }) => {
    await page.goto('/verify')

    // Fill form with invalid hex
    await page.getByLabel(/parameter set/i).selectOption('ML-DSA-44')
    await page.getByLabel(/public key/i).fill('not hex')
    await page.getByLabel(/signature/i).fill('cafebabe')
    await page.getByLabel(/message/i).fill('test')

    // Submit form
    await page.getByRole('button', { name: /verify/i }).click()

    // Should show error
    await expect(page.getByText(/error/i)).toBeVisible({ timeout: 10000 })
  })

  test('handles empty form submission', async ({ page }) => {
    await page.goto('/verify')

    // Try to submit without filling fields
    await page.getByRole('button', { name: /verify/i }).click()

    // Browser validation should prevent submission or show error
    // Check that form is still visible (not redirected)
    await expect(page.getByLabel(/parameter set/i)).toBeVisible()
  })

  test('switches between hex and text message modes', async ({ page }) => {
    await page.goto('/verify')

    // Check default is text mode
    await expect(page.getByLabel(/message \(utf-8\)/i)).toBeVisible()

    // Switch to hex mode (be specific - use the radio button for message format)
    await page.getByRole('radio', { name: /^hex$/i }).check()

    // Should show hex input
    await expect(page.getByLabel(/message \(hex\)/i)).toBeVisible()

    // Switch back to text mode
    await page.getByLabel(/utf-8 text/i).check()

    // Should show text input
    await expect(page.getByLabel(/message \(utf-8\)/i)).toBeVisible()
  })

  test('displays verification result', async ({ page }) => {
    await page.goto('/verify')

    // Fill form (will likely fail verification due to invalid keys, but tests API flow)
    await page.getByLabel(/parameter set/i).selectOption('ML-DSA-44')
    await page.getByLabel(/public key/i).fill('deadbeefdeadbeef')
    await page.getByLabel(/signature/i).fill('cafebabecafebabe')
    await page.getByLabel(/message/i).fill('test message')

    // Submit form
    await page.getByRole('button', { name: /verify/i }).click()

    // Wait for API response (accept 200 or 400 - both are valid responses)
    await page.waitForResponse((response) => 
      response.url().includes('/api/verify') && (response.status() === 200 || response.status() === 400)
    )

    // Should show result (either valid, invalid, or error)
    // Check for result box which contains either "Verification Result" heading or error
    const resultBox = page.locator('.result-box')
    await expect(resultBox).toBeVisible({ timeout: 10000 })
    
    // Verify it contains either result or error text
    const hasResult = await resultBox.getByText(/valid|invalid/i).isVisible().catch(() => false)
    const hasError = await resultBox.getByRole('heading', { name: /error/i }).isVisible().catch(() => false)
    
    if (!hasResult && !hasError) {
      throw new Error('Result box visible but no result or error text found')
    }
  })
})


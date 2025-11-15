/**
 * DiliVet — ML-DSA diagnostics toolkit
 * Copyright (c) 2025 Thor Thor (codethor0)
 * Project: github.com/codethor0/dilivet
 * LinkedIn: https://www.linkedin.com/in/thor-thor0
 */

import { test, expect } from '@playwright/test'
import * as path from 'path'
import * as fs from 'fs'
import { fileURLToPath } from 'url'

// Get __dirname equivalent in ES modules
const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

test.describe('Screenshot Capture', () => {
  test('captures Web UI dashboard screenshot', async ({ page }) => {
    // Ensure assets directory exists
    // From tests/e2e/tests/, go up three levels to repo root, then to docs/assets
    const repoRoot = path.resolve(__dirname, '..', '..', '..')
    const assetsDir = path.join(repoRoot, 'docs', 'assets')
    if (!fs.existsSync(assetsDir)) {
      fs.mkdirSync(assetsDir, { recursive: true })
    }

    // Set auth header if AUTH_TOKEN is provided (for lab profile)
    const authToken = process.env.AUTH_TOKEN
    if (authToken) {
      await page.setExtraHTTPHeaders({
        'Authorization': `Bearer ${authToken}`
      })
    }

    // Navigate to dashboard and wait for page to load
    await page.goto('/', { waitUntil: 'networkidle' })

    // Set viewport to a nice size
    await page.setViewportSize({ width: 1280, height: 720 })

    // Wait for the page to be interactive - check for any content
    // Try to find the dashboard heading, but don't fail if it's not there
    // The page might show an error or loading state, which is still useful to screenshot
    try {
      await expect(page.getByRole('heading', { name: /dashboard/i })).toBeVisible({ timeout: 5000 })
    } catch {
      // If dashboard heading not found, wait for any content to appear
      try {
        await expect(page.locator('body')).toBeVisible({ timeout: 5000 })
        // Wait a bit for any dynamic content to load
        await page.waitForTimeout(2000)
      } catch {
        // Even if nothing specific loads, wait a moment and take screenshot anyway
        await page.waitForTimeout(2000)
      }
    }

    // Additional wait for any animations or final rendering
    await page.waitForTimeout(1000)

    // Capture screenshot
    const screenshotPath = path.join(assetsDir, 'dilivet-web-ui.png')
    await page.screenshot({
      path: screenshotPath,
      fullPage: true,
    })

    // Verify screenshot was created
    expect(fs.existsSync(screenshotPath)).toBe(true)

    // Log success
    console.log(`Screenshot saved to: ${screenshotPath}`)
  })
})


/**
 * DiliVet — ML-DSA diagnostics toolkit
 * Copyright (c) 2025 Thor Thor (codethor0)
 * Project: github.com/codethor0/dilivet
 * LinkedIn: https://www.linkedin.com/in/thor-thor0
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import KATVerify from './KATVerify'
import * as api from '../api/client'

vi.mock('../api/client')

describe('KATVerify', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders KAT verification page', () => {
    render(<KATVerify />)
    expect(screen.getByText('KAT Verification')).toBeInTheDocument()
  })

  it('displays run button', () => {
    render(<KATVerify />)
    expect(screen.getByRole('button', { name: /Run KAT Verification/i })).toBeInTheDocument()
  })

  it('calls verifyKAT when button is clicked', async () => {
    const mockVerifyKAT = vi.spyOn(api, 'verifyKAT')
    mockVerifyKAT.mockResolvedValue({
      ok: true,
      totalVectors: 100,
      passed: 95,
      failed: 5,
    })

    render(<KATVerify />)
    const user = userEvent.setup()

    await user.click(screen.getByRole('button', { name: /Run KAT Verification/i }))

    await waitFor(() => {
      expect(mockVerifyKAT).toHaveBeenCalled()
    })
  })

  it('displays results after successful verification', async () => {
    const mockVerifyKAT = vi.spyOn(api, 'verifyKAT')
    mockVerifyKAT.mockResolvedValue({
      ok: true,
      totalVectors: 100,
      passed: 95,
      failed: 5,
      decodeFailures: 0,
    })

    render(<KATVerify />)
    const user = userEvent.setup()

    await user.click(screen.getByRole('button', { name: /Run KAT Verification/i }))

    // Wait for API call to complete
    await waitFor(() => {
      expect(mockVerifyKAT).toHaveBeenCalled()
    })

    // Wait for result box to appear first (longer timeout for CI)
    await waitFor(
      () => {
        expect(screen.getByText(/KAT Verification Results/i)).toBeInTheDocument()
      },
      { timeout: 5000 }
    )

    // Then verify the specific values (longer timeout for CI)
    await waitFor(
      () => {
        expect(screen.getByText(/Total Vectors:/i)).toBeInTheDocument()
        // Use more specific queries that match the actual rendered text
        expect(screen.getByText('100')).toBeInTheDocument()
        expect(screen.getByText('95')).toBeInTheDocument()
        expect(screen.getByText('5')).toBeInTheDocument()
      },
      { timeout: 5000 }
    )
  })

  it('displays error when verification fails', async () => {
    const mockVerifyKAT = vi.spyOn(api, 'verifyKAT')
    mockVerifyKAT.mockResolvedValue({
      ok: false,
      error: 'Failed to load vectors',
    })

    render(<KATVerify />)
    const user = userEvent.setup()

    await user.click(screen.getByRole('button', { name: /Run KAT Verification/i }))

    await waitFor(() => {
      expect(screen.getByText(/Error/)).toBeInTheDocument()
      expect(screen.getByText(/Failed to load vectors/)).toBeInTheDocument()
    })
  })

  it('shows loading state while running', async () => {
    const mockVerifyKAT = vi.spyOn(api, 'verifyKAT')
    mockVerifyKAT.mockImplementation(
      () => new Promise((resolve) => {
        setTimeout(() => resolve({ ok: true, totalVectors: 100, passed: 100, failed: 0 }), 100)
      })
    )

    render(<KATVerify />)
    const user = userEvent.setup()

    await user.click(screen.getByRole('button', { name: /Run KAT Verification/i }))

    expect(screen.getByText(/Running KAT Verification/)).toBeInTheDocument()

    await waitFor(() => {
      expect(screen.queryByText(/Running KAT Verification/)).not.toBeInTheDocument()
    })
  })
})


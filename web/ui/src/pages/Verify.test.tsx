/**
 * DiliVet — ML-DSA diagnostics toolkit
 * Copyright (c) 2025 Thor Thor (codethor0)
 * Project: github.com/codethor0/dilivet
 * LinkedIn: https://www.linkedin.com/in/thor-thor0
 */

import { describe, it, expect, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import Verify from './Verify'
import * as api from '../api/client'

vi.mock('../api/client')

describe('Verify', () => {
  it('renders the verify form', () => {
    render(<Verify />)
    expect(screen.getByText('Verify Signature')).toBeInTheDocument()
    expect(screen.getByLabelText('Parameter Set')).toBeInTheDocument()
    expect(screen.getByLabelText('Public Key (hex)')).toBeInTheDocument()
  })

  it('submits verification request', async () => {
    const mockVerify = vi.spyOn(api, 'verifySignature')
    mockVerify.mockResolvedValue({
      ok: true,
      result: 'valid',
    })

    render(<Verify />)
    const user = userEvent.setup()

    await user.type(screen.getByLabelText('Public Key (hex)'), 'deadbeef')
    await user.type(screen.getByLabelText('Signature (hex)'), 'cafebabe')
    await user.type(screen.getByLabelText('Message (UTF-8)'), 'test message')
    await user.click(screen.getByRole('button', { name: /verify/i }))

    await waitFor(() => {
      expect(mockVerify).toHaveBeenCalled()
    })
  })

  it('displays error on verification failure', async () => {
    const mockVerify = vi.spyOn(api, 'verifySignature')
    mockVerify.mockResolvedValue({
      ok: false,
      error: 'Invalid signature format',
    })

    render(<Verify />)
    const user = userEvent.setup()

    await user.type(screen.getByLabelText('Public Key (hex)'), 'deadbeef')
    await user.type(screen.getByLabelText('Signature (hex)'), 'cafebabe')
    await user.type(screen.getByLabelText('Message (UTF-8)'), 'test')
    await user.click(screen.getByRole('button', { name: /verify/i }))

    await waitFor(() => {
      expect(screen.getByText(/error/i)).toBeInTheDocument()
    })
  })
})


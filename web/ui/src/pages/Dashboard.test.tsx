import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import Dashboard from './Dashboard'
import * as api from '../api/client'

vi.mock('../api/client')

describe('Dashboard', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders dashboard title', () => {
    vi.spyOn(api, 'getHealth').mockResolvedValue({
      status: 'ok',
      version: '0.2.4',
    })

    render(<Dashboard />)
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
  })

  it('calls getHealth on mount', async () => {
    const mockGetHealth = vi.spyOn(api, 'getHealth')
    mockGetHealth.mockResolvedValue({
      status: 'ok',
      version: '0.2.4',
    })

    render(<Dashboard />)

    await waitFor(() => {
      expect(mockGetHealth).toHaveBeenCalled()
    })
  })

  it('displays health status and version', async () => {
    vi.spyOn(api, 'getHealth').mockResolvedValue({
      status: 'ok',
      version: '0.2.4',
    })

    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.getByText(/Status:/)).toBeInTheDocument()
      expect(screen.getByText(/Version:/)).toBeInTheDocument()
      expect(screen.getByText('0.2.4')).toBeInTheDocument()
    })
  })

  it('displays error when health check fails', async () => {
    vi.spyOn(api, 'getHealth').mockRejectedValue(new Error('Network error'))

    render(<Dashboard />)

    await waitFor(() => {
      expect(screen.getByText(/Error:/)).toBeInTheDocument()
      expect(screen.getByText(/Network error/)).toBeInTheDocument()
    })
  })

  it('shows loading state initially', () => {
    vi.spyOn(api, 'getHealth').mockImplementation(
      () => new Promise(() => {}) // Never resolves
    )

    render(<Dashboard />)
    expect(screen.getByText(/Checking server status/)).toBeInTheDocument()
  })
})


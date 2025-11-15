import { describe, it, expect, vi, beforeEach } from 'vitest'
import { getHealth, verifySignature, verifyKAT, type VerifyRequest, type KATVerifyRequest } from './client'

// Mock global fetch
const mockFetch = vi.fn()
globalThis.fetch = mockFetch as any

describe('API Client', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    globalThis.fetch = mockFetch as any
  })

  describe('getHealth', () => {
    it('returns health status and version', async () => {
      const mockResponse = { status: 'ok', version: '0.2.4' }
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await getHealth()

      expect(result).toEqual(mockResponse)
      expect(mockFetch).toHaveBeenCalledWith('/api/health')
    })

    it('throws on network error', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'))

      await expect(getHealth()).rejects.toThrow('Network error')
    })

    it('throws on non-200 status', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
      })

      await expect(getHealth()).rejects.toThrow('Health check failed')
    })
  })

  describe('verifySignature', () => {
    it('sends correct request and returns valid result', async () => {
      const mockRequest: VerifyRequest = {
        paramSet: 'ML-DSA-44',
        publicKeyHex: 'deadbeef',
        signatureHex: 'cafebabe',
        message: 'test',
      }
      const mockResponse = { ok: true, result: 'valid' }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await verifySignature(mockRequest)

      expect(result).toEqual(mockResponse)
      expect(mockFetch).toHaveBeenCalledWith('/api/verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(mockRequest),
      })
    })

    it('handles error response', async () => {
      const mockRequest: VerifyRequest = {
        paramSet: 'ML-DSA-44',
        publicKeyHex: 'deadbeef',
        signatureHex: 'cafebabe',
        message: 'test',
      }
      const mockResponse = { ok: false, error: 'Invalid signature format' }

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        statusText: 'Bad Request',
        json: async () => mockResponse,
      })

      const result = await verifySignature(mockRequest)

      expect(result.ok).toBe(false)
      expect(result.error).toBe('Invalid signature format')
    })

    it('handles network error', async () => {
      const mockRequest: VerifyRequest = {
        paramSet: 'ML-DSA-44',
        publicKeyHex: 'deadbeef',
        signatureHex: 'cafebabe',
        message: 'test',
      }

      mockFetch.mockRejectedValueOnce(new Error('Network error'))

      await expect(verifySignature(mockRequest)).rejects.toThrow('Network error')
    })

    it('handles malformed JSON response', async () => {
      const mockRequest: VerifyRequest = {
        paramSet: 'ML-DSA-44',
        publicKeyHex: 'deadbeef',
        signatureHex: 'cafebabe',
        message: 'test',
      }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => {
          throw new Error('Invalid JSON')
        },
      })

      await expect(verifySignature(mockRequest)).rejects.toThrow()
    })
  })

  describe('verifyKAT', () => {
    it('sends correct request and returns results', async () => {
      const mockRequest: KATVerifyRequest = {}
      const mockResponse = {
        ok: true,
        totalVectors: 100,
        passed: 95,
        failed: 5,
      }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await verifyKAT(mockRequest)

      expect(result).toEqual(mockResponse)
      expect(mockFetch).toHaveBeenCalledWith('/api/kat-verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(mockRequest),
      })
    })

    it('handles error response', async () => {
      const mockRequest: KATVerifyRequest = { vectorsPath: '/invalid/path' }
      const mockResponse = { ok: false, error: 'Failed to load vectors' }

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
        json: async () => mockResponse,
      })

      const result = await verifyKAT(mockRequest)

      expect(result.ok).toBe(false)
      expect(result.error).toBe('Failed to load vectors')
    })
  })
})


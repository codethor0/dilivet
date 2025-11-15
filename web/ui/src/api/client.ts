/**
 * DiliVet — ML-DSA diagnostics toolkit
 * Copyright (c) 2025 Thor Thor (codethor0)
 * Project: github.com/codethor0/dilivet
 * LinkedIn: https://www.linkedin.com/in/thor-thor0
 */

// @ts-ignore - Vite injects import.meta.env at build time
const API_BASE = import.meta.env?.VITE_API_BASE || '/api'

// Get auth token from environment variable or localStorage (for lab/hardened profiles)
function getAuthToken(): string | null {
  // Check environment variable first (for server-side or build-time)
  if (typeof window !== 'undefined') {
    // Browser environment: check localStorage
    const stored = localStorage.getItem('dilivet_auth_token')
    if (stored) return stored
    
    // Check if token is in URL params (for easy setup)
    const params = new URLSearchParams(window.location.search)
    const urlToken = params.get('token')
    if (urlToken) {
      localStorage.setItem('dilivet_auth_token', urlToken)
      return urlToken
    }
  }
  return null
}

// Helper to add auth headers if token is available
function getHeaders(includeAuth: boolean = true): HeadersInit {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  }
  if (includeAuth) {
    const token = getAuthToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
  }
  return headers
}

export interface HealthResponse {
  status: string
  version: string
}

export interface VerifyRequest {
  paramSet: string
  publicKeyHex: string
  signatureHex: string
  messageHex?: string
  message?: string
}

export interface VerifyResponse {
  ok: boolean
  result?: 'valid' | 'invalid'
  error?: string
}

export interface KATVerifyRequest {
  vectorsPath?: string
}

export interface KATVerifyResponse {
  ok: boolean
  totalVectors?: number
  passed?: number
  failed?: number
  decodeFailures?: number
  error?: string
  details?: KATVerifyDetail[]
}

export interface KATVerifyDetail {
  caseId: number
  passed: boolean
  parameterSet?: string
  reason?: string
}

export async function getHealth(): Promise<HealthResponse> {
  const res = await fetch(`${API_BASE}/health`, {
    headers: getHeaders(true),
  })
  if (!res.ok) {
    // If 401, try to get token from URL or show helpful error
    if (res.status === 401) {
      const params = new URLSearchParams(window.location.search)
      if (!params.get('token') && !localStorage.getItem('dilivet_auth_token')) {
        throw new Error('Authentication required. Add ?token=YOUR_TOKEN to the URL or set dilivet_auth_token in localStorage.')
      }
    }
    throw new Error(`Health check failed: ${res.statusText}`)
  }
  return res.json()
}

export async function verifySignature(req: VerifyRequest): Promise<VerifyResponse> {
  const res = await fetch(`${API_BASE}/verify`, {
    method: 'POST',
    headers: getHeaders(true),
    body: JSON.stringify(req),
  })

  const data = await res.json()
  if (!res.ok) {
    return {
      ok: false,
      error: data.error || `Request failed: ${res.statusText}`,
    }
  }
  return data
}

export async function verifyKAT(req: KATVerifyRequest = {}): Promise<KATVerifyResponse> {
  const res = await fetch(`${API_BASE}/kat-verify`, {
    method: 'POST',
    headers: getHeaders(true),
    body: JSON.stringify(req),
  })

  const data = await res.json()
  if (!res.ok) {
    return {
      ok: false,
      error: data.error || `Request failed: ${res.statusText}`,
    }
  }
  return data
}


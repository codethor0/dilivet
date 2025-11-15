// @ts-ignore - Vite injects import.meta.env at build time
const API_BASE = import.meta.env?.VITE_API_BASE || '/api'

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
  const res = await fetch(`${API_BASE}/health`)
  if (!res.ok) {
    throw new Error(`Health check failed: ${res.statusText}`)
  }
  return res.json()
}

export async function verifySignature(req: VerifyRequest): Promise<VerifyResponse> {
  const res = await fetch(`${API_BASE}/verify`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
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
    headers: {
      'Content-Type': 'application/json',
    },
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


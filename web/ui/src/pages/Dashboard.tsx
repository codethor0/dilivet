/**
 * DiliVet — ML-DSA diagnostics toolkit
 * Copyright (c) 2025 Thor Thor (codethor0)
 * Project: github.com/codethor0/dilivet
 * LinkedIn: https://www.linkedin.com/in/thor-thor0
 */

import { useEffect, useState } from 'react'
import { getHealth, HealthResponse, setAuthToken, getAuthToken } from '../api/client'
import './Dashboard.css'

function Dashboard() {
  const [health, setHealth] = useState<HealthResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showTokenInput, setShowTokenInput] = useState(false)
  const [tokenValue, setTokenValue] = useState('')
  const [authError, setAuthError] = useState(false)

  const checkHealth = () => {
    setLoading(true)
    setError(null)
    setAuthError(false)
    getHealth()
      .then((data) => {
        setHealth(data)
        setAuthError(false)
      })
      .catch((err) => {
        const isAuthError = err.message.includes('Authentication') || 
                           err.message.includes('Authorization') ||
                           err.message.includes('401') ||
                           err.message.includes('403')
        if (isAuthError) {
          setAuthError(true)
          setShowTokenInput(true)
          setError('Authentication required. Please enter your auth token.')
        } else {
          setError(err.message)
        }
      })
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    checkHealth()
  }, [])

  const handleTokenSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (tokenValue.trim()) {
      setAuthToken(tokenValue.trim())
      setTokenValue('')
      setShowTokenInput(false)
      // Retry health check with new token
      checkHealth()
    }
  }

  const handleTokenClear = () => {
    localStorage.removeItem('dilivet_auth_token')
    setTokenValue('')
    setShowTokenInput(true)
    checkHealth()
  }

  return (
    <div className="dashboard">
      <h2>Dashboard</h2>
      <div className="dashboard-content">
        <section className="intro">
          <h3>About DiliVet</h3>
          <p>
            DiliVet is a diagnostics and vetting toolkit for ML-DSA (Dilithium-like)
            signature implementations. This web interface provides access to core
            diagnostic features including signature verification and known-answer
            test (KAT) verification.
          </p>
          <p>
            <strong>Note:</strong> This is diagnostics tooling, not a production
            cryptographic library. Use in controlled environments only.
          </p>
        </section>

        <section className="health-status">
          <h3>Server Status</h3>
          {loading && <p>Checking server status...</p>}
          
          {authError && showTokenInput && (
            <div className="auth-box">
              <h4>Authentication Required</h4>
              <p>This server requires an authentication token. Please enter your token below:</p>
              <form onSubmit={handleTokenSubmit} className="token-form">
                <div className="token-input-group">
                  <input
                    type="text"
                    value={tokenValue}
                    onChange={(e) => setTokenValue(e.target.value)}
                    placeholder="Enter your auth token"
                    className="token-input"
                    autoFocus
                  />
                  <button type="submit" className="token-submit-btn">
                    Set Token
                  </button>
                </div>
                {getAuthToken() && (
                  <button
                    type="button"
                    onClick={handleTokenClear}
                    className="token-clear-btn"
                  >
                    Clear Saved Token
                  </button>
                )}
              </form>
              <p className="token-hint">
                <small>
                  Tip: You can also add <code>?token=YOUR_TOKEN</code> to the URL, 
                  or set it via browser console: <code>localStorage.setItem('dilivet_auth_token', 'YOUR_TOKEN')</code>
                </small>
              </p>
            </div>
          )}
          
          {error && !authError && (
            <div className="error-box">
              <strong>Error:</strong> {error}
            </div>
          )}
          
          {health && (
            <div className="health-info">
              <p>
                <strong>Status:</strong> <span className="status-ok">{health.status}</span>
              </p>
              <p>
                <strong>Version:</strong> {health.version}
              </p>
              {getAuthToken() && (
                <p className="auth-status">
                  <strong>Auth:</strong> <span className="status-ok">Token configured</span>
                  <button onClick={handleTokenClear} className="token-clear-link">
                    (clear)
                  </button>
                </p>
              )}
            </div>
          )}
        </section>

        <section className="features">
          <h3>Available Features</h3>
          <ul>
            <li>
              <strong>Verify Signature:</strong> Validate ML-DSA signatures against
              public keys for parameter sets ML-DSA-44, ML-DSA-65, and ML-DSA-87.
            </li>
            <li>
              <strong>KAT Verification:</strong> Run known-answer test vectors
              through structural checks to validate implementation correctness.
            </li>
          </ul>
        </section>
      </div>
    </div>
  )
}

export default Dashboard


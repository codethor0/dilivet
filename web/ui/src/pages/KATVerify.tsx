/**
 * DiliVet — ML-DSA diagnostics toolkit
 * Copyright (c) 2025 Thor Thor (codethor0)
 * Project: github.com/codethor0/dilivet
 * LinkedIn: https://www.linkedin.com/in/thor-thor0
 */

import { useState } from 'react'
import { verifyKAT, KATVerifyResponse } from '../api/client'
import './KATVerify.css'

function KATVerify() {
  const [result, setResult] = useState<KATVerifyResponse | null>(null)
  const [loading, setLoading] = useState(false)

  const handleRun = async () => {
    setLoading(true)
    setResult(null)

    try {
      const res = await verifyKAT()
      setResult(res)
    } catch (err) {
      setResult({
        ok: false,
        error: err instanceof Error ? err.message : 'Unknown error',
      })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="kat-verify">
      <h2>KAT Verification</h2>
      <div className="kat-info">
        <p>
          Run known-answer test (KAT) verification using ACVP sigVer vectors.
          This will validate the implementation against standardized test vectors.
        </p>
      </div>

      <button
        onClick={handleRun}
        disabled={loading}
        className="run-button"
      >
        {loading ? 'Running KAT Verification...' : 'Run KAT Verification'}
      </button>

      {result && (
        <div className={`result-box ${result.ok ? 'result-success' : 'result-error'}`}>
          {result.ok ? (
            <>
              <h3>KAT Verification Results</h3>
              <div className="results-grid">
                <div className="result-item">
                  <strong>Total Vectors:</strong> {result.totalVectors}
                </div>
                <div className="result-item">
                  <strong>Passed:</strong>{' '}
                  <span className="passed">{result.passed}</span>
                </div>
                <div className="result-item">
                  <strong>Failed:</strong>{' '}
                  <span className="failed">{result.failed}</span>
                </div>
                {result.decodeFailures !== undefined && result.decodeFailures > 0 && (
                  <div className="result-item">
                    <strong>Decode Failures:</strong>{' '}
                    <span className="failed">{result.decodeFailures}</span>
                  </div>
                )}
              </div>
              {result.details && result.details.length > 0 && (
                <div className="details-section">
                  <h4>Test Details</h4>
                  <div className="details-table">
                    <table>
                      <thead>
                        <tr>
                          <th>Case ID</th>
                          <th>Parameter Set</th>
                          <th>Status</th>
                          <th>Reason</th>
                        </tr>
                      </thead>
                      <tbody>
                        {result.details.slice(0, 50).map((detail, idx) => (
                          <tr key={idx}>
                            <td>{detail.caseId}</td>
                            <td>{detail.parameterSet || '-'}</td>
                            <td>
                              <span className={detail.passed ? 'passed' : 'failed'}>
                                {detail.passed ? 'Passed' : 'Failed'}
                              </span>
                            </td>
                            <td>{detail.reason || '-'}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                    {result.details.length > 50 && (
                      <p className="details-note">
                        Showing first 50 of {result.details.length} test cases
                      </p>
                    )}
                  </div>
                </div>
              )}
            </>
          ) : (
            <>
              <h3>Error</h3>
              <p>{result.error}</p>
            </>
          )}
        </div>
      )}
    </div>
  )
}

export default KATVerify


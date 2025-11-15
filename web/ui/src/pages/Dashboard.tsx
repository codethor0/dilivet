import { useEffect, useState } from 'react'
import { getHealth, HealthResponse } from '../api/client'
import './Dashboard.css'

function Dashboard() {
  const [health, setHealth] = useState<HealthResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    getHealth()
      .then(setHealth)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

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
          {error && (
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


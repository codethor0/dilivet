import { useState } from 'react'
import { verifySignature, VerifyResponse } from '../api/client'
import './Verify.css'

function Verify() {
  const [paramSet, setParamSet] = useState('ML-DSA-44')
  const [publicKeyHex, setPublicKeyHex] = useState('')
  const [signatureHex, setSignatureHex] = useState('')
  const [messageMode, setMessageMode] = useState<'hex' | 'text'>('text')
  const [messageHex, setMessageHex] = useState('')
  const [message, setMessage] = useState('')
  const [result, setResult] = useState<VerifyResponse | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setResult(null)

    try {
      const req = {
        paramSet,
        publicKeyHex: publicKeyHex.trim(),
        signatureHex: signatureHex.trim(),
        ...(messageMode === 'hex'
          ? { messageHex: messageHex.trim() }
          : { message: message }),
      }

      const res = await verifySignature(req)
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
    <div className="verify">
      <h2>Verify Signature</h2>
      <form onSubmit={handleSubmit} className="verify-form">
        <div className="form-group">
          <label htmlFor="paramSet">Parameter Set</label>
          <select
            id="paramSet"
            value={paramSet}
            onChange={(e) => setParamSet(e.target.value)}
            required
          >
            <option value="ML-DSA-44">ML-DSA-44</option>
            <option value="ML-DSA-65">ML-DSA-65</option>
            <option value="ML-DSA-87">ML-DSA-87</option>
          </select>
        </div>

        <div className="form-group">
          <label htmlFor="publicKeyHex">Public Key (hex)</label>
          <textarea
            id="publicKeyHex"
            value={publicKeyHex}
            onChange={(e) => setPublicKeyHex(e.target.value)}
            placeholder="Enter hex-encoded public key"
            rows={3}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="signatureHex">Signature (hex)</label>
          <textarea
            id="signatureHex"
            value={signatureHex}
            onChange={(e) => setSignatureHex(e.target.value)}
            placeholder="Enter hex-encoded signature"
            rows={3}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="messageMode">Message Format</label>
          <div className="radio-group">
            <label>
              <input
                type="radio"
                value="text"
                checked={messageMode === 'text'}
                onChange={() => setMessageMode('text')}
              />
              UTF-8 Text
            </label>
            <label>
              <input
                type="radio"
                value="hex"
                checked={messageMode === 'hex'}
                onChange={() => setMessageMode('hex')}
              />
              Hex
            </label>
          </div>
        </div>

        {messageMode === 'hex' ? (
          <div className="form-group">
            <label htmlFor="messageHex">Message (hex)</label>
            <textarea
              id="messageHex"
              value={messageHex}
              onChange={(e) => setMessageHex(e.target.value)}
              placeholder="Enter hex-encoded message"
              rows={3}
              required
            />
          </div>
        ) : (
          <div className="form-group">
            <label htmlFor="message">Message (UTF-8)</label>
            <textarea
              id="message"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="Enter message text"
              rows={3}
              required
            />
          </div>
        )}

        <button type="submit" disabled={loading} className="submit-button">
          {loading ? 'Verifying...' : 'Verify'}
        </button>
      </form>

      {result && (
        <div className={`result-box ${result.ok ? 'result-success' : 'result-error'}`}>
          {result.ok ? (
            <>
              <h3>Verification Result</h3>
              <p>
                <strong>Status:</strong>{' '}
                <span className={result.result === 'valid' ? 'valid' : 'invalid'}>
                  {result.result === 'valid' ? 'Valid signature' : 'Invalid signature'}
                </span>
              </p>
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

export default Verify


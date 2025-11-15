// DiliVet Web - Load test for /api/verify endpoint
// Run with: k6 run verify_load.js

import http from 'k6/http'
import { check, sleep } from 'k6'

export const options = {
  stages: [
    { duration: '10s', target: 10 },  // Ramp up to 10 users
    { duration: '30s', target: 50 }, // Ramp up to 50 users
    { duration: '20s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% of requests should be below 2s
    http_req_failed: ['rate<0.1'],     // Error rate should be less than 10%
  },
}

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080'

export default function () {
  // Test payload (will fail verification but tests API structure)
  const payload = JSON.stringify({
    paramSet: 'ML-DSA-44',
    publicKeyHex: 'deadbeefdeadbeef',
    signatureHex: 'cafebabecafebabe',
    message: 'test message',
  })

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  }

  const res = http.post(`${BASE_URL}/api/verify`, payload, params)

  check(res, {
    'status is 200 or 400': (r) => r.status === 200 || r.status === 400,
    'response has ok field': (r) => {
      try {
        const body = JSON.parse(r.body)
        return 'ok' in body
      } catch {
        return false
      }
    },
    'response time < 2s': (r) => r.timings.duration < 2000,
  })

  sleep(1)
}


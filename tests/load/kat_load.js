// DiliVet Web - Load test for /api/kat-verify endpoint
// Run with: k6 run kat_load.js

import http from 'k6/http'
import { check, sleep } from 'k6'

export const options = {
  stages: [
    { duration: '10s', target: 2 },  // Ramp up to 2 users (KAT is heavy)
    { duration: '60s', target: 5 },  // Ramp up to 5 users
    { duration: '20s', target: 0 },  // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<30000'], // 95% of requests should be below 30s (KAT takes time)
    http_req_failed: ['rate<0.1'],      // Error rate should be less than 10%
  },
}

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080'

export default function () {
  const payload = JSON.stringify({})

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  }

  const res = http.post(`${BASE_URL}/api/kat-verify`, payload, params)

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response has ok field': (r) => {
      try {
        const body = JSON.parse(r.body)
        return 'ok' in body
      } catch {
        return false
      }
    },
    'response has totalVectors': (r) => {
      try {
        const body = JSON.parse(r.body)
        return 'totalVectors' in body
      } catch {
        return false
      }
    },
  })

  sleep(5) // Longer sleep for KAT (it's CPU intensive)
}


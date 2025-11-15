// DiliVet Web - Load test for /api/health endpoint
// Run with: k6 run health_load.js

import http from 'k6/http'
import { check, sleep } from 'k6'

export const options = {
  stages: [
    { duration: '10s', target: 50 },  // Ramp up to 50 users
    { duration: '30s', target: 100 },  // Ramp up to 100 users
    { duration: '20s', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<100'],  // 95% of requests should be below 100ms
    http_req_failed: ['rate<0.01'],     // Error rate should be less than 1%
  },
}

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080'

export default function () {
  const res = http.get(`${BASE_URL}/api/health`)

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response has status field': (r) => {
      try {
        const body = JSON.parse(r.body)
        return body.status === 'ok'
      } catch {
        return false
      }
    },
    'response time < 100ms': (r) => r.timings.duration < 100,
  })

  sleep(0.1)
}


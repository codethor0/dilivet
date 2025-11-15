<!--
DiliVet – ML-DSA diagnostics and vetting toolkit
Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)
-->

## Wycheproof-style plan for ML-DSA validation

### Goals

1. Curate adversarial cases that complement the official ACVP vectors.
2. Provide an open contribution model so implementers can reproduce issues.
3. Keep the format simple: text-based JSON/req files with minimal boilerplate.

### Categories to curate

| Category | Description | Example sources |
|----------|-------------|-----------------|
| Encoding traps | DER/CMS-like length/offset quirks, leading-zero trimming mistakes | Alternative encodings observed in other PQC libraries |
| Range violations | Coefficients outside expected bounds, forged hints | Mutations produced by fuzzers or failed constant-time primitives |
| Degenerate keys | All-zero or repeated coefficients, invalid public matrix | Hands-on experimentation with `code/poly` helpers |
| Boundary sizes | Maximum-length messages, multi-precision overflow edges | Synthetic data from `code/kat.EdgeMsgs` |
| Replay / reuse | Re-used nonce data or per-message randomness collisions | Deterministic signers that forget to reseed |
| Serialization errors | Unexpected whitespace, extra keys, wrong casing | Interoperability with scripting languages and CLI pipelines |

Every category should include:

- A human-readable description of the failure mode.
- The minimal data needed to reproduce (message/public key/secret key).
- Expected result (accept, reject, or raise specific error code).

### JSON case format

```json
{
  "suite": "mldsa-wycheproof-alpha1",
  "cases": [
    {
      "id": "range-out-of-bounds-001",
      "description": "Coefficient q+1 in t1 causes rejection.",
      "msg": "<hex>",
      "pk": "<hex>",
      "sk": "<optional hex>",
      "sig": "<optional hex>",
      "expected": "reject"
    }
  ]
}
```

Guidelines:

- Keep keys lower-case, match the minimal set used by `code/kat`.
- Optional fields (`sk`, `sig`) may be omitted; deterministic signers can fill
  them in during verification.
- `expected` values: `accept`, `reject`, `error:<code>`.

### Contribution checklist

1. **Reproducer** – include a snippet (`go test`, `dilivet exec`, etc.) that
   demonstrates the failure with a current or downstream build.
2. **Minimality** – prove the case cannot be reduced further without losing
   the failure (e.g., use `go test -run TestMinimization`).
3. **Label** – choose a clear `id` and `description`; reference spec clauses
   when relevant.
4. **Metadata** – record the environment (toolchain, OS, commit SHA) in PR
   notes for historical context.
5. **Regression test** – where possible, add a unit test using `code/kat.Verify`
   to ensure we do not regress.
6. **License** – contributors agree to MIT/CC0-style data sharing for test
   vectors so they can be embedded downstream.

### Review process

- Triagers reproduce the issue, confirm it fits an existing category (or
  propose a new one), and ensure the data passes linting (`go test ./code/kat`).
- Merged cases update the public catalogue under `testdata/wycheproof`.
- Release tags should reference the catalogue version so downstream projects
  can pin to a known set of adversarial vectors.


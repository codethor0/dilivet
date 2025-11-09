## Fuzzing strategy

### Why we fuzz

Fuzzing gives us high-signal coverage across the parsing and verification paths
that routinely ingest untrusted data: KAT fixtures, interoperability payloads,
and edge-case messages. The harnesses in `code/kat` and the CLI adapters are
hardening layers designed to stay online even when fed malformed or malicious
inputs. Fuzzing helps surface panics, allocator blowups, and logic branches that
escape handwritten unit tests.

### Targets

- `FuzzDecodePublicKey` — probes the request/response parser used by the CLI
  and tooling. The target ensures we do not panic or leak resources when the
  KAT format is malformed.
- `FuzzVerify` — exercises the deterministic verification scaffolding while
  mutating messages, keys, and signatures. The focus is on resilience: the hook
  must reject unexpected data without triggering undefined behaviour.

Feel free to add more targets for higher-level adapters (e.g., future
net/http endpoints) — keep them under `./fuzz` so we can test them in one pass.

### Running locally

PR authors should run a short fuzz sweep before opening a review:

```bash
go test -fuzz=Fuzz -fuzztime=30s ./fuzz
```

This finds obvious panics without consuming an afternoon of CPU time.
Long-running fuzz campaigns should land in nightly automation (see below).

### CI integration

- `ci.yml` runs unit tests, vet, lint, and uploads coverage for every push/PR.
- `fuzz.yml` runs a one-minute fuzz sweep on each PR _and_ on a nightly cron.

Nightly fuzzing catches regressions that sneak in after the short PR run. If a
nightly fuzz job fails, triage it like a bug: capture the minimized corpus,
land a fix, and add a regression test so we avoid repeating the same mistake.

### Corpus hygiene

The native Go fuzzer persists crashing inputs in the `testdata/fuzz` tree.
Please keep the corpus small and informative — remove obsolete inputs once a
regression turns into a deterministic unit test.


# LinkedIn Post for DiliVet v0.2.4

DiliVet v0.2.4 is live.

This release focuses on correctness and robustness for ML-DSA (Dilithium-like) diagnostics:

• Fixed the useHint() implementation to fully match FIPS 204 Algorithm 3
• Added property-based tests and fuzzing around the hint logic and packing routines
• Benchmarked matrix A expansion for all parameter sets
• Introduced stress scripts for large inputs and weird encodings
• Updated the audit document with full coverage and results

All tests, race checks, fuzz runs, and stress tests pass. DiliVet is now production-ready as a diagnostics and vetting toolkit for ML-DSA signature implementations.

Repo: https://github.com/codethor0/dilivet  
Latest release: https://github.com/codethor0/dilivet/releases/tag/v0.2.4


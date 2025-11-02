# BUG-001: Double Highbits Computation

## Description
The `SignBug` function incorrectly computes `highBits` twice during the signing process, leading to invalid signatures that fail verification.

## Root Cause
In the signing algorithm, after computing `w0 = w - A·c`, the code should compute hints using `w0`. However, due to the bug, it computes `highBits` on the original `w` vector again instead of using the correct `w0` values.

## Impact
- Signatures created with `SignBug` will fail verification by standard ML-DSA implementations
- The bug is subtle and may not be immediately obvious during code review
- Affects the correctness of the post-quantum signature scheme

## Detection
The bug can be detected by:
1. Signing a message with the buggy implementation
2. Attempting to verify with a correct implementation
3. The verification will fail, indicating the presence of the bug

## Fix
Replace the second `highBits` computation with proper handling of `w0` values in the hint computation.

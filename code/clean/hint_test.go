// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package mldsa

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/codethor0/dilivet/code/poly"
)

// decompose splits a coefficient r into (r0, r1) where r = r1 * 2*gamma2 + r0 and |r0| <= gamma2.
// This is the mathematical decomposition used in FIPS 204.
func decompose(r uint32, gamma2 int) (r0, r1 int32) {
	canon := poly.Canonical(r)
	r1 = canon / int32(2*gamma2)
	if canon < 0 {
		r1--
	}
	// Adjust r0 to be in range [-gamma2, gamma2]
	r0 = canon - r1*int32(2*gamma2)
	return r0, r1
}

// makeHintForCoeff determines if a coefficient at index i needs a hint.
// Returns true if |r0| > gamma2 (needs adjustment).
func makeHintForCoeff(r0 int32, gamma2 int) bool {
	return r0 > int32(gamma2) || r0 < -int32(gamma2)
}

// TestUseHint_Property1_ReconstructionCorrectness tests that useHint correctly reconstructs w.
// Property: For all valid w, useHint(w, hasHint) produces w1' such that:
//   - r = w1' * 2*gamma2 + r0' where |r0'| <= gamma2
//   - If hasHint is false, w1' should equal the standard decomposition w1
func TestUseHint_Property1_ReconstructionCorrectness(t *testing.T) {
	gamma2Values := []int{
		ParamsMLDSA44.Gamma2,
		ParamsMLDSA65.Gamma2,
		ParamsMLDSA87.Gamma2,
	}

	for _, gamma2 := range gamma2Values {
		t.Run(fmt.Sprintf("gamma2=%d", gamma2), func(t *testing.T) {
			// Test specific edge cases
			testCases := []struct {
				name string
				r    uint32
			}{
				{"zero", 0},
				{"q-1", poly.Q - 1},
				{"q/2", poly.Q / 2},
				{"gamma2", uint32(gamma2)},
				{"2*gamma2", uint32(2 * gamma2)},
				{"q-gamma2", poly.Q - uint32(gamma2)},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					canon := poly.Canonical(tc.r)
					r0, w1 := decompose(tc.r, gamma2)
					hasHint := makeHintForCoeff(r0, gamma2)
					w1PrimeUint := useHintForCoeff(tc.r, hasHint, gamma2)

					// Convert w1' back to signed form for reconstruction
					w1PrimeSigned := int32(w1PrimeUint)
					if w1PrimeSigned > int32(poly.Q)/2 {
						w1PrimeSigned -= int32(poly.Q)
					}

					// Reconstruct r0' = r - w1' * 2*gamma2
					r0Prime := canon - w1PrimeSigned*int32(2*gamma2)

					// Verify |r0'| <= gamma2
					if r0Prime > int32(gamma2) || r0Prime < -int32(gamma2) {
						t.Errorf("useHint reconstruction invalid: r=%d w1'=%d r0'=%d (|r0'|=%d > gamma2=%d) hasHint=%v",
							tc.r, w1PrimeUint, r0Prime, abs(r0Prime), gamma2, hasHint)
					}

					// Verify reconstruction: r = w1' * 2*gamma2 + r0' (mod q)
					reconstructed := (w1PrimeSigned*int32(2*gamma2) + r0Prime) % int32(poly.Q)
					if reconstructed < 0 {
						reconstructed += int32(poly.Q)
					}
					expected := int32(poly.ModQ(tc.r))
					if expected != reconstructed {
						t.Errorf("useHint reconstruction mismatch: r=%d expected=%d reconstructed=%d hasHint=%v",
							tc.r, expected, reconstructed, hasHint)
					}

					// If no hint, w1' should equal standard decomposition w1
					if !hasHint {
						w1Uint := uint32(w1)
						if w1 < 0 {
							w1Uint = uint32(w1 + int32(poly.Q))
						}
						w1Uint %= poly.Q
						if w1Uint != w1PrimeUint {
							t.Errorf("useHint without hint should equal decompose: r=%d w1=%d(%d) w1'=%d",
								tc.r, w1, w1Uint, w1PrimeUint)
						}
					}
				})
			}

			// Property test with random values
			rng := rand.New(rand.NewSource(42))
			for i := 0; i < 1000; i++ {
				r := uint32(rng.Intn(int(poly.Q)))
				canon := poly.Canonical(r)
				r0, w1 := decompose(r, gamma2)
				hasHint := makeHintForCoeff(r0, gamma2)
				w1PrimeUint := useHintForCoeff(r, hasHint, gamma2)

				// Convert w1' back to signed form
				w1PrimeSigned := int32(w1PrimeUint)
				if w1PrimeSigned > int32(poly.Q)/2 {
					w1PrimeSigned -= int32(poly.Q)
				}

				// Reconstruct r0'
				r0Prime := canon - w1PrimeSigned*int32(2*gamma2)

				// Verify |r0'| <= gamma2
				if r0Prime > int32(gamma2) || r0Prime < -int32(gamma2) {
					t.Errorf("useHint reconstruction invalid: r=%d w1'=%d r0'=%d (|r0'|=%d > gamma2=%d) hasHint=%v",
						r, w1PrimeUint, r0Prime, abs(r0Prime), gamma2, hasHint)
				}

				// If no hint, w1' should equal w1
				if !hasHint {
					w1Uint := uint32(w1)
					if w1 < 0 {
						w1Uint = uint32(w1 + int32(poly.Q))
					}
					w1Uint %= poly.Q
					if w1Uint != w1PrimeUint {
						t.Errorf("useHint without hint should equal decompose: r=%d w1=%d(%d) w1'=%d",
							r, w1, w1Uint, w1PrimeUint)
					}
				}
			}
		})
	}
}

// abs returns the absolute value of an int32.
func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// TestUseHint_Property2_NoOpWhenNoHints tests that hints are no-ops when there are none.
// Property: For all w where |r0| <= gamma2, useHint(w, false) == decompose(w).w1
func TestUseHint_Property2_NoOpWhenNoHints(t *testing.T) {
	gamma2Values := []int{
		ParamsMLDSA44.Gamma2,
		ParamsMLDSA65.Gamma2,
		ParamsMLDSA87.Gamma2,
	}

	for _, gamma2 := range gamma2Values {
		t.Run(fmt.Sprintf("gamma2=%d", gamma2), func(t *testing.T) {
			rng := rand.New(rand.NewSource(42))
			for i := 0; i < 1000; i++ {
				r := uint32(rng.Intn(int(poly.Q)))
				r0, w1 := decompose(r, gamma2)

				// Only test cases where hint would be false (r0 in range)
				if r0 > int32(gamma2) || r0 < -int32(gamma2) {
					continue
				}

				w1PrimeUint := useHintForCoeff(r, false, gamma2)

				// Convert w1 to uint32 for comparison
				w1Uint := uint32(w1)
				if w1 < 0 {
					w1Uint = uint32(w1 + int32(poly.Q))
				}
				w1Uint %= poly.Q

				if w1Uint != w1PrimeUint {
					t.Errorf("useHint with hasHint=false should equal decompose: r=%d r0=%d w1=%d(%d) w1'=%d",
						r, r0, w1, w1Uint, w1PrimeUint)
				}
			}
		})
	}
}

// TestUseHint_Property3_Bounds tests that useHint returns values in the allowed range.
func TestUseHint_Property3_Bounds(t *testing.T) {
	gamma2Values := []int{
		ParamsMLDSA44.Gamma2,
		ParamsMLDSA65.Gamma2,
		ParamsMLDSA87.Gamma2,
	}

	for _, gamma2 := range gamma2Values {
		t.Run(fmt.Sprintf("gamma2=%d", gamma2), func(t *testing.T) {
			rng := rand.New(rand.NewSource(42))
			for i := 0; i < 1000; i++ {
				r := uint32(rng.Intn(int(poly.Q)))
				hasHint := rng.Intn(2) == 1 // Random hint

				w1Prime := useHintForCoeff(r, hasHint, gamma2)

				// w1' should be in valid range [0, q-1]
				if w1Prime >= poly.Q {
					t.Errorf("useHint out of bounds: r=%d hasHint=%v w1'=%d (>= %d)",
						r, hasHint, w1Prime, poly.Q)
				}
			}
		})
	}
}

// normalizeToModQ reduces x to the valid coefficient range [0, q-1].
func normalizeToModQ(x int64) uint32 {
	result := int64(x) % int64(poly.Q)
	if result < 0 {
		result += int64(poly.Q)
	}
	return uint32(result)
}

// FuzzMakeUseHintRoundTrip fuzzes the relationship between decompose, makeHint, and useHint.
// This ensures that useHint correctly reconstructs w1' for all valid inputs.
func FuzzMakeUseHintRoundTrip(f *testing.F) {
	gamma2Values := []int{
		ParamsMLDSA44.Gamma2,
		ParamsMLDSA65.Gamma2,
		ParamsMLDSA87.Gamma2,
	}

	// Seed with edge values
	for _, gamma2 := range gamma2Values {
		f.Add(int64(0), gamma2)
		f.Add(int64(poly.Q-1), gamma2)
		f.Add(int64(poly.Q/2), gamma2)
		f.Add(int64(gamma2), gamma2)
		f.Add(int64(2*gamma2), gamma2)
		f.Add(int64(poly.Q-gamma2), gamma2)
	}

	f.Fuzz(func(t *testing.T, x int64, gamma2 int) {
		// Validate gamma2
		validGamma2 := false
		for _, g := range gamma2Values {
			if gamma2 == g {
				validGamma2 = true
				break
			}
		}
		if !validGamma2 {
			t.Skip("invalid gamma2")
		}

		// Normalize x to valid coefficient range
		r := normalizeToModQ(x)
		canon := poly.Canonical(r)

		// Decompose
		r0, _ := decompose(r, gamma2)
		hasHint := makeHintForCoeff(r0, gamma2)

		// Use hint to reconstruct
		w1PrimeUint := useHintForCoeff(r, hasHint, gamma2)

		// Convert w1' back to signed form
		w1PrimeSigned := int32(w1PrimeUint)
		if w1PrimeSigned > int32(poly.Q)/2 {
			w1PrimeSigned -= int32(poly.Q)
		}

		// Reconstruct r0'
		r0Prime := canon - w1PrimeSigned*int32(2*gamma2)

		// Verify |r0'| <= gamma2
		if r0Prime > int32(gamma2) || r0Prime < -int32(gamma2) {
			t.Fatalf("useHint reconstruction invalid: r=%d w1'=%d r0'=%d (|r0'|=%d > gamma2=%d) hasHint=%v",
				r, w1PrimeUint, r0Prime, abs(r0Prime), gamma2, hasHint)
		}

		// Verify reconstruction: r = w1' * 2*gamma2 + r0' (mod q)
		reconstructed := (w1PrimeSigned*int32(2*gamma2) + r0Prime) % int32(poly.Q)
		if reconstructed < 0 {
			reconstructed += int32(poly.Q)
		}
		expected := int32(poly.ModQ(r))
		if expected != reconstructed {
			t.Fatalf("useHint reconstruction mismatch: r=%d expected=%d reconstructed=%d hasHint=%v",
				r, expected, reconstructed, hasHint)
		}
	})
}

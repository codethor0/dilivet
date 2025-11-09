// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package poly

import "errors"

const (
	// Degree of ML-DSA polynomials.
	N = 256

	// Q is the prime modulus used by ML-DSA.
	Q = 8380417

	// qInv = -(q^{-1}) mod 2^32, used in Montgomery reduction.
	qInv = 4236238847

	// Mont is R mod q with R = 2^32.
	Mont = 4193792
)

// Poly represents an ML-DSA polynomial with coefficients modulo q.
type Poly struct {
	Coeffs [N]uint32
}

var errNotImplemented = errors.New("poly: not implemented")

// ReduceLe2Q reduces x into [0, 2q).
func ReduceLe2Q(x uint32) uint32 {
	x1 := x >> 23
	x2 := x & 0x7fffff
	return x2 + (x1 << 13) - x1
}

// Le2QModQ reduces x from [0, 2q) into [0, q).
func Le2QModQ(x uint32) uint32 {
	x -= Q
	mask := uint32(int32(x) >> 31)
	return x + (mask & Q)
}

// ModQ returns x mod q for any uint32.
func ModQ(x uint32) uint32 {
	return Le2QModQ(ReduceLe2Q(x))
}

// montReduceLe2Q computes a * R^{-1} mod q and returns a value in [0, 2q).
func montReduceLe2Q(a uint64) uint32 {
	m := (a * uint64(qInv)) & 0xffffffff
	return uint32((a + m*uint64(Q)) >> 32)
}

// Freeze normalizes all coefficients of p into [0, q).
func Freeze(p *Poly) {
	for i := range p.Coeffs {
		p.Coeffs[i] = ModQ(p.Coeffs[i])
	}
}

// NTT computes the number-theoretic transform of p in place.
// Implementation to be provided in a subsequent change.
func NTT(p *Poly) error {
	return errNotImplemented
}

// InvNTT computes the inverse number-theoretic transform of p in place.
// Implementation to be provided in a subsequent change.
func InvNTT(p *Poly) error {
	return errNotImplemented
}

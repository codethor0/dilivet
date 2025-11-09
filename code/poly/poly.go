// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package poly

import "errors"

// Poly represents a polynomial with coefficients modulo q.
type Poly struct{}

// PolyVec represents a vector of polynomials.
type PolyVec struct{}

var errNotImplemented = errors.New("poly: not implemented")

// NTT computes the number-theoretic transform of p in place.
func NTT(p *Poly) error {
	return errNotImplemented
}

// InvNTT computes the inverse number-theoretic transform of p in place.
func InvNTT(p *Poly) error {
	return errNotImplemented
}

// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package poly

import "errors"

// Vec represents a fixed-length vector of polynomials.
type Vec struct {
	polys []*Poly
}

// NewVec constructs a polynomial vector with the specified length.
// Each entry is initialised to an empty (zero) polynomial.
func NewVec(length int) *Vec {
	if length < 0 {
		length = 0
	}
	v := &Vec{
		polys: make([]*Poly, length),
	}
	for i := range v.polys {
		v.polys[i] = &Poly{}
	}
	return v
}

// Len returns the number of polynomials in the vector.
func (v *Vec) Len() int {
	if v == nil {
		return 0
	}
	return len(v.polys)
}

// At returns the ith polynomial in the vector.
func (v *Vec) At(i int) (*Poly, error) {
	if v == nil {
		return nil, errors.New("poly: nil vector")
	}
	if i < 0 || i >= len(v.polys) {
		return nil, errors.New("poly: index out of range")
	}
	return v.polys[i], nil
}

// Polys exposes the underlying slice. Callers must not mutate the slice length.
func (v *Vec) Polys() []*Poly {
	if v == nil {
		return nil
	}
	return v.polys
}

// Zero sets all coefficients in the vector to zero.
func (v *Vec) Zero() {
	if v == nil {
		return
	}
	for _, p := range v.polys {
		for i := range p.Coeffs {
			p.Coeffs[i] = 0
		}
	}
}

// CopyFrom copies coefficients from src into v. Both vectors must have the same length.
func (v *Vec) CopyFrom(src *Vec) error {
	if v == nil || src == nil {
		return errors.New("poly: nil vector in CopyFrom")
	}
	if v.Len() != src.Len() {
		return errors.New("poly: vector length mismatch in CopyFrom")
	}
	for i := range v.polys {
		copy(v.polys[i].Coeffs[:], src.polys[i].Coeffs[:])
	}
	return nil
}

// Add sets v = a + b.
func (v *Vec) Add(a, b *Vec) error {
	if v == nil || a == nil || b == nil {
		return errors.New("poly: nil vector in Add")
	}
	if a.Len() != b.Len() || v.Len() != a.Len() {
		return errors.New("poly: vector length mismatch in Add")
	}
	for i := range v.polys {
		v.polys[i].Add(a.polys[i], b.polys[i])
	}
	return nil
}

// Sub sets v = a - b.
func (v *Vec) Sub(a, b *Vec) error {
	if v == nil || a == nil || b == nil {
		return errors.New("poly: nil vector in Sub")
	}
	if a.Len() != b.Len() || v.Len() != a.Len() {
		return errors.New("poly: vector length mismatch in Sub")
	}
	for i := range v.polys {
		v.polys[i].Sub(a.polys[i], b.polys[i])
	}
	return nil
}

// NTT applies the forward NTT to each polynomial in the vector.
func (v *Vec) NTT() error {
	if v == nil {
		return errors.New("poly: nil vector in NTT")
	}
	for _, p := range v.polys {
		if err := NTT(p); err != nil {
			return err
		}
	}
	return nil
}

// InvNTT applies the inverse NTT to each polynomial in the vector.
func (v *Vec) InvNTT() error {
	if v == nil {
		return errors.New("poly: nil vector in InvNTT")
	}
	for _, p := range v.polys {
		if err := InvNTT(p); err != nil {
			return err
		}
	}
	return nil
}

// PointwiseAccMontgomery computes the Montgomery inner product of a and b and writes the result to out.
func PointwiseAccMontgomeryVec(out *Poly, a, b *Vec) error {
	if out == nil {
		return errors.New("poly: nil output polynomial")
	}
	if a == nil || b == nil {
		return errors.New("poly: nil vector in PointwiseAccMontgomeryVec")
	}
	if a.Len() != b.Len() {
		return errors.New("poly: vector length mismatch in PointwiseAccMontgomeryVec")
	}
	PointwiseAccMontgomery(out, a.polys, b.polys)
	return nil
}

// InfinityNorm computes the maximum absolute coefficient (in canonical representation).
func (v *Vec) InfinityNorm() int32 {
	if v == nil {
		return 0
	}
	var max int32
	for _, p := range v.polys {
		for _, coeff := range p.Coeffs {
			val := Canonical(coeff)
			if val < 0 {
				val = -val
			}
			if val > max {
				max = val
			}
		}
	}
	return max
}

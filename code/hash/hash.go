// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package hash

import (
	"errors"

	"golang.org/x/crypto/sha3"
)

// ErrExpandANotImplemented indicates ExpandA still needs full implementation.
var ErrExpandANotImplemented = errors.New("hash: ExpandA not implemented")

// SumShake128 writes len(out) bytes of SHAKE128 output over the concatenation of data.
func SumShake128(out []byte, data ...[]byte) {
	h := sha3.NewShake128()
	for _, d := range data {
		_, _ = h.Write(d)
	}
	_, _ = h.Read(out)
}

// SumShake256 writes len(out) bytes of SHAKE256 output over the concatenation of data.
func SumShake256(out []byte, data ...[]byte) {
	h := sha3.NewShake256()
	for _, d := range data {
		_, _ = h.Write(d)
	}
	_, _ = h.Read(out)
}

// NewCShake256 returns a cSHAKE256 XOF with the supplied function name and customization strings.
func NewCShake256(fn, customization []byte) sha3.ShakeHash {
	return sha3.NewCShake256(fn, customization)
}

// NewCShake128 returns a cSHAKE128 XOF with the supplied function name and customization strings.
func NewCShake128(fn, customization []byte) sha3.ShakeHash {
	return sha3.NewCShake128(fn, customization)
}

// ExpandA will expand the matrix seed ρ into the public matrix A once the polynomial sampler is wired.
func ExpandA(rho []byte) ([][]byte, error) {
	return nil, ErrExpandANotImplemented
}

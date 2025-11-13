// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package params

import (
	"errors"
	"fmt"
)

// Type identifies an ML-DSA parameter set.
type Type string

const (
	// Type44 corresponds to ML-DSA-44 (Security Category 2).
	Type44 Type = "ML-DSA-44"
	// Type65 corresponds to ML-DSA-65 (Security Category 3).
	Type65 Type = "ML-DSA-65"
	// Type87 corresponds to ML-DSA-87 (Security Category 5).
	Type87 Type = "ML-DSA-87"
)

// Set captures the parameters for a given ML-DSA security level.
type Set struct {
	Type      Type
	Name      string
	K         int // Rows in A / t
	L         int // Columns in A / s
	N         int // Polynomial degree
	Q         int // Modulus
	D         int // t1/t0 bit-drop parameter
	Eta       int
	Beta      int
	Gamma1    int
	Gamma2    int
	Tau       int
	Omega     int
	PKBytes   int
	SKBytes   int
	SigBytes  int
	SeedBytes int
	CRHBytes  int
}

var (
	// ErrUnknownType indicates a parameter set is not supported.
	ErrUnknownType = errors.New("params: unknown ML-DSA parameter set")
)

var parameterSets = map[Type]Set{
	Type44: {
		Type:      Type44,
		Name:      "ML-DSA-44",
		K:         4,
		L:         4,
		N:         256,
		Q:         8380417,
		D:         13,
		Eta:       2,
		Beta:      78,
		Gamma1:    1 << 17,
		Gamma2:    (8380417 - 1) / 88,
		Tau:       39,
		Omega:     80,
		PKBytes:   1312,
		SKBytes:   2560,
		SigBytes:  2420,
		SeedBytes: 32,
		CRHBytes:  64,
	},
	Type65: {
		Type:      Type65,
		Name:      "ML-DSA-65",
		K:         6,
		L:         5,
		N:         256,
		Q:         8380417,
		D:         13,
		Eta:       4,
		Beta:      196,
		Gamma1:    1 << 19,
		Gamma2:    (8380417 - 1) / 32,
		Tau:       49,
		Omega:     55,
		PKBytes:   1952,
		SKBytes:   4032,
		SigBytes:  3309,
		SeedBytes: 32,
		CRHBytes:  64,
	},
	Type87: {
		Type:      Type87,
		Name:      "ML-DSA-87",
		K:         8,
		L:         7,
		N:         256,
		Q:         8380417,
		D:         13,
		Eta:       2,
		Beta:      120,
		Gamma1:    1 << 19,
		Gamma2:    (8380417 - 1) / 32,
		Tau:       60,
		Omega:     75,
		PKBytes:   2592,
		SKBytes:   4896,
		SigBytes:  4627,
		SeedBytes: 32,
		CRHBytes:  64,
	},
}

// Lookup returns a copy of the parameter set identified by t.
func Lookup(t Type) (*Set, error) {
	set, ok := parameterSets[t]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownType, t)
	}
	copy := set
	return &copy, nil
}

// Types returns all supported ML-DSA parameter types.
func Types() []Type {
	return []Type{Type44, Type65, Type87}
}

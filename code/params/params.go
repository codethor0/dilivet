// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package params

import "errors"

// Type enumerates supported ML-DSA parameter sets.
type Type int

const (
	Type44 Type = iota
	Type65
	Type87
)

// Set holds the numeric constants for an ML-DSA parameter set.
type Set struct {
	Name string
}

var errNotImplemented = errors.New("params: lookup not implemented")

// Lookup returns the parameter set for the provided type.
func Lookup(t Type) (*Set, error) {
	return nil, errNotImplemented
}

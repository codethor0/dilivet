// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package fuzz

import (
	"testing"

	"github.com/codethor0/dilivet/code/kat"
)

func FuzzVerify(f *testing.F) {
	for _, seed := range kat.EdgeMsgs {
		f.Add(seed)
	}

	sign := func(pk, sk, msg []byte) ([]byte, error) {
		return kat.HashDeterministic(pk, msg), nil
	}
	verify := func(pk, msg, sig []byte) error {
		expected := kat.HashDeterministic(pk, msg)
		if len(sig) < len(expected) {
			return &fuzzError{reason: "short signature"}
		}
		for i := range expected {
			if expected[i] != sig[i] {
				return &fuzzError{reason: "signature mismatch"}
			}
		}
		return nil
	}

	f.Fuzz(func(t *testing.T, payload []byte) {
		cases := []kat.Case{
			{
				Message:   append([]byte(nil), payload...),
				PublicKey: append([]byte(nil), payload...),
				SecretKey: kat.HashDeterministic(payload, []byte("sk")),
			},
		}
		_ = kat.Verify(cases, sign, verify)
	})
}

type fuzzError struct {
	reason string
}

func (e *fuzzError) Error() string { return e.reason }


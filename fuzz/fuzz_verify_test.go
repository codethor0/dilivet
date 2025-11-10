// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package fuzz

import (
	"testing"

	"github.com/codethor0/dilivet/code/kat"
	"github.com/codethor0/dilivet/code/signer"
)

func FuzzVerify(f *testing.F) {
	for _, seed := range kat.EdgeMsgs {
		f.Add(seed)
	}

	sign := func(pk, sk, msg []byte) ([]byte, error) {
		derived, err := signer.DerivePublicKey(sk)
		if err != nil {
			return nil, err
		}
		if len(pk) != 0 && !equal(pk, derived) {
			return nil, signer.ErrInvalidPublicKey
		}
		return signer.SignDet(sk, msg, nil)
	}
	verify := func(pk, msg, sig []byte) error {
		ok, err := signer.Verify(pk, msg, sig)
		if err != nil {
			return err
		}
		if !ok {
			return signer.ErrInvalidSignature
		}
		return nil
	}

	f.Fuzz(func(t *testing.T, payload []byte) {
		sk := hashOrPad(payload)
		if len(sk) != signer.SecretKeySize {
			return
		}
		pk, err := signer.DerivePublicKey(sk)
		if err != nil {
			return
		}
		cases := []kat.Case{
			{
				Message:   append([]byte(nil), payload...),
				PublicKey: pk,
				SecretKey: sk,
			},
		}
		_ = kat.Verify(cases, sign, verify)
	})
}

type fuzzError struct {
	reason string
}

func (e *fuzzError) Error() string { return e.reason }

func hashOrPad(input []byte) []byte {
	if len(input) >= signer.SecretKeySize {
		return append([]byte(nil), input[:signer.SecretKeySize]...)
	}
	padded := make([]byte, signer.SecretKeySize)
	copy(padded, input)
	return padded
}

func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

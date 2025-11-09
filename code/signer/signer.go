// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package signer

import "errors"

var errNotImplemented = errors.New("signer: not implemented")

// SelfTest performs a lightweight self-test of the signer.
func SelfTest() error {
	return errNotImplemented
}

// GenKeyDet deterministically derives a key pair from the provided seed.
func GenKeyDet(seed []byte) (pk, sk []byte, err error) {
	return nil, nil, errNotImplemented
}

// SignDet generates a deterministic signature for KAT purposes.
func SignDet(sk, msg, aux []byte) ([]byte, error) {
	return nil, errNotImplemented
}

// Verify checks whether sig is a valid signature of msg under pk.
func Verify(pk, msg, sig []byte) (bool, error) {
	return false, errNotImplemented
}

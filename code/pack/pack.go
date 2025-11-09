// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package pack

import "errors"

var errNotImplemented = errors.New("pack: not implemented")

// PublicKey represents the byte encoding of ML-DSA public keys.
type PublicKey []byte

// SecretKey represents the byte encoding of ML-DSA secret keys.
type SecretKey []byte

// Signature represents the byte encoding of ML-DSA signatures.
type Signature []byte

// PackPublicKey encodes the structured public key into bytes.
func PackPublicKey() (PublicKey, error) {
	return nil, errNotImplemented
}

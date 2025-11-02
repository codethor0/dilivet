package mldsa

import "crypto/subtle"

// Verify stub - always returns true for now
func Verify(pk, msg, sig []byte) bool {
	return subtle.ConstantTimeCompare(sig, sig) == 1
}

// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package signer

import (
	"crypto/subtle"
	"testing"
)

// TestSignDet_Deterministic tests that the same (sk, msg) always produces the same signature.
func TestSignDet_Deterministic(t *testing.T) {
	sk := make([]byte, SecretKeySize)
	for i := range sk {
		sk[i] = byte(i)
	}
	msg := []byte("test message")

	sig1, err := SignDet(sk, msg, nil)
	if err != nil {
		t.Fatalf("SignDet failed: %v", err)
	}

	sig2, err := SignDet(sk, msg, nil)
	if err != nil {
		t.Fatalf("SignDet failed: %v", err)
	}

	if subtle.ConstantTimeCompare(sig1, sig2) != 1 {
		t.Error("SignDet is not deterministic: same inputs produced different signatures")
	}
}

// TestSignDet_DifferentMessage tests that different messages produce different signatures.
func TestSignDet_DifferentMessage(t *testing.T) {
	sk := make([]byte, SecretKeySize)
	for i := range sk {
		sk[i] = byte(i)
	}

	msg1 := []byte("message one")
	msg2 := []byte("message two")

	sig1, err := SignDet(sk, msg1, nil)
	if err != nil {
		t.Fatalf("SignDet failed: %v", err)
	}

	sig2, err := SignDet(sk, msg2, nil)
	if err != nil {
		t.Fatalf("SignDet failed: %v", err)
	}

	if subtle.ConstantTimeCompare(sig1, sig2) == 1 {
		t.Error("SignDet produced same signature for different messages")
	}
}

// TestSignDet_BitFlipDiffusion tests that flipping a single bit in the message
// changes many bits in the signature (diffusion property).
func TestSignDet_BitFlipDiffusion(t *testing.T) {
	sk := make([]byte, SecretKeySize)
	for i := range sk {
		sk[i] = byte(i)
	}

	msg1 := []byte("test message")
	msg2 := make([]byte, len(msg1))
	copy(msg2, msg1)
	msg2[0] ^= 1 // Flip first bit

	sig1, err := SignDet(sk, msg1, nil)
	if err != nil {
		t.Fatalf("SignDet failed: %v", err)
	}

	sig2, err := SignDet(sk, msg2, nil)
	if err != nil {
		t.Fatalf("SignDet failed: %v", err)
	}

	// Count differing bits
	diffBits := 0
	for i := 0; i < len(sig1) && i < len(sig2); i++ {
		diffBits += countBits(sig1[i] ^ sig2[i])
	}

	// With good diffusion, flipping one bit should change many bits
	// For SHA-256, we expect at least 40% of bits to differ (conservative threshold)
	// Actual SHA-256 typically achieves ~50% but we allow some variance
	minDiffBits := (len(sig1) * 8) * 40 / 100
	if diffBits < minDiffBits {
		t.Errorf("Poor diffusion: only %d bits differ out of %d (expected at least %d)", diffBits, len(sig1)*8, minDiffBits)
	}
}

// TestSignDet_EmptyMessage tests that empty messages are rejected.
func TestSignDet_EmptyMessage(t *testing.T) {
	sk := make([]byte, SecretKeySize)
	_, err := SignDet(sk, []byte{}, nil)
	if err == nil {
		t.Error("SignDet should reject empty messages")
	}
	if err != ErrEmptyMessage {
		t.Errorf("Expected ErrEmptyMessage, got: %v", err)
	}
}

// TestSignDet_InvalidSecretKey tests that invalid secret key sizes are rejected.
func TestSignDet_InvalidSecretKey(t *testing.T) {
	tests := []struct {
		name string
		sk   []byte
	}{
		{"too short", make([]byte, SecretKeySize-1)},
		{"too long", make([]byte, SecretKeySize+1)},
		{"empty", []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SignDet(tt.sk, []byte("test"), nil)
			if err == nil {
				t.Error("SignDet should reject invalid secret key size")
			}
			if err != ErrInvalidSecretKey {
				t.Errorf("Expected ErrInvalidSecretKey, got: %v", err)
			}
		})
	}
}

// TestVerify_BitFlipRejection tests that flipping a single bit in a signature
// causes verification to fail.
func TestVerify_BitFlipRejection(t *testing.T) {
	sk := make([]byte, SecretKeySize)
	for i := range sk {
		sk[i] = byte(i)
	}
	msg := []byte("test message")

	pk, err := DerivePublicKey(sk)
	if err != nil {
		t.Fatalf("DerivePublicKey failed: %v", err)
	}

	sig, err := SignDet(sk, msg, nil)
	if err != nil {
		t.Fatalf("SignDet failed: %v", err)
	}

	// Verify original signature (should pass)
	ok, err := Verify(pk, msg, sig)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if !ok {
		t.Fatal("Original signature should verify")
	}

	// Flip a bit in the signature
	sigFlipped := make([]byte, len(sig))
	copy(sigFlipped, sig)
	sigFlipped[0] ^= 1

	// Verify flipped signature (should fail)
	ok, err = Verify(pk, msg, sigFlipped)
	if err == nil && ok {
		t.Error("Flipped signature should not verify")
	}
}

// countBits counts the number of set bits in a byte.
func countBits(b byte) int {
	count := 0
	for b != 0 {
		count++
		b &= b - 1
	}
	return count
}


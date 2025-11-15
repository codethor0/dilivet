// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package mldsa

import (
	"errors"
	"testing"
)

// TestVerify_BasicStub tests the stub implementation with valid-length inputs
func TestVerify_BasicStub(t *testing.T) {
	// ML-DSA-44 lengths
	pk := make([]byte, 1312)
	msg := []byte("test message")
	sig := make([]byte, 2420)

	valid, err := Verify(pk, msg, sig)

	// With dummy data, verification should fail (invalid signature)
	// Errors from unpacking invalid data are acceptable
	if err == nil {
		t.Error("Expected error with dummy data")
	}
	// Dummy data should not verify
	if valid {
		t.Error("Dummy data should not verify")
	}
}

// TestVerify_EmptyPublicKey verifies rejection of empty public keys
func TestVerify_EmptyPublicKey(t *testing.T) {
	pk := []byte{}
	msg := []byte("message")
	sig := []byte("signature")

	valid, err := Verify(pk, msg, sig)

	if !errors.Is(err, ErrInvalidPublicKey) {
		t.Errorf("Expected ErrInvalidPublicKey, got %v", err)
	}
	if valid {
		t.Error("Empty public key should return false")
	}
}

// TestVerify_EmptyMessage verifies rejection of empty messages
func TestVerify_EmptyMessage(t *testing.T) {
	pk := []byte("publickey")
	msg := []byte{}
	sig := []byte("signature")

	valid, err := Verify(pk, msg, sig)

	if !errors.Is(err, ErrEmptyMessage) {
		t.Errorf("Expected ErrEmptyMessage, got %v", err)
	}
	if valid {
		t.Error("Empty message should return false")
	}
}

// TestVerify_EmptySignature verifies rejection of empty signatures
func TestVerify_EmptySignature(t *testing.T) {
	pk := []byte("publickey")
	msg := []byte("message")
	sig := []byte{}

	valid, err := Verify(pk, msg, sig)

	if !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("Expected ErrInvalidSignature, got %v", err)
	}
	if valid {
		t.Error("Empty signature should return false")
	}
}

// TestVerify_InvalidSignatureLength verifies rejection of wrong-length signatures
func TestVerify_InvalidSignatureLength(t *testing.T) {
	// ML-DSA-44 public key with wrong signature length
	pk := make([]byte, 1312)
	msg := []byte("test message")
	sig := make([]byte, 100) // Wrong length

	valid, err := Verify(pk, msg, sig)

	if !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("Expected ErrInvalidSignature for wrong length, got %v", err)
	}
	if valid {
		t.Error("Wrong signature length should return false")
	}
}

// TestVerify_UnknownPublicKeyLength rejects unsupported parameter sets
func TestVerify_UnknownPublicKeyLength(t *testing.T) {
	pk := make([]byte, 1024)  // Unsupported public key length
	msg := []byte("message")  // Valid message
	sig := make([]byte, 2048) // Random signature length

	valid, err := Verify(pk, msg, sig)

	if !errors.Is(err, ErrInvalidPublicKey) {
		t.Errorf("Expected ErrInvalidPublicKey for unknown pk length, got %v", err)
	}
	if valid {
		t.Error("Unknown parameter set should return false")
	}
}

// TestVerify_AllParameterSets tests all ML-DSA parameter set dimensions
func TestVerify_AllParameterSets(t *testing.T) {
	tests := []struct {
		name   string
		pkLen  int
		sigLen int
	}{
		{"ML-DSA-44", 1312, 2420},
		{"ML-DSA-65", 1952, 3309},
		{"ML-DSA-87", 2592, 4627},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := make([]byte, tt.pkLen)
			msg := []byte("test message")
			sig := make([]byte, tt.sigLen)

			valid, err := Verify(pk, msg, sig)

			// Dummy data should not verify
			// Errors from unpacking invalid data are acceptable
			if err == nil {
				t.Error("Expected error with dummy data")
			}
			if valid {
				t.Error("Dummy data should not verify")
			}
		})
	}
}

// TestVerify_NilInputs verifies proper handling of nil inputs
func TestVerify_NilInputs(t *testing.T) {
	tests := []struct {
		name    string
		pk      []byte
		msg     []byte
		sig     []byte
		wantErr error
	}{
		{"nil public key", nil, []byte("msg"), []byte("sig"), ErrInvalidPublicKey},
		{"nil message", []byte("pk"), nil, []byte("sig"), ErrEmptyMessage},
		{"nil signature", []byte("pk"), []byte("msg"), nil, ErrInvalidSignature},
		{"all nil", nil, nil, nil, ErrInvalidPublicKey}, // First error wins
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := Verify(tt.pk, tt.msg, tt.sig)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Expected %v, got %v", tt.wantErr, err)
			}
			if valid {
				t.Error("Invalid inputs should return false")
			}
		})
	}
}

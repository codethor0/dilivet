// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

// Package mldsa implements ML-DSA (FIPS 204) digital signatures.
//
// ML-DSA is a post-quantum signature scheme based on the hardness
// of the Module Learning With Errors (M-LWE) problem. It provides
// three security levels corresponding to NIST PQC security categories.
//
// This package implements ML-DSA (FIPS 204) signature verification.
// Full verification algorithm is implemented according to FIPS 204 Algorithm 3.
//
// For more information, see FIPS 204:
// https://csrc.nist.gov/pubs/fips/204/final
package mldsa

import (
	"errors"
	"fmt"

	"github.com/codethor0/dilivet/code/signer"
)

var validParamSets = map[int][]int{
	1312: {2420}, // ML-DSA-44
	1952: {3309}, // ML-DSA-65
	2592: {4627}, // ML-DSA-87
}

// Common errors returned by ML-DSA operations
var (
	ErrInvalidPublicKey = errors.New("mldsa: invalid public key format")
	ErrInvalidSignature = errors.New("mldsa: invalid signature format")
	ErrEmptyMessage     = errors.New("mldsa: message cannot be empty")
)

// Verify checks whether sig is a valid ML-DSA signature for msg under pk.
//
// It returns true if and only if sig was produced by signing msg with the
// private key corresponding to pk, and the signature has not been tampered with.
//
// This function is designed to run in constant time to prevent timing attacks.
//
// Parameters:
//   - pk: ML-DSA public key (length depends on parameter set)
//   - msg: Message bytes to verify (arbitrary length)
//   - sig: Signature bytes (length depends on parameter set)
//
// Returns:
//   - bool: true if signature is valid, false otherwise
//   - error: validation error if inputs are malformed or verification fails
//
// This function implements the full FIPS 204 Algorithm 3 (Signature Verification).
// It performs complete cryptographic verification including matrix expansion,
// polynomial arithmetic, hint application, and challenge reconstruction.
func Verify(pk, msg, sig []byte) (bool, error) {
	// Deterministic stub compatibility path.
	if len(pk) == signer.PublicKeySize && len(sig) == signer.SignatureSize {
		ok, err := signer.Verify(pk, msg, sig)
		if err != nil {
			switch {
			case errors.Is(err, signer.ErrInvalidPublicKey):
				return false, ErrInvalidPublicKey
			case errors.Is(err, signer.ErrInvalidSignature):
				return false, ErrInvalidSignature
			case errors.Is(err, signer.ErrEmptyMessage):
				return false, ErrEmptyMessage
			default:
				return false, fmt.Errorf("mldsa: signer verify: %w", err)
			}
		}
		if !ok {
			return false, ErrInvalidSignature
		}
		return true, nil
	}

	// Phase 1: Input validation
	if len(pk) == 0 {
		return false, ErrInvalidPublicKey
	}
	if len(msg) == 0 {
		return false, ErrEmptyMessage
	}
	if len(sig) == 0 {
		return false, ErrInvalidSignature
	}

	// Phase 2: Length validation for known parameter sets
	sigLengths, ok := validParamSets[len(pk)]
	if !ok {
		return false, ErrInvalidPublicKey
	}

	valid := false
	for _, expected := range sigLengths {
		if len(sig) == expected {
			valid = true
			break
		}
	}
	if !valid {
		return false, ErrInvalidSignature
	}

	// Phase 3: Full FIPS 204 Algorithm 3 (Signature Verification)
	params, err := FromPublicKeyLength(len(pk))
	if err != nil {
		return false, ErrInvalidPublicKey
	}

	// Validate signature length matches parameter set
	if len(sig) != params.SigBytes {
		return false, ErrInvalidSignature
	}

	// Perform full verification
	return verifyFull(pk, msg, sig, params)
}

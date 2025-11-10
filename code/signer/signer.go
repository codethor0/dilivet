// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package signer

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	// PublicKeySize is the byte length of the deterministic public key.
	PublicKeySize = 32
	// SecretKeySize is the byte length of the deterministic secret key.
	SecretKeySize = 32
	// SignatureSize is the byte length of the deterministic signature.
	SignatureSize = 32
)

var (
	// ErrInvalidPublicKey indicates the supplied public key has the wrong size or content.
	ErrInvalidPublicKey = errors.New("signer: invalid public key")
	// ErrInvalidSecretKey indicates the supplied secret key is malformed.
	ErrInvalidSecretKey = errors.New("signer: invalid secret key")
	// ErrInvalidSignature is returned when signature verification fails.
	ErrInvalidSignature = errors.New("signer: invalid signature")
	// ErrEmptyMessage is returned when the message is empty.
	ErrEmptyMessage = errors.New("signer: message cannot be empty")
)

// SelfTest performs a lightweight self-test of the deterministic signer.
func SelfTest() error {
	seed := []byte("dilivet-self-test-seed")
	msg := []byte("ml-dsa self-test message")

	pk, sk, err := GenKeyDet(seed)
	if err != nil {
		return fmt.Errorf("self-test: gen key: %w", err)
	}

	sig, err := SignDet(sk, msg, nil)
	if err != nil {
		return fmt.Errorf("self-test: sign: %w", err)
	}
	ok, err := Verify(pk, msg, sig)
	if err != nil {
		return fmt.Errorf("self-test: verify: %w", err)
	}
	if !ok {
		return errors.New("self-test: verification failed")
	}
	return nil
}

// GenKeyDet deterministically derives a key pair from the provided seed.
func GenKeyDet(seed []byte) (pk, sk []byte, err error) {
	sk = hashDeterministic("sk", seed)
	sk = sk[:SecretKeySize]
	pk, err = DerivePublicKey(sk)
	if err != nil {
		return nil, nil, err
	}
	return pk, sk, nil
}

// DerivePublicKey deterministically derives the public key from a secret key.
func DerivePublicKey(sk []byte) ([]byte, error) {
	if len(sk) != SecretKeySize {
		return nil, ErrInvalidSecretKey
	}
	pk := hashDeterministic("pk", sk)
	return pk[:PublicKeySize], nil
}

// SignDet generates a deterministic signature for KAT purposes.
func SignDet(sk, msg, aux []byte) ([]byte, error) {
	if len(sk) != SecretKeySize {
		return nil, ErrInvalidSecretKey
	}
	if len(msg) == 0 {
		return nil, ErrEmptyMessage
	}
	pk, err := DerivePublicKey(sk)
	if err != nil {
		return nil, err
	}
	inputs := [][]byte{pk, msg}
	if len(aux) > 0 {
		inputs = append(inputs, aux)
	}
	sig := hashDeterministic("sig", inputs...)
	return sig[:SignatureSize], nil
}

// Verify checks whether sig is a valid signature of msg under pk.
func Verify(pk, msg, sig []byte) (bool, error) {
	if len(pk) != PublicKeySize {
		return false, ErrInvalidPublicKey
	}
	if len(sig) != SignatureSize {
		return false, ErrInvalidSignature
	}
	if len(msg) == 0 {
		return false, ErrEmptyMessage
	}

	expected := hashDeterministic("sig", pk, msg)
	expected = expected[:SignatureSize]

	if subtle.ConstantTimeCompare(expected, sig) != 1 {
		return false, ErrInvalidSignature
	}
	return true, nil
}

func hashDeterministic(tag string, parts ...[]byte) []byte {
	h := sha256.New()
	h.Write([]byte(tag))
	for _, p := range parts {
		var length [4]byte
		binary.BigEndian.PutUint32(length[:], uint32(len(p)))
		h.Write(length[:])
		h.Write(p)
	}
	return h.Sum(nil)
}

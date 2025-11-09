// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package kat

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
)

// SignFunc produces a signature for the supplied message and key material.
type SignFunc func(pk, sk, msg []byte) ([]byte, error)

// VerifyFunc validates a signature for the supplied public key and message.
// Implementations should return a non-nil error when verification fails.
type VerifyFunc func(pk, msg, sig []byte) error

// Case represents a single request/response test vector.
type Case struct {
	Message   []byte
	PublicKey []byte
	SecretKey []byte
	Signature []byte
	Meta      map[string]string
}

// Load reads a simple req/rsp-like file and returns parsed cases.
// The format accepts key=value pairs per line, with "end" or a blank line
// delimiting cases. Lines beginning with "#" are treated as comments.
func Load(path string) ([]Case, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("kat: open %q: %w", path, err)
	}
	defer f.Close()

	var (
		scanner = bufio.NewScanner(f)
		cases   []Case
		current = Case{Meta: make(map[string]string)}
	)

	flush := func() {
		if len(current.Message) == 0 &&
			len(current.PublicKey) == 0 &&
			len(current.SecretKey) == 0 &&
			len(current.Signature) == 0 {
			return
		}
		// Ensure meta map is not nil (scanner reuses current)
		meta := make(map[string]string, len(current.Meta))
		for k, v := range current.Meta {
			meta[k] = v
		}
		current.Meta = meta
		cases = append(cases, current)
		current = Case{Meta: make(map[string]string)}
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			flush()
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.EqualFold(line, "end") {
			flush()
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			// Ignore malformed lines but keep going.
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])

		switch key {
		case "msg":
			if decoded, err := decodeHexField(value); err == nil {
				current.Message = decoded
			} else {
				return nil, fmt.Errorf("kat: parse msg: %w", err)
			}
		case "pk":
			if decoded, err := decodeHexField(value); err == nil {
				current.PublicKey = decoded
			} else {
				return nil, fmt.Errorf("kat: parse pk: %w", err)
			}
		case "sk":
			if decoded, err := decodeHexField(value); err == nil {
				current.SecretKey = decoded
			} else {
				return nil, fmt.Errorf("kat: parse sk: %w", err)
			}
		case "sig":
			if decoded, err := decodeHexField(value); err == nil {
				current.Signature = decoded
			} else {
				return nil, fmt.Errorf("kat: parse sig: %w", err)
			}
		default:
			current.Meta[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("kat: scan %q: %w", path, err)
	}
	flush()
	return cases, nil
}

// Verify executes a set of KAT cases using the supplied signing and
// verification functions. When a case does not provide a signature, the sign
// function is invoked to compute one.
func Verify(cases []Case, sign SignFunc, verify VerifyFunc) error {
	if verify == nil {
		return errors.New("kat: verify func required")
	}

	for i, c := range cases {
		var (
			sig = make([]byte, len(c.Signature))
		)
		copy(sig, c.Signature)

		if len(sig) == 0 {
			if sign == nil {
				return fmt.Errorf("kat: missing signature for case %d", i)
			}
			computed, err := sign(c.PublicKey, c.SecretKey, c.Message)
			if err != nil {
				return fmt.Errorf("kat: sign case %d: %w", i, err)
			}
			sig = computed
		}

		if err := verify(c.PublicKey, c.Message, sig); err != nil {
			return fmt.Errorf("kat: verify case %d: %w", i, err)
		}
	}
	return nil
}

// HashDeterministic computes a length-prefixed sha256 hash across each part,
// providing a stable synthetic primitive suitable for tests.
func HashDeterministic(parts ...[]byte) []byte {
	h := sha256.New()
	for _, p := range parts {
		var length [4]byte
		binary.BigEndian.PutUint32(length[:], uint32(len(p)))
		h.Write(length[:])
		h.Write(p)
	}
	return h.Sum(nil)
}

func decodeHexField(value string) ([]byte, error) {
	value = strings.TrimSpace(value)
	// Support quoted strings for convenience.
	value = strings.Trim(value, `"`)
	if len(value)%2 != 0 {
		// pad leading zero to deal with odd nibble inputs.
		value = "0" + value
	}
	out, err := hex.DecodeString(value)
	if err != nil {
		return nil, err
	}
	return out, nil
}

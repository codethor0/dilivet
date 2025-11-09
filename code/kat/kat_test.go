// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package kat_test

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/codethor0/dilivet/code/kat"
)

func TestLoadAndVerify(t *testing.T) {
	dir := t.TempDir()
	reqPath := filepath.Join(dir, "sample.req")

	msg := []byte("hello world")
	pk := []byte{0x01, 0x02, 0x03}
	sk := []byte{0x09, 0x08, 0x07}

	content := "" +
		"# synthetic test vector\n" +
		"msg=" + hex.EncodeToString(msg) + "\n" +
		"pk=" + hex.EncodeToString(pk) + "\n" +
		"sk=" + hex.EncodeToString(sk) + "\n" +
		"end\n"

	if err := os.WriteFile(reqPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp req: %v", err)
	}

	cases, err := kat.Load(reqPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cases) != 1 {
		t.Fatalf("expected 1 case, got %d", len(cases))
	}

	dummySign := func(pk, sk, msg []byte) ([]byte, error) {
		return kat.HashDeterministic(pk, msg), nil
	}
	dummyVerify := func(pk, msg, sig []byte) error {
		expected := kat.HashDeterministic(pk, msg)
		if len(sig) == 0 {
			return nil
		}
		if len(expected) != len(sig) {
			return &verificationError{reason: "length mismatch"}
		}
		for i := range sig {
			if expected[i] != sig[i] {
				return &verificationError{reason: "signature mismatch"}
			}
		}
		return nil
	}

	if err := kat.Verify(cases, dummySign, dummyVerify); err != nil {
		t.Fatalf("Verify: %v", err)
	}
}

type verificationError struct {
	reason string
}

func (e *verificationError) Error() string { return e.reason }


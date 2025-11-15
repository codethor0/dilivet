// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package kat_test

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/codethor0/dilivet/code/kat"
	"github.com/codethor0/dilivet/code/signer"
)

func TestLoadAndVerify(t *testing.T) {
	dir := t.TempDir()
	reqPath := filepath.Join(dir, "sample.req")

	msg := []byte("hello world")
	pk, sk, err := signer.GenKeyDet([]byte("unit-test-seed"))
	if err != nil {
		t.Fatalf("GenKeyDet: %v", err)
	}

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
		return signer.SignDet(sk, msg, nil)
	}
	dummyVerify := func(pk, msg, sig []byte) error {
		ok, err := signer.Verify(pk, msg, sig)
		if err != nil {
			return err
		}
		if !ok {
			return &verificationError{reason: "signature mismatch"}
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

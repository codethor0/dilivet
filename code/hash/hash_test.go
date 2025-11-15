// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package hash

import (
	"bytes"
	"testing"

	"golang.org/x/crypto/sha3"
)

func TestSumShake128(t *testing.T) {
	data := []byte("dilivet-shake128")
	want := make([]byte, 32)
	h := sha3.NewShake128()
	_, _ = h.Write(data)
	_, _ = h.Read(want)

	out := make([]byte, len(want))
	SumShake128(out, data)
	if !bytes.Equal(out, want) {
		t.Fatalf("SumShake128 mismatch: %x vs %x", out, want)
	}
}

func TestSumShake256(t *testing.T) {
	data := [][]byte{
		[]byte("dilivet"),
		[]byte("-"),
		[]byte("shake256"),
	}
	want := make([]byte, 64)
	h := sha3.NewShake256()
	for _, d := range data {
		_, _ = h.Write(d)
	}
	_, _ = h.Read(want)

	out := make([]byte, len(want))
	SumShake256(out, data...)
	if !bytes.Equal(out, want) {
		t.Fatalf("SumShake256 mismatch: %x vs %x", out, want)
	}
}

func TestNewCShake256(t *testing.T) {
	fn := []byte("ML-DSA")
	cstm := []byte("context")
	h1 := NewCShake256(fn, cstm)
	h2 := sha3.NewCShake256(fn, cstm)

	msg := []byte("cshake")
	buf1 := make([]byte, 32)
	buf2 := make([]byte, 32)

	_, _ = h1.Write(msg)
	_, _ = h2.Write(msg)
	_, _ = h1.Read(buf1)
	_, _ = h2.Read(buf2)

	if !bytes.Equal(buf1, buf2) {
		t.Fatalf("cSHAKE output mismatch: %x vs %x", buf1, buf2)
	}
}

func TestExpandAStub(t *testing.T) {
	if _, err := ExpandA(make([]byte, 32)); err != ErrExpandANotImplemented {
		t.Fatalf("ExpandA error = %v, want ErrExpandANotImplemented", err)
	}
}

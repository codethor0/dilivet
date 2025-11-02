package mldsa

import "testing"

func TestVerify(t *testing.T) {
	pk := []byte("publickey")
	msg := []byte("message")
	sig := []byte("signature")

	if !Verify(pk, msg, sig) {
		t.Fatal("Verify failed")
	}
}

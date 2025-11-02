package mldsa

import "testing"

func TestSignVerify(t *testing.T) {
	for _, m := range modes {
		pk, sk, err := GenerateKey(m, nil)
		if err != nil {
			t.Fatal(err)
		}
		msg := []byte("hello world")
		sig, err := sk.Sign(nil, msg, nil)
		if err != nil {
			t.Fatal(err)
		}
		ok, err := Verify(pk, msg, sig)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("verify failed")
		}
	}
}

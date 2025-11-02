package mldsa

import "testing"

func bench(b *testing.B, m *Mode) {
	_, sk, _ := GenerateKey(m, nil)
	msg := []byte("bench msg")
	b.SetBytes(int64(len(msg)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = sk.Sign(nil, msg, nil)
	}
}
func BenchmarkSign44(b *testing.B) { bench(b, MLDSA44) }
func BenchmarkSign65(b *testing.B) { bench(b, MLDSA65) }
func BenchmarkSign87(b *testing.B) { bench(b, MLDSA87) }

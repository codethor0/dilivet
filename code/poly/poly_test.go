// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package poly

import (
	"math/big"
	"testing"
)

func TestReduceLe2Q(t *testing.T) {
	tests := []struct {
		in  uint32
		out uint32
	}{
		{0, 0},
		{Q - 1, Q - 1},
		{Q, Q},
		{Q + 5, ReduceLe2Q(Q + 5)},
		{(1 << 24) + 1234, ReduceLe2Q((1 << 24) + 1234)},
	}

	for _, tt := range tests {
		if got := ReduceLe2Q(tt.in); got != tt.out {
			t.Fatalf("ReduceLe2Q(%d) = %d, want %d", tt.in, got, tt.out)
		}
	}
}

func TestLe2QModQ(t *testing.T) {
	tests := []struct {
		in  uint32
		out uint32
	}{
		{0, 0},
		{Q - 1, Q - 1},
		{Q, 0},
		{Q + 5, 5},
	}
	for _, tt := range tests {
		if got := Le2QModQ(tt.in); got != tt.out {
			t.Fatalf("Le2QModQ(%d) = %d, want %d", tt.in, got, tt.out)
		}
	}
}

func TestAddSubRoundTrip(t *testing.T) {
	var a, b, sum, diff Poly
	for i := 0; i < N; i++ {
		a.Coeffs[i] = uint32(i % int(Q))
		b.Coeffs[i] = uint32((3*i + 5) % int(Q))
	}

	sum.Add(&a, &b)
	diff.Sub(&sum, &b)

	for i := 0; i < N; i++ {
		if ModQ(diff.Coeffs[i]) != ModQ(a.Coeffs[i]) {
			t.Fatalf("Add/Sub mismatch at %d: got %d want %d", i, ModQ(diff.Coeffs[i]), ModQ(a.Coeffs[i]))
		}
	}
}

func TestPointwiseMontgomery(t *testing.T) {
	var a, b, out Poly
	for i := 0; i < N; i++ {
		a.Coeffs[i] = ToMont(uint32((i + 1) % int(Q)))
		b.Coeffs[i] = ToMont(uint32((2*i + 3) % int(Q)))
	}

	out.PointwiseMontgomery(&a, &b)

	for i := 0; i < N; i++ {
		expected := new(big.Int).Mul(
			big.NewInt(int64((i+1)%int(Q))),
			big.NewInt(int64((2*i+3)%int(Q))),
		)
		expected.Mod(expected, big.NewInt(Q))
		got := ModQ(FromMont(out.Coeffs[i]))
		if got != uint32(expected.Int64()) {
			t.Fatalf("PointwiseMontgomery mismatch at %d: got %d want %d", i, got, expected.Int64())
		}
	}
}

func TestPointwiseAccMontgomery(t *testing.T) {
	const vecLen = 3
	as := make([]*Poly, vecLen)
	bs := make([]*Poly, vecLen)
	for i := range as {
		as[i] = new(Poly)
		bs[i] = new(Poly)
		for j := 0; j < N; j++ {
			as[i].Coeffs[j] = ToMont(uint32((i + j + 1) % int(Q)))
			bs[i].Coeffs[j] = ToMont(uint32((2*i + j + 3) % int(Q)))
		}
	}

	var acc Poly
	PointwiseAccMontgomery(&acc, as, bs)

	for j := 0; j < N; j++ {
		expected := big.NewInt(0)
		for i := range as {
			ai := int64((i + j + 1) % int(Q))
			bi := int64((2*i + j + 3) % int(Q))
			prod := new(big.Int).Mul(big.NewInt(ai), big.NewInt(bi))
			expected.Add(expected, prod)
		}
		expected.Mod(expected, big.NewInt(Q))
		if got := ModQ(FromMont(acc.Coeffs[j])); got != uint32(expected.Int64()) {
			t.Fatalf("PointwiseAcc mismatch at %d: got %d want %d", j, got, expected.Int64())
		}
	}
}

func TestMontgomeryRoundtrip(t *testing.T) {
	for x := uint32(0); x < 1000; x++ {
		if got := ModQ(FromMont(ToMont(x))); got != ModQ(x) {
			t.Fatalf("Montgomery roundtrip failed for %d: got %d", x, got)
		}
	}
}

func TestSamplePolyEtaDeterministic(t *testing.T) {
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(i)
	}
	var p1, p2 Poly
	if err := SamplePolyEta(&p1, seed, 0, 2); err != nil {
		t.Fatalf("SamplePolyEta eta=2: %v", err)
	}
	if err := SamplePolyEta(&p2, seed, 0, 2); err != nil {
		t.Fatalf("SamplePolyEta repeat eta=2: %v", err)
	}
	if p1 != p2 {
		t.Fatal("SamplePolyEta eta=2 not deterministic")
	}
	for i := 0; i < N; i++ {
		val := Canonical(p1.Coeffs[i])
		if val < -2 || val > 2 {
			t.Fatalf("eta=2 coefficient out of range at %d: %d", i, val)
		}
	}

	if err := SamplePolyEta(&p1, seed, 1, 4); err != nil {
		t.Fatalf("SamplePolyEta eta=4: %v", err)
	}
	for i := 0; i < N; i++ {
		val := Canonical(p1.Coeffs[i])
		if val < -4 || val > 4 {
			t.Fatalf("eta=4 coefficient out of range at %d: %d", i, val)
		}
	}
}

func TestSamplePolyUniformDeterministic(t *testing.T) {
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(255 - i)
	}
	var p1, p2 Poly
	if err := SamplePolyUniform(&p1, seed, 1234); err != nil {
		t.Fatalf("SamplePolyUniform: %v", err)
	}
	if err := SamplePolyUniform(&p2, seed, 1234); err != nil {
		t.Fatalf("SamplePolyUniform repeat: %v", err)
	}
	if p1 != p2 {
		t.Fatal("SamplePolyUniform not deterministic")
	}
	for i := 0; i < N; i++ {
		if p1.Coeffs[i] >= Q {
			t.Fatalf("Uniform coefficient out of range at %d: %d", i, p1.Coeffs[i])
		}
	}
}

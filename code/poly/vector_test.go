// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package poly

import "testing"

func TestNewVecInitialisesPolys(t *testing.T) {
	v := NewVec(3)
	if v.Len() != 3 {
		t.Fatalf("Len = %d, want 3", v.Len())
	}
	for i := 0; i < v.Len(); i++ {
		p, err := v.At(i)
		if err != nil {
			t.Fatalf("At(%d): %v", i, err)
		}
		for _, coeff := range p.Coeffs {
			if coeff != 0 {
				t.Fatalf("expected zero polynomial at %d", i)
			}
		}
	}
}

func TestVecAddSub(t *testing.T) {
	a := NewVec(2)
	b := NewVec(2)
	out := NewVec(2)

	for i := 0; i < N; i++ {
		a.polys[0].Coeffs[i] = uint32(i % int(Q))
		a.polys[1].Coeffs[i] = uint32((2*i + 1) % int(Q))
		b.polys[0].Coeffs[i] = uint32((3*i + 5) % int(Q))
		b.polys[1].Coeffs[i] = uint32((4*i + 7) % int(Q))
	}

	if err := out.Add(a, b); err != nil {
		t.Fatalf("Add: %v", err)
	}
	recovered := NewVec(2)
	if err := recovered.Sub(out, b); err != nil {
		t.Fatalf("Sub: %v", err)
	}
	for i := 0; i < recovered.Len(); i++ {
		for j := 0; j < N; j++ {
			if ModQ(recovered.polys[i].Coeffs[j]) != ModQ(a.polys[i].Coeffs[j]) {
				t.Fatalf("Add/Sub mismatch at poly=%d coeff=%d", i, j)
			}
		}
	}
}

func TestVecNTTInverse(t *testing.T) {
	v := NewVec(1)
	for i := 0; i < N; i++ {
		v.polys[0].Coeffs[i] = uint32(i % int(Q))
	}
	orig := NewVec(1)
	_ = orig.CopyFrom(v)

	if err := v.NTT(); err != nil {
		t.Fatalf("Vec.NTT: %v", err)
	}
	if err := v.InvNTT(); err != nil {
		t.Fatalf("Vec.InvNTT: %v", err)
	}
	Freeze(v.polys[0])
	for i := 0; i < N; i++ {
		got := ModQ(FromMont(v.polys[0].Coeffs[i]))
		want := ModQ(orig.polys[0].Coeffs[i])
		if got != want {
			t.Fatalf("NTT roundtrip mismatch at coeff %d: got %d want %d", i, got, want)
		}
	}
}

func TestPointwiseAccMontgomeryVec(t *testing.T) {
	a := NewVec(2)
	b := NewVec(2)
	for i := 0; i < 2; i++ {
		for j := 0; j < N; j++ {
			a.polys[i].Coeffs[j] = ToMont(uint32((i + j + 1) % int(Q)))
			b.polys[i].Coeffs[j] = ToMont(uint32((2*i + j + 2) % int(Q)))
		}
	}
	out := &Poly{}
	if err := PointwiseAccMontgomeryVec(out, a, b); err != nil {
		t.Fatalf("PointwiseAccMontgomeryVec: %v", err)
	}

	for j := 0; j < N; j++ {
		expected := int64(0)
		for i := 0; i < 2; i++ {
			valA := int64((i + j + 1) % int(Q))
			valB := int64((2*i + j + 2) % int(Q))
			expected = (expected + valA*valB) % int64(Q)
		}
		if got := ModQ(FromMont(out.Coeffs[j])); got != uint32(expected) {
			t.Fatalf("pointwise acc mismatch at coeff %d: got %d want %d", j, got, expected)
		}
	}
}

func TestInfinityNorm(t *testing.T) {
	v := NewVec(1)
	v.polys[0].Coeffs[0] = ModQ(uint32(Q - 3))
	v.polys[0].Coeffs[1] = 5
	norm := v.InfinityNorm()
	if norm != 5 {
		t.Fatalf("InfinityNorm = %d, want 5", norm)
	}
}

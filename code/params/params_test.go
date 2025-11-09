// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package params

import "testing"

func TestLookup(t *testing.T) {
	tests := []struct {
		typ      Type
		wantName string
		wantPK   int
		wantSig  int
	}{
		{Type44, "ML-DSA-44", 1312, 2420},
		{Type65, "ML-DSA-65", 1952, 3309},
		{Type87, "ML-DSA-87", 2592, 4627},
	}

	for _, tt := range tests {
		t.Run(string(tt.typ), func(t *testing.T) {
			set, err := Lookup(tt.typ)
			if err != nil {
				t.Fatalf("Lookup(%s) error = %v", tt.typ, err)
			}
			if set.Name != tt.wantName {
				t.Fatalf("Name = %q, want %q", set.Name, tt.wantName)
			}
			if set.PKBytes != tt.wantPK {
				t.Fatalf("PKBytes = %d, want %d", set.PKBytes, tt.wantPK)
			}
			if set.SigBytes != tt.wantSig {
				t.Fatalf("SigBytes = %d, want %d", set.SigBytes, tt.wantSig)
			}
			if set.Q != 8380417 || set.N != 256 || set.D != 13 {
				t.Fatalf("unexpected base constants: Q=%d N=%d D=%d", set.Q, set.N, set.D)
			}
		})
	}
}

func TestLookupUnknown(t *testing.T) {
	if _, err := Lookup(Type("ML-DSA-unknown")); err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestTypes(t *testing.T) {
	got := Types()
	if len(got) != 3 {
		t.Fatalf("Types length = %d, want 3", len(got))
	}
	if got[0] != Type44 || got[1] != Type65 || got[2] != Type87 {
		t.Fatalf("Types order mismatch: %v", got)
	}
}

// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package mldsa

import (
	"testing"
)

func TestParams_StandardSets(t *testing.T) {
	tests := []struct {
		name   string
		params *Params
		want   struct {
			pkBytes  int
			sigBytes int
			k        int
			l        int
		}
	}{
		{
			name:   "ML-DSA-44",
			params: ParamsMLDSA44,
			want: struct {
				pkBytes  int
				sigBytes int
				k        int
				l        int
			}{1312, 2420, 4, 4},
		},
		{
			name:   "ML-DSA-65",
			params: ParamsMLDSA65,
			want: struct {
				pkBytes  int
				sigBytes int
				k        int
				l        int
			}{1952, 3309, 6, 5},
		},
		{
			name:   "ML-DSA-87",
			params: ParamsMLDSA87,
			want: struct {
				pkBytes  int
				sigBytes int
				k        int
				l        int
			}{2592, 4627, 8, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.params.PKBytes != tt.want.pkBytes {
				t.Errorf("PKBytes = %d, want %d", tt.params.PKBytes, tt.want.pkBytes)
			}
			if tt.params.SigBytes != tt.want.sigBytes {
				t.Errorf("SigBytes = %d, want %d", tt.params.SigBytes, tt.want.sigBytes)
			}
			if tt.params.K != tt.want.k {
				t.Errorf("K = %d, want %d", tt.params.K, tt.want.k)
			}
			if tt.params.L != tt.want.l {
				t.Errorf("L = %d, want %d", tt.params.L, tt.want.l)
			}
			if tt.params.Name != tt.name {
				t.Errorf("Name = %s, want %s", tt.params.Name, tt.name)
			}
		})
	}
}

func TestParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  *Params
		wantErr bool
	}{
		{
			name:    "valid ML-DSA-44",
			params:  ParamsMLDSA44,
			wantErr: false,
		},
		{
			name:    "valid ML-DSA-65",
			params:  ParamsMLDSA65,
			wantErr: false,
		},
		{
			name:    "valid ML-DSA-87",
			params:  ParamsMLDSA87,
			wantErr: false,
		},
		{
			name:    "nil params",
			params:  nil,
			wantErr: true,
		},
		{
			name: "invalid custom params",
			params: &Params{
				Name:    "Invalid",
				PKBytes: 999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.ValidateParams()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParams_Constants(t *testing.T) {
	// Verify ML-DSA constants
	if q != 8380417 {
		t.Errorf("q = %d, want 8380417", q)
	}
	if n != 256 {
		t.Errorf("n = %d, want 256", n)
	}
	if d != 13 {
		t.Errorf("d = %d, want 13", d)
	}
	if SeedBytes != 32 {
		t.Errorf("SeedBytes = %d, want 32", SeedBytes)
	}
	if CRHBytes != 64 {
		t.Errorf("CRHBytes = %d, want 64", CRHBytes)
	}
}

func TestParams_Gamma2Calculation(t *testing.T) {
	// Verify Gamma2 calculations are correct
	tests := []struct {
		name   string
		params *Params
		want   int
	}{
		{"ML-DSA-44", ParamsMLDSA44, (q - 1) / 88},
		{"ML-DSA-65", ParamsMLDSA65, (q - 1) / 32},
		{"ML-DSA-87", ParamsMLDSA87, (q - 1) / 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.params.Gamma2 != tt.want {
				t.Errorf("Gamma2 = %d, want %d", tt.params.Gamma2, tt.want)
			}
		})
	}
}

// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package pack

import (
	"fmt"
	"testing"

	"github.com/codethor0/dilivet/code/poly"
)

// TestPackUnpackBitsRoundTripProperty tests that PackBits and UnpackBits are exact inverses
// for various bit widths and value ranges.
func TestPackUnpackBitsRoundTripProperty(t *testing.T) {
	bitWidths := []int{1, 2, 4, 8, 10, 12, 16, 20, 24, 32}
	lengths := []int{1, 10, 100, 256, poly.N}

	for _, bits := range bitWidths {
		for _, length := range lengths {
			t.Run(fmt.Sprintf("bits=%d,len=%d", bits, length), func(t *testing.T) {
				// Generate test values within the bit width
				var maxVal uint32
				if bits == 32 {
					maxVal = ^uint32(0) // Max uint32
				} else {
					maxVal = uint32((1 << bits) - 1)
				}
				vals := make([]uint32, length)
				for i := range vals {
					if maxVal == ^uint32(0) {
						vals[i] = uint32(i)
					} else {
						vals[i] = uint32(i) % (maxVal + 1)
					}
				}

				packed, err := PackBits(vals, bits)
				if err != nil {
					t.Fatalf("PackBits failed: %v", err)
				}

				unpacked, err := UnpackBits(packed, bits, length)
				if err != nil {
					t.Fatalf("UnpackBits failed: %v", err)
				}

				if len(unpacked) != len(vals) {
					t.Fatalf("length mismatch: got %d, want %d", len(unpacked), len(vals))
				}

				for i := range vals {
					if vals[i] != unpacked[i] {
						t.Errorf("roundtrip mismatch at index %d: got %d, want %d", i, unpacked[i], vals[i])
					}
				}
			})
		}
	}
}

// TestPackUnpackPolyCoeffsRoundTripProperty tests round-trip for polynomial coefficients.
func TestPackUnpackPolyCoeffsRoundTripProperty(t *testing.T) {
	bitWidths := []int{10, 12, 16, 20}

	for _, bits := range bitWidths {
		t.Run(fmt.Sprintf("bits=%d", bits), func(t *testing.T) {
			var p poly.Poly
			maxVal := uint32((1 << bits) - 1)
			for i := 0; i < poly.N; i++ {
				p.Coeffs[i] = uint32(i) % (maxVal + 1)
			}

			packed, err := PackPolyCoeffs(&p, bits)
			if err != nil {
				t.Fatalf("PackPolyCoeffs failed: %v", err)
			}

			unpacked, err := UnpackPolyCoeffs(packed, bits)
			if err != nil {
				t.Fatalf("UnpackPolyCoeffs failed: %v", err)
			}

			for i := 0; i < poly.N; i++ {
				if p.Coeffs[i] != unpacked.Coeffs[i] {
					t.Errorf("roundtrip mismatch at coeff %d: got %d, want %d", i, unpacked.Coeffs[i], p.Coeffs[i])
				}
			}
		})
	}
}

// TestPackBitsBoundaryConditions tests edge cases for packing.
func TestPackBitsBoundaryConditions(t *testing.T) {
	tests := []struct {
		name    string
		vals    []uint32
		bits    int
		wantErr bool
	}{
		{
			name:    "empty slice",
			vals:    []uint32{},
			bits:    8,
			wantErr: false,
		},
		{
			name:    "all zeros",
			vals:    make([]uint32, 100),
			bits:    8,
			wantErr: false,
		},
		{
			name:    "all max values",
			vals:    []uint32{255, 255, 255},
			bits:    8,
			wantErr: false,
		},
		{
			name:    "overflow",
			vals:    []uint32{256},
			bits:    8,
			wantErr: true,
		},
		{
			name:    "single value",
			vals:    []uint32{42},
			bits:    8,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := PackBits(tt.vals, tt.bits)
			if (err != nil) != tt.wantErr {
				t.Errorf("PackBits() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUnpackBitsBoundaryConditions tests edge cases for unpacking.
func TestUnpackBitsBoundaryConditions(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		bits    int
		length  int
		wantErr bool
	}{
		{
			name:    "empty input",
			data:    []byte{},
			bits:    8,
			length:  0,
			wantErr: false,
		},
		{
			name:    "too short",
			data:    []byte{0x01},
			bits:    12,
			length:  2,
			wantErr: true,
		},
		{
			name:    "exact length",
			data:    []byte{0x01, 0x02},
			bits:    8,
			length:  2,
			wantErr: false,
		},
		{
			name:    "single byte, single value",
			data:    []byte{0xFF},
			bits:    8,
			length:  1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UnpackBits(tt.data, tt.bits, tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackBits() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPackBitsLengthProperty tests that packed length matches expected byte count.
func TestPackBitsLengthProperty(t *testing.T) {
	tests := []struct {
		vals    []uint32
		bits    int
		wantLen int
	}{
		{
			vals:    make([]uint32, 8),
			bits:    8,
			wantLen: 8,
		},
		{
			vals:    make([]uint32, 8),
			bits:    4,
			wantLen: 4,
		},
		{
			vals:    make([]uint32, 8),
			bits:    10,
			wantLen: 10, // 8 * 10 = 80 bits = 10 bytes
		},
		{
			vals:    make([]uint32, 7),
			bits:    10,
			wantLen: 9, // 7 * 10 = 70 bits = 9 bytes (rounded up)
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("len=%d,bits=%d", len(tt.vals), tt.bits), func(t *testing.T) {
			packed, err := PackBits(tt.vals, tt.bits)
			if err != nil {
				t.Fatalf("PackBits failed: %v", err)
			}
			if len(packed) != tt.wantLen {
				t.Errorf("packed length = %d, want %d", len(packed), tt.wantLen)
			}
		})
	}
}


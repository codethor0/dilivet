// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package pack

import (
	"testing"

	"github.com/codethor0/dilivet/code/poly"
)

func TestPackUnpackBitsRoundTrip(t *testing.T) {
	input := []uint32{1, 2, 3, 4, 5, 1023}
	bits := 10

	packed, err := PackBits(input, bits)
	if err != nil {
		t.Fatalf("PackBits error: %v", err)
	}
	unpacked, err := UnpackBits(packed, bits, len(input))
	if err != nil {
		t.Fatalf("UnpackBits error: %v", err)
	}

	for i := range input {
		if input[i] != unpacked[i] {
			t.Fatalf("roundtrip mismatch at index %d: got %d want %d", i, unpacked[i], input[i])
		}
	}
}

func TestPackBitsOverflow(t *testing.T) {
	_, err := PackBits([]uint32{8}, 3)
	if err == nil {
		t.Fatal("expected overflow error")
	}
}

func TestUnpackBitsShortInput(t *testing.T) {
	_, err := UnpackBits([]byte{0x01}, 12, 2)
	if err == nil {
		t.Fatal("expected ErrInvalidLength")
	}
}

func TestPackPolyCoeffsRoundTrip(t *testing.T) {
	var p poly.Poly
	for i := 0; i < poly.N; i++ {
		p.Coeffs[i] = uint32(i % 1024)
	}

	bits := 10
	packed, err := PackPolyCoeffs(&p, bits)
	if err != nil {
		t.Fatalf("PackPolyCoeffs error: %v", err)
	}

	unpacked, err := UnpackPolyCoeffs(packed, bits)
	if err != nil {
		t.Fatalf("UnpackPolyCoeffs error: %v", err)
	}

	for i := 0; i < poly.N; i++ {
		if p.Coeffs[i] != unpacked.Coeffs[i] {
			t.Fatalf("roundtrip mismatch at coeff %d: got %d want %d", i, unpacked.Coeffs[i], p.Coeffs[i])
		}
	}
}

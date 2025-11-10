// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package pack

import (
	"errors"
	"fmt"

	"github.com/codethor0/dilivet/code/poly"
)

var (
	// ErrInvalidBits is returned when a bit-width is out of the supported range.
	ErrInvalidBits = errors.New("pack: invalid bit width")
	// ErrOverflow is returned when a value does not fit within the requested width.
	ErrOverflow = errors.New("pack: value does not fit in requested bit width")
	// ErrInvalidLength is returned when the supplied byte slice is too short.
	ErrInvalidLength = errors.New("pack: input too short for requested unpack length")
)

// ErrInvalidHint indicates an invalid hint encoding.
var ErrInvalidHint = errors.New("pack: invalid hint encoding")

// mask returns a uint64 mask with the lowest bits bits set.
func mask(bits int) uint64 {
	if bits == 64 {
		return ^uint64(0)
	}
	return (uint64(1) << bits) - 1
}

// PackBits packs the provided coefficients into a byte slice using the supplied bit width.
// Values must already be reduced to fit within the provided width.
func PackBits(vals []uint32, bits int) ([]byte, error) {
	if bits <= 0 || bits > 32 {
		return nil, ErrInvalidBits
	}
	if len(vals) == 0 {
		return []byte{}, nil
	}

	totalBits := len(vals) * bits
	outLen := (totalBits + 7) / 8
	out := make([]byte, outLen)

	var acc uint64
	var accBits uint
	var idx int
	limit := uint64(1)<<bits - 1

	for _, v := range vals {
		if uint64(v) > limit {
			return nil, fmt.Errorf("%w: value %d exceeds %d bits", ErrOverflow, v, bits)
		}
		acc |= uint64(v) << accBits
		accBits += uint(bits)

		for accBits >= 8 {
			if idx >= len(out) {
				panic("pack: internal error overflow") // should never happen
			}
			out[idx] = byte(acc & 0xff)
			idx++
			acc >>= 8
			accBits -= 8
		}
	}

	if accBits > 0 {
		if idx >= len(out) {
			panic("pack: internal error final byte")
		}
		out[idx] = byte(acc & 0xff)
	}

	return out, nil
}

// UnpackBits unpacks exactly count values from data, each encoded with the supplied bit width.
func UnpackBits(data []byte, bits int, count int) ([]uint32, error) {
	if bits <= 0 || bits > 32 {
		return nil, ErrInvalidBits
	}
	if count < 0 {
		return nil, errors.New("pack: negative count")
	}
	if count == 0 {
		return []uint32{}, nil
	}
	result := make([]uint32, count)

	var acc uint64
	var accBits uint
	var idx int
	mask := mask(bits)

	for i := 0; i < count; i++ {
		for accBits < uint(bits) {
			if idx >= len(data) {
				return nil, ErrInvalidLength
			}
			acc |= uint64(data[idx]) << accBits
			accBits += 8
			idx++
		}

		val := uint32(acc & mask)
		result[i] = val
		acc >>= uint(bits)
		accBits -= uint(bits)
	}

	return result, nil
}

// PackPolyCoeffs encodes the polynomial coefficients into a bit-packed byte slice.
func PackPolyCoeffs(p *poly.Poly, bits int) ([]byte, error) {
	if p == nil {
		return nil, errors.New("pack: nil polynomial")
	}
	return PackBits(p.Coeffs[:], bits)
}

// UnpackPolyCoeffs decodes a bit-packed buffer into a polynomial, using the supplied bit width.
func UnpackPolyCoeffs(data []byte, bits int) (*poly.Poly, error) {
	vals, err := UnpackBits(data, bits, poly.N)
	if err != nil {
		return nil, err
	}
	var p poly.Poly
	for i := 0; i < poly.N; i++ {
		p.Coeffs[i] = vals[i]
	}
	return &p, nil
}

// UnpackPolyLeGamma1 decodes a polynomial from the two's-complement centered representation
// used for z coefficients in ML-DSA. Each coefficient is stored as a signed
// little-endian integer with b bits (typically 18 or 20 depending on parameter set).
func UnpackPolyLeGamma1(data []byte, bits int) (*poly.Poly, error) {
	if bits <= 0 || bits > 24 {
		return nil, ErrInvalidBits
	}
	totalBits := bits * poly.N
	expectedLen := (totalBits + 7) / 8
	if len(data) < expectedLen {
		return nil, ErrInvalidLength
	}

	var p poly.Poly
	bitMask := uint32(1<<bits) - 1
	byteIdx := 0
	bitIdx := 0
	for i := 0; i < poly.N; i++ {
		var coeff uint32
		remaining := bits
		shift := 0
		for remaining > 0 {
			if byteIdx >= len(data) {
				return nil, ErrInvalidLength
			}
			available := 8 - bitIdx
			take := remaining
			if take > available {
				take = available
			}
			chunk := (uint32(data[byteIdx]) >> bitIdx) & ((1 << take) - 1)
			coeff |= chunk << shift
			remaining -= take
			shift += take
			bitIdx += take
			if bitIdx == 8 {
				bitIdx = 0
				byteIdx++
			}
		}
		coeff &= bitMask
		// Convert from two's complement centered form into canonical representative.
		signBit := uint32(1) << (bits - 1)
		if coeff&signBit != 0 {
			coeff = coeff - (1 << bits)
		}
		if coeff >= poly.Q {
			coeff = coeff % poly.Q
		} else {
			coeff = (coeff + poly.Q) % poly.Q
		}
		p.Coeffs[i] = coeff
	}
	return &p, nil
}

// UnpackHint decodes the hint vector h from the compressed representation used in ML-DSA.
// The encoding consists of omega indices followed by padding with 0xFF.
func UnpackHint(data []byte, omega int) ([]uint8, error) {
	if omega < 0 || omega > poly.N {
		return nil, ErrInvalidHint
	}
	h := make([]uint8, 0, omega)
	for i, b := range data {
		if b == 0xFF {
			// padding byte, all subsequent bytes must also be 0xFF
			for _, rem := range data[i:] {
				if rem != 0xFF {
					return nil, ErrInvalidHint
				}
			}
			break
		}
		if int(b) >= poly.N {
			return nil, ErrInvalidHint
		}
		h = append(h, b)
		if len(h) > omega {
			return nil, ErrInvalidHint
		}
	}
	return h, nil
}

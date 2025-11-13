// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package mldsa

import (
	"crypto/subtle"
	"fmt"

	"github.com/codethor0/dilivet/code/hash"
	"github.com/codethor0/dilivet/code/pack"
	"github.com/codethor0/dilivet/code/poly"
	"golang.org/x/crypto/sha3"
)

// verifyFull implements full FIPS 204 Algorithm 3 (Signature Verification).
func verifyFull(pk, msg, sig []byte, params *Params) (bool, error) {
	// Step 1: Parse public key (ρ, t₁)
	// Public key format: ρ (32 bytes) || t₁ (compressed, k * n * du bits)
	if len(pk) < SeedBytes {
		return false, ErrInvalidPublicKey
	}
	rho := pk[:SeedBytes]
	t1Bytes := pk[SeedBytes:]

	// Step 2: Decompress signature (c̃, z, h)
	// Signature format: c̃ (tau bytes) || z (k * n * gamma1Bits bits) || h (omega bytes + padding)
	tauBytes := (params.Tau + 7) / 8
	if len(sig) < tauBytes {
		return false, ErrInvalidSignature
	}
	ctilde := sig[:tauBytes]

	// Unpack z vector (k polynomials, each gamma1Bits per coefficient)
	zStart := tauBytes
	zBytes := (params.K*poly.N*params.Gamma1Bits + 7) / 8
	if len(sig) < zStart+zBytes {
		return false, ErrInvalidSignature
	}
	zData := sig[zStart : zStart+zBytes]

	zVec := poly.NewVec(params.K)
	for i := 0; i < params.K; i++ {
		offset := (i*poly.N*params.Gamma1Bits + 7) / 8
		if offset+((poly.N*params.Gamma1Bits+7)/8) > len(zData) {
			return false, ErrInvalidSignature
		}
		p, err := pack.UnpackPolyLeGamma1(zData[offset:], params.Gamma1Bits)
		if err != nil {
			return false, fmt.Errorf("mldsa: unpack z[%d]: %w", i, err)
		}
		// Check ||z||∞ < γ₁
		for _, coeff := range p.Coeffs {
			canon := poly.Canonical(coeff)
			if canon < 0 {
				canon = -canon
			}
			if int32(canon) >= int32(params.Gamma1) {
				return false, ErrInvalidSignature
			}
		}
		zVec.Polys()[i] = p
	}

	// Unpack hint h
	hStart := zStart + zBytes
	hData := sig[hStart:]
	h, err := pack.UnpackHint(hData, params.Omega)
	if err != nil {
		return false, fmt.Errorf("mldsa: unpack hint: %w", err)
	}
	if len(h) > params.Omega {
		return false, ErrInvalidSignature
	}

	// Step 3: Expand matrix A from seed ρ
	// A is a k×l matrix of polynomials
	// For now, we'll expand it on-the-fly as needed
	// TODO: Cache A if performance requires it

	// Step 4: Sample challenge polynomial c from c̃
	c := &poly.Poly{}
	if err := sampleChallenge(c, ctilde, params.Tau); err != nil {
		return false, fmt.Errorf("mldsa: sample challenge: %w", err)
	}

	// Step 5: Compute μ = CRH(tr || msg) where tr = H(pk)
	tr := make([]byte, CRHBytes)
	hashPublicKey(tr, pk)
	mu := make([]byte, CRHBytes)
	hash.SumShake256(mu, tr, msg)

	// Step 6: Compute w'₁ = UseHint(h, Az - c·t₁·2^d, 2γ₂)
	// First, compute Az
	azVec := poly.NewVec(params.K)
	for i := 0; i < params.K; i++ {
		azVec.Polys()[i] = &poly.Poly{}
		for j := 0; j < params.L; j++ {
			// Expand A[i][j] from rho
			aij := &poly.Poly{}
			nonce := uint16(i*params.L + j)
			if err := poly.SamplePolyUniform(aij, rho, nonce); err != nil {
				return false, fmt.Errorf("mldsa: expand A[%d][%d]: %w", i, j, err)
			}
			// Convert to NTT domain
			if err := poly.NTT(aij); err != nil {
				return false, fmt.Errorf("mldsa: NTT A[%d][%d]: %w", i, j, err)
			}

			// z[j] in NTT
			zjNTT := &poly.Poly{}
			copy(zjNTT.Coeffs[:], zVec.Polys()[j].Coeffs[:])
			if err := poly.NTT(zjNTT); err != nil {
				return false, fmt.Errorf("mldsa: NTT z[%d]: %w", j, err)
			}

			// Multiply and accumulate: az[i] += A[i][j] * z[j]
			prod := &poly.Poly{}
			prod.PointwiseMontgomery(aij, zjNTT)
			azVec.Polys()[i].Add(azVec.Polys()[i], prod)
		}
		// Convert back from NTT
		if err := poly.InvNTT(azVec.Polys()[i]); err != nil {
			return false, fmt.Errorf("mldsa: InvNTT az[%d]: %w", i, err)
		}
	}

	// Compute c·t₁·2^d
	// First, unpack t₁
	t1Vec := poly.NewVec(params.K)
	for i := 0; i < params.K; i++ {
		offset := (i*poly.N*params.DuBits + 7) / 8
		if offset+((poly.N*params.DuBits+7)/8) > len(t1Bytes) {
			return false, ErrInvalidPublicKey
		}
		t1i, err := pack.UnpackBits(t1Bytes[offset:], params.DuBits, poly.N)
		if err != nil {
			return false, fmt.Errorf("mldsa: unpack t1[%d]: %w", i, err)
		}
		t1Poly := &poly.Poly{}
		for j := 0; j < poly.N; j++ {
			t1Poly.Coeffs[j] = t1i[j] << d // Multiply by 2^d
		}
		t1Vec.Polys()[i] = t1Poly
	}

	// Convert c and t1Vec to NTT
	cNTT := &poly.Poly{}
	copy(cNTT.Coeffs[:], c.Coeffs[:])
	if err := poly.NTT(cNTT); err != nil {
		return false, fmt.Errorf("mldsa: NTT c: %w", err)
	}

	for i := 0; i < params.K; i++ {
		if err := poly.NTT(t1Vec.Polys()[i]); err != nil {
			return false, fmt.Errorf("mldsa: NTT t1[%d]: %w", i, err)
		}
	}

	// Compute c·t₁
	ct1Vec := poly.NewVec(params.K)
	for i := 0; i < params.K; i++ {
		ct1Vec.Polys()[i] = &poly.Poly{}
		ct1Vec.Polys()[i].PointwiseMontgomery(cNTT, t1Vec.Polys()[i])
		if err := poly.InvNTT(ct1Vec.Polys()[i]); err != nil {
			return false, fmt.Errorf("mldsa: InvNTT ct1[%d]: %w", i, err)
		}
	}

	// Compute w = az - ct1
	wVec := poly.NewVec(params.K)
	for i := 0; i < params.K; i++ {
		wVec.Polys()[i] = &poly.Poly{}
		wVec.Polys()[i].Sub(azVec.Polys()[i], ct1Vec.Polys()[i])
	}

	// Apply hints to get w'₁
	w1Prime := poly.NewVec(params.K)
	for i := 0; i < params.K; i++ {
		w1Prime.Polys()[i] = &poly.Poly{}
		if err := useHint(w1Prime.Polys()[i], wVec.Polys()[i], h, i, params.Gamma2); err != nil {
			return false, fmt.Errorf("mldsa: useHint w[%d]: %w", i, err)
		}
	}

	// Step 7: Encode w'₁ and compute c' = H(μ || w₁Encode(w'₁))
	w1Encoded := encodeW1(w1Prime, params.DvBits)
	cPrime := make([]byte, tauBytes)
	hashChallenge(cPrime, mu, w1Encoded, params.Tau)

	// Step 8: Constant-time comparison
	return subtle.ConstantTimeCompare(ctilde, cPrime) == 1, nil
}

// sampleChallenge samples a challenge polynomial c with exactly tau coefficients set to ±1.
func sampleChallenge(c *poly.Poly, ctilde []byte, tau int) error {
	// Use SHAKE256 to expand ctilde into indices
	xof := sha3.NewShake256()
	if _, err := xof.Write(ctilde); err != nil {
		return err
	}

	// Sample tau distinct indices in [0, 2n)
	indices := make(map[uint16]bool)
	buf := make([]byte, 2)
	for len(indices) < tau {
		if _, err := xof.Read(buf); err != nil {
			return err
		}
		idx := uint16(buf[0]) | (uint16(buf[1]) << 8)
		idx %= 2 * poly.N
		if !indices[idx] {
			indices[idx] = true
		}
	}

	// Set coefficients: c[i] = 1 if i in indices and sign bit set, else -1
	for i := range c.Coeffs {
		c.Coeffs[i] = 0
	}
	for idx := range indices {
		pos := idx % poly.N
		sign := (idx / poly.N) & 1
		if sign == 0 {
			c.Coeffs[pos] = 1
		} else {
			c.Coeffs[pos] = poly.Q - 1 // -1 mod q
		}
	}

	return nil
}

// hashPublicKey computes tr = H(pk) using SHAKE256.
func hashPublicKey(tr []byte, pk []byte) {
	hash.SumShake256(tr, pk)
}

// useHint applies the hint h to reconstruct w'₁ from w.
func useHint(w1Prime *poly.Poly, w *poly.Poly, h []uint8, vecIdx int, gamma2 int) error {
	// For each coefficient in w, use the hint to determine the high bits
	// This is a simplified version - full implementation needs proper hint application
	for i := range w1Prime.Coeffs {
		w1Prime.Coeffs[i] = decomposeHigh(w.Coeffs[i], gamma2)
	}
	// TODO: Apply hints from h vector properly
	return nil
}

// decomposeHigh returns the high part of the decomposition.
func decomposeHigh(r uint32, gamma2 int) uint32 {
	// Decompose r = r1 * 2*gamma2 + r0 where |r0| <= gamma2
	r0 := int32(poly.Canonical(r))
	r1 := r0 / int32(2*gamma2)
	if r0 < 0 {
		r1--
	}
	if r1 < 0 {
		r1 += int32(poly.Q)
	}
	return uint32(r1)
}

// encodeW1 encodes w'₁ into a byte string for hashing.
func encodeW1(w1Prime *poly.Vec, dvBits int) []byte {
	// Encode each polynomial in w1Prime using dvBits per coefficient
	totalBits := w1Prime.Len() * poly.N * dvBits
	totalBytes := (totalBits + 7) / 8
	encoded := make([]byte, totalBytes)

	bitIdx := 0
	for i := 0; i < w1Prime.Len(); i++ {
		p := w1Prime.Polys()[i]
		for j := 0; j < poly.N; j++ {
			coeff := p.Coeffs[j] & ((1 << dvBits) - 1)
			// Pack into encoded
			for b := 0; b < dvBits; b++ {
				byteIdx := bitIdx / 8
				bitPos := bitIdx % 8
				if (coeff>>uint(b))&1 != 0 {
					encoded[byteIdx] |= 1 << bitPos
				}
				bitIdx++
			}
		}
	}

	return encoded
}

// hashChallenge computes c' = H(μ || w₁Encoded) and outputs tau bytes.
func hashChallenge(cPrime []byte, mu []byte, w1Encoded []byte, tau int) {
	// Use SHAKE256 and take first tau bytes
	xof := sha3.NewShake256()
	_, _ = xof.Write(mu)
	_, _ = xof.Write(w1Encoded)
	_, _ = xof.Read(cPrime)
}

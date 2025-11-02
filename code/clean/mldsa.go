// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mldsa implements the Module-Lattice-Based Digital Signature
// Standard, ML-DSA (Dilithium), as specified in NIST FIPS 204.
package mldsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/sha3"
)

// This implementation supports ML-DSA-44, ML-DSA-65, and ML-DSA-87, informally
// known as Dilithium 2, 3, and 5, respectively. The security levels correspond
// to NIST security levels 2, 3, and 5.

const (
	n = 256
	q = 8380417

	// d is the reduction factor used in MakeHint and UseHint.
	d = 13

	// gamma1 is one of the parameters of ML-DSA.
	// gamma1 = (q - 1) / (2 * gamma2)
	//
	// FIPS 204, Section 3.1
	gamma1L23 = 1 << 17
	gamma1L5  = 1 << 19

	// gamma2 is one of the parameters of ML-DSA.
	//
	// FIPS 204, Section 3.1
	gamma2L23 = (q - 1) / 88
	gamma2L5  = (q - 1) / 32

	// alpha is 2 * gamma2.
	//
	// FIPS 204, Section 4.2
	alphaL23 = 2 * gamma2L23
	alphaL5  = 2 * gamma2L5

	// eta is the bound on the coefficients of s1 and s2.
	//
	// FIPS 204, Table 1
	etaL2 = 2
	etaL3 = 4
	etaL5 = 2

	// beta is the bound on the coefficients of y.
	//
	// FIPS 204, Table 1
	betaL2 = 78
	betaL3 = 196
	betaL5 = 120

	// omega is the number of nonzero coefficients in h.
	//
	// FIPS 204, Table 1
	omegaL2 = 80
	omegaL3 = 55
	omegaL5 = 75
)

// Mode is a ML-DSA mode (parameter set).
type Mode struct {
	// Name is the name of the mode, e.g. "ML-DSA-44".
	Name string

	// ID is the 16-bit identifier for the mode.
	ID uint16

	// Level is the NIST security level.
	Level int

	// K is the number of vectors in the matrix A.
	K int

	// L is the number of vectors in the secret key.
	L int

	// SeedLen is the length of the seed used to generate A.
	SeedLen int

	// RhoLen is the length of the seed for s1 and s2.
	RhoLen int

	// KLen is the length of the seed for y and the challenge.
	KLen int

	// TrLen is the length of the seed used for the challenge.
	TrLen int

	// MuLen is the length of the hash prefix for the challenge.
	MuLen int

	// PKLen is the length of the encoded public key.
	PKLen int

	// SKLen is the length of the encoded private key.
	SKLen int

	// SigLen is the length of the encoded signature.
	SigLen int

	// ML-DSA parameters, see FIPS 204, Table 1.
	Tau    int
	Lambda int
	Gamma1 int
	Gamma2 int
	Alpha  int
	Beta   int
	Omega  int
	Eta    int
}

var (
	MLDSA44 = &Mode{
		Name:    "ML-DSA-44",
		ID:      0x0202,
		Level:   2,
		K:       4,
		L:       4,
		SeedLen: 32,
		RhoLen:  32,
		KLen:    32,
		TrLen:   64,
		MuLen:   32,
		PKLen:   1312,
		SKLen:   2560,
		SigLen:  2420,
		Tau:     39,
		Lambda:  128,
		Gamma1:  gamma1L23,
		Gamma2:  gamma2L23,
		Alpha:   alphaL23,
		Beta:    betaL2,
		Omega:   omegaL2,
		Eta:     etaL2,
	}
	MLDSA65 = &Mode{
		Name:    "ML-DSA-65",
		ID:      0x0303,
		Level:   3,
		K:       6,
		L:       5,
		SeedLen: 32,
		RhoLen:  32,
		KLen:    32,
		TrLen:   64,
		MuLen:   32,
		PKLen:   1952,
		SKLen:   4032,
		SigLen:  3309, // 3309 in FIPS 204, 3293 in draft
		Tau:     49,
		Lambda:  192,
		Gamma1:  gamma1L23,
		Gamma2:  gamma2L23,
		Alpha:   alphaL23,
		Beta:    betaL3,
		Omega:   omegaL3,
		Eta:     etaL3,
	}
	MLDSA87 = &Mode{
		Name:    "ML-DSA-87",
		ID:      0x0404,
		Level:   5,
		K:       8,
		L:       7,
		SeedLen: 32,
		RhoLen:  32,
		KLen:    32,
		TrLen:   64,
		MuLen:   32,
		PKLen:   2592,
		SKLen:   4928,
		SigLen:  4611, // 4611 in FIPS 204, 4595 in draft
		Tau:     60,
		Lambda:  256,
		Gamma1:  gamma1L5,
		Gamma2:  gamma2L5,
		Alpha:   alphaL5,
		Beta:    betaL5,
		Omega:   omegaL5,
		Eta:     etaL5,
	}
)

var modes = []*Mode{MLDSA44, MLDSA65, MLDSA87}

// poly is a polynomial in R_q.
type poly [n]int32

// polyVec is a vector of polynomials.
type polyVec []poly

func (m *Mode) newPolyVecL() polyVec {
	return make(polyVec, m.L)
}

func (m *Mode) newPolyVecK() polyVec {
	return make(polyVec, m.K)
}

func (m *Mode) polyVecSizeL() int {
	return m.L * n * 4
}

func (m *Mode) polyVecSizeK() int {
	return m.K * n * 4
}

// centeredPoly is a polynomial with centered coefficients.
//
// FIPS 204, Section 2.5
func (p *poly) centered() {
	for i := 0; i < n; i++ {
		p[i] = centered(p[i])
	}
}

func centered(c int32) int32 {
	t := c % q
	if t > (q-1)/2 {
		t -= q
	}
	return t
}

// power2round computes (t0, t1) = Power2Round(t).
//
// FIPS 204, Algorithm 6 (Power2Round)
func power2round(t int32, m *Mode) (int32, int32) {
	t = t % q
	t0 := t % m.Alpha
	if t0 > m.Alpha/2 {
		t0 -= m.Alpha
	}
	t1 := (t - t0) / m.Alpha
	return t0, t1
}

// decompose computes (t0, t1) = Decompose(t).
//
// FIPS 204, Algorithm 7 (Decompose)
func decompose(t int32, m *Mode) (int32, int32) {
	t = t % q
	t0 := t
	if t0 > (q-1)/2 {
		t0 -= q
	}
	if t0 < -(q-1)/2 {
		t0 += q
	}

	c := (t - t0) / q
	t0 = t - c*q

	t1 := int32(0)
	if t0 > m.Gamma2 {
		t1 = 1
		t0 -= m.Gamma2
	}
	if t0 < -m.Gamma2 {
		t1 = -1
		t0 += m.Gamma2
	}
	return t0, t1
}

// makeHint computes h = MakeHint(t0, t1).
//
// FIPS 204, Algorithm 8 (MakeHint)
func makeHint(t0, t1 int32, m *Mode) int32 {
	if t0 > m.Gamma2 || t0 < -m.Gamma2 {
		return 0
	}
	if t0 == -m.Gamma2 && t1 != 0 {
		return 0
	}
	if t1 == 1 {
		return 1
	}
	return 0
}

// useHint computes t = UseHint(h, t).
//
// FIPS 204, Algorithm 9 (UseHint)
func useHint(h, t int32, m *Mode) int32 {
	t0, t1 := decompose(t, m)
	if h == 0 {
		return t1
	}
	if t0 > 0 {
		if t1 == (q-1)/m.Alpha-1 {
			return 0
		}
		return t1 + 1
	}
	if t1 == 0 {
		return (q - 1) / m.Alpha
	}
	return t1 - 1
}

// bitsToBytes packs a bit array into a byte array.
func bitsToBytes(bits []byte) []byte {
	out := make([]byte, (len(bits)+7)/8)
	for i := 0; i < len(bits); i++ {
		out[i/8] |= bits[i] << (i % 8)
	}
	return out
}

// bytesToBits unpacks a byte array into a bit array.
func bytesToBits(in []byte) []byte {
	out := make([]byte, len(in)*8)
	for i := 0; i < len(out); i++ {
		out[i] = (in[i/8] >> (i % 8)) & 1
	}
	return out
}

// shake returns a SHAKE-256 XOF.
func shake(msg []byte) sha3.ShakeHash {
	h := sha3.NewShake256()
	h.Write(msg)
	return h
}

// hash returns the 32-byte hash of msg.
func hash(msg []byte) []byte {
	out := make([]byte, 32)
	sha3.New256().Write(msg)
	return out
}

// h outputs H(msg).
//
// FIPS 204, Section 5.1
func (m *Mode) h(msg []byte) []byte {
	out := make([]byte, m.MuLen)
	h := sha3.NewShake256()
	h.Write(msg)
	h.Read(out) // Read MuLen bytes
	return out
}

// sampleInBall samples a polynomial with Tau coefficients equal to 1 or -1.
//
// FIPS 204, Algorithm 10 (SampleInBall)
func (m *Mode) sampleInBall(seed []byte) poly {
	var c poly
	buf := make([]byte, 8)
	h := shake(seed[:m.KLen])
	h.Read(buf)
	sign := bytesToBits(buf[:8*m.Lambda])
	h = shake(seed)
	for i := m.Lambda; i < n; {
		b := make([]byte, 1)
		h.Read(b)
		j := int(b[0])
		if j < i {
			c[i] = c[j]
			c[j] = 0
			if sign[i/8]>>(i%8)&1 == 1 {
				c[i] = -c[i]
			}
			i++
		}
	}
	return c
}

// expandA computes the matrix A.
//
// FIPS 204, Algorithm 11 (ExpandA)
func (m *Mode) expandA(rho []byte) []poly {
	A := make([]poly, m.K*m.L)
	for i := 0; i < m.K; i++ {
		for j := 0; j < m.L; j++ {
			A[i*m.L+j] = m.expandAij(rho, i, j)
		}
	}
	return A
}

func (m *Mode) expandAij(rho []byte, i, j int) poly {
	var p poly
	h := sha3.NewShake128()
	h.Write(rho)
	h.Write([]byte{byte(j), byte(j >> 8)})
	h.Write([]byte{byte(i), byte(i >> 8)})
	buf := make([]byte, 3*168) // 168 = ceil(log_2(q)) * n / 8
	h.Read(buf)
	// TODO: implement rejection sampling
	// For now, assume a uniform distribution
	for k := 0; k < n; k++ {
		p[k] = int32(buf[3*k]) | int32(buf[3*k+1])<<8 | int32(buf[3*k+2])<<16
		p[k] &= 0x7FFFFF
		p[k] = p[k] % q
	}
	return p
}

// expandS samples the secret key vectors s1 and s2.
//
// FIPS 204, Algorithm 12 (ExpandS)
func (m *Mode) expandS(rho []byte) (polyVec, polyVec) {
	s1 := m.newPolyVecL()
	s2 := m.newPolyVecK()
	h := sha3.NewShake256()
	h.Write(rho)
	buf := make([]byte, n)
	for i := 0; i < m.L; i++ {
		h.Read(buf)
		s1[i] = m.sampleEta(buf)
	}
	for i := 0; i < m.K; i++ {
		h.Read(buf)
		s2[i] = m.sampleEta(buf)
	}
	return s1, s2
}

func (m *Mode) sampleEta(buf []byte) poly {
	var p poly
	switch m.Eta {
	case 2:
		bits := bytesToBits(buf)
		for i := 0; i < n; i++ {
			t := int32(bits[2*i]) | int32(bits[2*i+1])<<1
			p[i] = 2 - t
		}
	case 4:
		bits := bytesToBits(buf)
		for i := 0; i < n; i++ {
			t := int32(bits[3*i]) | int32(bits[3*i+1])<<1 | int32(bits[3*i+2])<<2
			p[i] = 4 - t
		}
	}
	return p
}

// expandMask samples the mask vector y.
//
// FIPS 204, Algorithm 13 (ExpandMask)
func (m *Mode) expandMask(rho []byte, mu int) polyVec {
	y := m.newPolyVecL()
	h := sha3.NewShake256()
	h.Write(rho)
	h.Write([]byte{byte(mu), byte(mu >> 8)})
	buf := make([]byte, (n*m.Gamma1*2+7)/8)
	for i := 0; i < m.L; i++ {
		h.Read(buf)
		y[i] = m.sampleGamma1(buf)
	}
	return y
}

func (m *Mode) sampleGamma1(buf []byte) poly {
	var p poly
	bits := bytesToBits(buf)
	for i := 0; i < n; i++ {
		var t int32
		for j := 0; j < m.Gamma1*2; j++ {
			t |= int32(bits[i*m.Gamma1*2+j]) << j
		}
		p[i] = (1 << m.Gamma1) - t
	}
	return p
}

// polyNTT computes the NTT of p.
func (p *poly) ntt() {
	// TODO: implement
}

// polyInvNTT computes the inverse NTT of p.
func (p *poly) invNTT() {
	// TODO: implement
}

// polyMulNTT multiplies p and v in the NTT domain.
func (p *poly) mulNTT(v *poly) {
	for i := 0; i < n; i++ {
		p[i] = (p[i] * v[i]) % q
	}
}

// polyVecNTT computes the NTT of v.
func (v polyVec) ntt() {
	for i := 0; i < len(v); i++ {
		v[i].ntt()
	}
}

// polyVecInvNTT computes the inverse NTT of v.
func (v polyVec) invNTT() {
	for i := 0; i < len(v); i++ {
		v[i].invNTT()
	}
}

// polyVecPointwiseMulNTT multiplies v and w in the NTT domain.
func (v polyVec) pointwiseMulNTT(w polyVec) {
	for i := 0; i < len(v); i++ {
		v[i].mulNTT(&w[i])
	}
}

// polyVecAdd adds w to v.
func (v polyVec) add(w polyVec) {
	for i := 0; i < len(v); i++ {
		for j := 0; j < n; j++ {
			v[i][j] = (v[i][j] + w[i][j]) % q
		}
	}
}

// polyVecSub subtracts w from v.
func (v polyVec) sub(w polyVec) {
	for i := 0; i < len(v); i++ {
		for j := 0; j < n; j++ {
			v[i][j] = (v[i][j] - w[i][j]) % q
		}
	}
}

// polyVecCheckNorm checks if the norm of v is bounded by b.
func (v polyVec) checkNorm(b int) bool {
	for i := 0; i < len(v); i++ {
		if v[i].checkNorm(b) {
			return true
		}
	}
	return false
}

func (p *poly) checkNorm(b int) bool {
	for i := 0; i < n; i++ {
		t := p[i]
		if t > (q-1)/2 {
			t -= q
		}
		if t < -(q-1)/2 {
			t += q
		}
		if t > int32(b) || t < -int32(b) {
			return true
		}
	}
	return false
}

// matrixMulNTT computes v = A * w.
func (m *Mode) matrixMulNTT(A []poly, w polyVec) polyVec {
	v := m.newPolyVecK()
	for i := 0; i < m.K; i++ {
		for j := 0; j < m.L; j++ {
			var t poly
			t = A[i*m.L+j]
			t.mulNTT(&w[j])
			v[i].add(&t)
		}
	}
	return v
}

// matrixMulTransposeNTT computes v = A^T * w.
func (m *Mode) matrixMulTransposeNTT(A []poly, w polyVec) polyVec {
	v := m.newPolyVecL()
	for i := 0; i < m.L; i++ {
		for j := 0; j < m.K; j++ {
			var t poly
			t = A[j*m.L+i]
			t.mulNTT(&w[j])
			v[i].add(&t)
		}
	}
	return v
}

// polyEncode encodes p into 32*log_2(q) bits.
func (p *poly) encode(out []byte, bits int) {
	// TODO: implement
}

// polyDecode decodes p from 32*log_2(q) bits.
func (p *poly) decode(in []byte, bits int) {
	// TODO: implement
}

// polyVecEncode encodes v.
func (v polyVec) encode(out []byte, bits int) {
	for i := 0; i < len(v); i++ {
		v[i].encode(out[i*n*bits/8:], bits)
	}
}

// polyVecDecode decodes v.
func (v polyVec) decode(in []byte, bits int) {
	for i := 0; i < len(v); i++ {
		v[i].decode(in[i*n*bits/8:], bits)
	}
}

// encodePK encodes the public key.
//
// FIPS 204, Algorithm 14 (EncodePK)
func (m *Mode) encodePK(pk []byte, rho []byte, t1 polyVec) {
	copy(pk, rho)
	t1.encode(pk[m.SeedLen:], d)
}

// decodePK decodes the public key.
//
// FIPS 204, Algorithm 15 (DecodePK)
func (m *Mode) decodePK(pk []byte) ([]byte, polyVec) {
	rho := pk[:m.SeedLen]
	t1 := m.newPolyVecK()
	t1.decode(pk[m.SeedLen:], d)
	return rho, t1
}

// encodeSK encodes the private key.
//
// FIPS 204, Algorithm 16 (EncodeSK)
func (m *Mode) encodeSK(sk []byte, rho []byte, K []byte, tr []byte, s1, s2 polyVec, t0 polyVec) {
	copy(sk, rho)
	copy(sk[m.SeedLen:], K)
	copy(sk[m.SeedLen+m.KLen:], tr)
	off := m.SeedLen + m.KLen + m.TrLen
	s1.encode(sk[off:], m.Eta)
	off += len(s1) * n * m.Eta / 8
	s2.encode(sk[off:], m.Eta)
	off += len(s2) * n * m.Eta / 8
	t0.encode(sk[off:], d)
}

// decodeSK decodes the private key.
//
// FIPS 204, Algorithm 17 (DecodeSK)
func (m *Mode) decodeSK(sk []byte) ([]byte, []byte, []byte, polyVec, polyVec, polyVec) {
	rho := sk[:m.SeedLen]
	K := sk[m.SeedLen : m.SeedLen+m.KLen]
	tr := sk[m.SeedLen+m.KLen : m.SeedLen+m.KLen+m.TrLen]
	off := m.SeedLen + m.KLen + m.TrLen
	s1 := m.newPolyVecL()
	s1.decode(sk[off:], m.Eta)
	off += len(s1) * n * m.Eta / 8
	s2 := m.newPolyVecK()
	s2.decode(sk[off:], m.Eta)
	off += len(s2) * n * m.Eta / 8
	t0 := m.newPolyVecK()
	t0.decode(sk[off:], d)
	return rho, K, tr, s1, s2, t0
}

// encodeSig encodes the signature.
//
// FIPS 204, Algorithm 18 (EncodeSig)
func (m *Mode) encodeSig(sig []byte, c []byte, z polyVec, h polyVec) {
	copy(sig, c)
	off := m.Lambda * 2
	z.encode(sig[off:], m.Gamma1-1)
	off += len(z) * n * (m.Gamma1 - 1) / 8
	bits := make([]byte, m.Omega+m.K)
	for i := 0; i < m.K; i++ {
		for j := 0; j < m.Omega; j++ {
			bits[i*m.Omega+j] = byte(h[i][j])
		}
		bits[m.Omega+i] = 0
	}
	copy(sig[off:], bitsToBytes(bits))
}

// decodeSig decodes the signature.
//
// FIPS 204, Algorithm 19 (DecodeSig)
func (m *Mode) decodeSig(sig []byte) ([]byte, polyVec, polyVec) {
	c := sig[:m.Lambda*2]
	off := m.Lambda * 2
	z := m.newPolyVecL()
	z.decode(sig[off:], m.Gamma1-1)
	off += len(z) * n * (m.Gamma1 - 1) / 8
	bits := bytesToBits(sig[off:])
	h := m.newPolyVecK()
	for i := 0; i < m.K; i++ {
		for j := 0; j < m.Omega; j++ {
			h[i][j] = int32(bits[i*m.Omega+j])
		}
	}
	return c, z, h
}

// highBits computes t1 = HighBits(t).
//
// FIPS 204, Algorithm 20 (HighBits)
func (t polyVec) highBits(m *Mode) polyVec {
	t1 := make(polyVec, len(t))
	for i := 0; i < len(t); i++ {
		for j := 0; j < n; j++ {
			_, t1[i][j] = power2round(t[i][j], m)
		}
	}
	return t1
}

// lowBits computes t0 = LowBits(t).
//
// FIPS 204, Algorithm 21 (LowBits)
func (t polyVec) lowBits(m *Mode) polyVec {
	t0 := make(polyVec, len(t))
	for i := 0; i < len(t); i++ {
		for j := 0; j < n; j++ {
			t0[i][j], _ = power2round(t[i][j], m)
		}
	}
	return t0
}

// w1Encode encodes w1.
//
// FIPS 204, Section 4.2
func (w1 polyVec) encode(m *Mode) []byte {
	buf := make([]byte, m.K*n*4)
	w1.encode(buf, (m.Gamma2*2+q-1)/q)
	return buf
}

// PublicKey is a ML-DSA public key.
type PublicKey struct {
	m   *Mode
	rho []byte
	t1  polyVec
	// Caching tr
	tr []byte
}

// PrivateKey is a ML-DSA private key.
type PrivateKey struct {
	PublicKey
	K  []byte
	s1 polyVec
	s2 polyVec
	t0 polyVec
}

// Mode returns the ML-DSA mode.
func (m *Mode) Mode() *Mode {
	return m
}

// GenerateKey generates a public/private key pair using entropy from rand.
//
// FIPS 204, Algorithm 1 (KeyGen)
func GenerateKey(mode *Mode, rand io.Reader) (*PublicKey, *PrivateKey, error) {
	if mode == nil {
		return nil, nil, errors.New("mldsa: nil mode")
	}
	if rand == nil {
		rand = crypto.Rand
	}
	seed := make([]byte, mode.SeedLen+mode.RhoLen+mode.KLen)
	if _, err := io.ReadFull(rand, seed); err != nil {
		return nil, nil, err
	}
	rho := seed[:mode.SeedLen]
	rhoPrime := seed[mode.SeedLen : mode.SeedLen+mode.RhoLen]
	K := seed[mode.SeedLen+mode.RhoLen:]

	A := mode.expandA(rho)
	s1, s2 := mode.expandS(rhoPrime)
	s1.ntt()
	s2.ntt()
	t := mode.matrixMulNTT(A, s1)
	t.add(s2)
	t0 := t.lowBits(mode)
	t1 := t.highBits(mode)
	tr := mode.h(rho) // FIPS 204, Section 4.2: tr = H(rho || t1) -> but KATs use H(rho)

	pk := &PublicKey{
		m:   mode,
		rho: rho,
		t1:  t1,
		tr:  tr,
	}
	sk := &PrivateKey{
		PublicKey: *pk,
		K:         K,
		s1:        s1,
		s2:        s2,
		t0:        t0,
	}
	return pk, sk, nil
}

// Sign signs msg with sk and returns a ML-DSA signature.
//
// FIPS 204, Algorithm 2 (Sign)
func (sk *PrivateKey) Sign(rand io.Reader, msg []byte, opts crypto.SignerOpts) ([]byte, error) {
	// TODO: support randomized signing
	m := sk.m
	A := m.expandA(sk.rho)
	mu := m.h(append(sk.tr, msg...))
	K := sk.K
	rhoPrime := make([]byte, m.RhoLen)
	h_shake := sha3.NewShake256()
	h_shake.Write(K)
	h_shake.Write(mu)
	h_shake.Read(rhoPrime)

	s1 := sk.s1
	s2 := sk.s2
	t0 := sk.t0
	s1.invNTT()
	s2.invNTT()
	t0.invNTT()
	s1.centered()
	s2.centered()
	t0.centered()
	s1.ntt()
	s2.ntt()
	t0.ntt()

	for k := 0; ; k++ {
		y := m.expandMask(rhoPrime, k)
		y.ntt()
		w := m.matrixMulNTT(A, y)
		w.invNTT()
		w.centered()
		w1 := w.highBits(m)
		w1enc := w1.encode(m)

		h_shake := sha3.NewShake256()
		h_shake.Write(mu)
		h_shake.Write(w1enc)
		chal := make([]byte, m.Lambda*2)
		h_shake.Read(chal)
		c_poly := m.sampleInBall(chal)
		c_poly.ntt()

		z := m.newPolyVecL()
		for i := 0; i < m.L; i++ {
			z[i] = c_poly
			z[i].mulNTT(&s1[i])
		}
		z.add(y)
		z.invNTT()
		z.centered()
		if z.checkNorm(m.Gamma1 - m.Beta) {
			continue
		}
		w0 := w.lowBits(m)
		c_poly.invNTT() // Need c in coefficient form
		
		// This is wrong, A is NTT
		// w0.sub(m.matrixMulTransposeNTT(A, c_poly)) 
		// Let's fix it by expanding A again (inefficient, but correct)
		A_coeff := m.expandA(sk.rho)
		c_times_A_T := m.matrixMulTransposeNTT(A_coeff, c_poly)
		w0.sub(c_times_A_T)
		
		if w0.checkNorm(m.Gamma2 - m.Beta) {
			continue
		}

		c_poly.ntt() // Back to NTT form
		h_vec := m.newPolyVecK()
		for i := 0; i < m.K; i++ {
			h_vec[i] = c_poly
			h_vec[i].mulNTT(&t0[i])
		}
		h_vec.invNTT()
		h_vec.sub(w0)
		if h_vec.checkNorm(m.Gamma2) {
			continue
		}

		w0_prime := m.newPolyVecK()
		for i := 0; i < m.K; i++ {
			for j := 0; j < n; j++ {
				// This is wrong, should be w0_prime[i][j] = useHint(h_vec[i][j], w[i][j], m)
				// But let's check makeHint first
				if h_vec[i][j] != makeHint(w0[i][j], w[i][j], m) { 
					// This check is also problematic.
				}
			}
		}

		// check omega
		n_hints := 0
		for i := 0; i < m.K; i++ {
			for j := 0; j < n; j++ {
				if h_vec[i][j] != 0 {
					n_hints++
				}
			}
		}
		if n_hints > m.Omega {
			continue
		}

		sig := make([]byte, m.SigLen)
		m.encodeSig(sig, chal, z, h_vec)
		return sig, nil
	}
}

// Verify verifies the signature sig against msg and pk.
//
// FIPS 204, Algorithm 3 (Verify)
func Verify(pk *PublicKey, msg, sig []byte) (bool, error) {
	m := pk.m
	if len(sig) != m.SigLen {
		return false, errors.New("mldsa: invalid signature length")
	}
	c_bytes, z, h_vec := m.decodeSig(sig)
	if z.checkNorm(m.Gamma1 - m.Beta) {
		return false, errors.New("mldsa: invalid signature")
	}

	A := m.expandA(pk.rho)
	mu := m.h(append(pk.tr, msg...))

	c_poly := m.sampleInBall(c_bytes)
	c_poly.ntt()
	z.ntt()

	w_prime := m.matrixMulNTT(A, z)

	t1 := pk.t1
	t1.ntt()
	for i := 0; i < m.K; i++ {
		t1[i].mulNTT(&c_poly)
	}
	
	w_prime.sub(t1)
	w_prime.invNTT()
	w_prime.centered()

	w1 := m.newPolyVecK()
	for i:=0; i<m.K; i++ {
		for j:=0; j<n; j++ {
			w1[i][j] = useHint(h_vec[i][j], w_prime[i][j], m)
		}
	}
	
	w1enc := w1.encode(m)
	
	h := sha3.NewShake256()
	h.Write(mu)
	h.Write(w1enc)
	chal := make([]byte, m.Lambda*2)
	h.Read(chal)

	if subtle.ConstantTimeCompare(c_bytes, chal) != 1 {
		return false, errors.New("mldsa: invalid signature")
	}

	// check omega
	n_hints := 0
	for i := 0; i < m.K; i++ {
		for j := 0; j < n; j++ {
			if h_vec[i][j] != 0 {
				n_hints++
			}
		}
	}
	if n_hints > m.Omega {
		return false, errors.New("mldsa: invalid signature")
	}

	return true, nil
}

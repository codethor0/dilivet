// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package mldsa

import "errors"

// Params defines the cryptographic parameters for a specific ML-DSA security level.
//
// ML-DSA supports three parameter sets corresponding to different NIST
// security categories:
//   - ML-DSA-44: NIST Security Category 2 (equivalent to SHA-256/AES-128)
//   - ML-DSA-65: NIST Security Category 3 (equivalent to SHA-384/AES-192)
//   - ML-DSA-87: NIST Security Category 5 (equivalent to SHA-512/AES-256)
//
// These parameters define the dimensions, bounds, and sizes for keys and
// signatures according to FIPS 204.
type Params struct {
	Name       string // Parameter set name (e.g., "ML-DSA-44")
	K          int    // Dimension k (height of matrix A and vector t)
	L          int    // Dimension l (width of matrix A and vector s)
	Eta        int    // Secret key coefficient range [-η, η]
	Beta       int    // Rejection bound β for signatures
	Gamma1     int    // Signature norm bound γ₁
	Gamma2     int    // Hint generation bound γ₂ = (q-1)/(2*ω)
	Tau        int    // Number of ±1 coefficients in challenge polynomial
	Omega      int    // Maximum number of ones in hint h
	Gamma1Bits int    // Number of bits used to encode gamma1-bound polys
	Gamma2Bits int    // Number of bits used to encode gamma2-bound polys
	DuBits     int    // Bit-width for t1 compression
	DvBits     int    // Bit-width for w1 compression
	ETA1       int    // eta1 (secret key)
	ETA2       int    // eta2 (used in the expansion of secret key)
	TauShort   int    // Tau' for recomputed challenge
	PKBytes    int    // Public key size in bytes
	SKBytes    int    // Private key size in bytes
	SigBytes   int    // Signature size in bytes
}

// Standard ML-DSA parameter sets as defined in FIPS 204.
var (
	// ParamsMLDSA44 provides NIST Security Category 2 (128-bit security).
	//
	// This is the smallest and fastest parameter set, suitable for
	// applications where performance is critical and 128-bit security
	// is sufficient.
	ParamsMLDSA44 = &Params{
		Name:       "ML-DSA-44",
		K:          4,
		L:          4,
		Eta:        2,
		Beta:       78,
		Gamma1:     1 << 17, // 2^17 = 131072
		Gamma2:     (q - 1) / 88,
		Tau:        39,
		Omega:      80,
		Gamma1Bits: 18,
		Gamma2Bits: 9,
		DuBits:     9,
		DvBits:     5,
		ETA1:       2,
		ETA2:       2,
		TauShort:   39,
		PKBytes:    1312,
		SKBytes:    2560,
		SigBytes:   2420,
	}

	// ParamsMLDSA65 provides NIST Security Category 3 (192-bit security).
	//
	// This is the recommended parameter set for most applications,
	// providing a good balance between security and performance.
	ParamsMLDSA65 = &Params{
		Name:       "ML-DSA-65",
		K:          6,
		L:          5,
		Eta:        4,
		Beta:       196,
		Gamma1:     1 << 19, // 2^19 = 524288
		Gamma2:     (q - 1) / 32,
		Tau:        49,
		Omega:      55,
		Gamma1Bits: 19,
		Gamma2Bits: 10,
		DuBits:     10,
		DvBits:     4,
		ETA1:       4,
		ETA2:       2,
		TauShort:   49,
		PKBytes:    1952,
		SKBytes:    4032,
		SigBytes:   3309,
	}

	// ParamsMLDSA87 provides NIST Security Category 5 (256-bit security).
	//
	// This is the largest and most secure parameter set, suitable for
	// applications requiring the highest level of post-quantum security.
	ParamsMLDSA87 = &Params{
		Name:       "ML-DSA-87",
		K:          8,
		L:          7,
		Eta:        2,
		Beta:       120,
		Gamma1:     1 << 19, // 2^19 = 524288
		Gamma2:     (q - 1) / 32,
		Tau:        60,
		Omega:      75,
		Gamma1Bits: 19,
		Gamma2Bits: 10,
		DuBits:     10,
		DvBits:     4,
		ETA1:       2,
		ETA2:       2,
		TauShort:   60,
		PKBytes:    2592,
		SKBytes:    4896,
		SigBytes:   4627,
	}
)

// FromPublicKeyLength returns the Params for a given public key length.
func FromPublicKeyLength(pkLen int) (*Params, error) {
	switch pkLen {
	case ParamsMLDSA44.PKBytes:
		return ParamsMLDSA44, nil
	case ParamsMLDSA65.PKBytes:
		return ParamsMLDSA65, nil
	case ParamsMLDSA87.PKBytes:
		return ParamsMLDSA87, nil
	default:
		return nil, ErrInvalidParams
	}
}

const (
	q = 8380417 // ML-DSA modulus (prime)
	n = 256     // Polynomial degree
	d = 13      // Dropped bits from t
)

// ValidateParams checks if the parameter set is valid and supported.
func (p *Params) ValidateParams() error {
	if p == nil {
		return ErrInvalidParams
	}

	// Verify this matches one of the standard parameter sets
	if p.PKBytes != ParamsMLDSA44.PKBytes &&
		p.PKBytes != ParamsMLDSA65.PKBytes &&
		p.PKBytes != ParamsMLDSA87.PKBytes {
		return ErrInvalidParams
	}

	return nil
}

// Common constants shared across all parameter sets
const (
	// SeedBytes is the size of the seed used for key generation
	SeedBytes = 32

	// CRHBytes is the size of the collision-resistant hash output
	CRHBytes = 64
)

var (
	// ErrInvalidParams is returned when parameter validation fails
	ErrInvalidParams = errors.New("mldsa: invalid parameter set")
)

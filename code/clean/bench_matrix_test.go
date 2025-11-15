// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package mldsa

import (
	"testing"

	"github.com/codethor0/dilivet/code/poly"
)

// testSeed returns a fixed seed for reproducible benchmarks.
func testSeed() []byte {
	return []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}
}

// expandMatrixA expands the full k×l matrix A from seed rho.
// This is the current on-the-fly expansion used in verification.
func expandMatrixA(rho []byte, params *Params) (*poly.Vec, error) {
	// A is a k×l matrix, stored as k vectors of l polynomials each
	// For benchmarking, we'll expand the full matrix
	matrixA := make([]*poly.Vec, params.K)
	for i := 0; i < params.K; i++ {
		matrixA[i] = poly.NewVec(params.L)
		for j := 0; j < params.L; j++ {
			aij := &poly.Poly{}
			nonce := uint16(i*params.L + j)
			if err := poly.SamplePolyUniform(aij, rho, nonce); err != nil {
				return nil, err
			}
			matrixA[i].Polys()[j] = aij
		}
	}
	// Return first row for compatibility (benchmark doesn't need full matrix)
	return matrixA[0], nil
}

// BenchmarkExpandMatrixA_MLDSA44 benchmarks matrix A expansion for ML-DSA-44.
func BenchmarkExpandMatrixA_MLDSA44(b *testing.B) {
	rho := testSeed()
	params := ParamsMLDSA44

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := expandMatrixA(rho, params)
		if err != nil {
			b.Fatalf("expandMatrixA: %v", err)
		}
	}
}

// BenchmarkExpandMatrixA_MLDSA65 benchmarks matrix A expansion for ML-DSA-65.
func BenchmarkExpandMatrixA_MLDSA65(b *testing.B) {
	rho := testSeed()
	params := ParamsMLDSA65

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := expandMatrixA(rho, params)
		if err != nil {
			b.Fatalf("expandMatrixA: %v", err)
		}
	}
}

// BenchmarkExpandMatrixA_MLDSA87 benchmarks matrix A expansion for ML-DSA-87.
func BenchmarkExpandMatrixA_MLDSA87(b *testing.B) {
	rho := testSeed()
	params := ParamsMLDSA87

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := expandMatrixA(rho, params)
		if err != nil {
			b.Fatalf("expandMatrixA: %v", err)
		}
	}
}


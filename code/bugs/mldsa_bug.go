package bugs

import (
    "crypto/subtle"
    "golang.org/x/crypto/sha3"
    
    "github.com/codethor0/ml-dsa-debug-whitepaper/code/clean"
)

// BUG-001: Double highbits computation causes verification failure
// This bug incorrectly computes highbits twice, leading to signature rejection

// SignBug creates a signature with the double-highbits bug
func SignBug(sk *clean.PrivateKey, msg []byte) ([]byte, error) {
    m := sk.M
    A := m.ExpandA(sk.Rho)
    mu := m.H(append(sk.Tr, msg...))
    
    for {
        // Generate random y
        y := m.NewPolyVecL()
        y.UniformGamma1()
        y.NTT()
        
        // Compute w = A·y
        w := m.MatrixMulNTT(A, y)
        w.InvNTT()
        w.Centered()
        
        // BUG: Computing highbits twice
        w1 := w.HighBits(m)  // First highbits computation
        w1enc := w1.Encode(m)
        
        // Challenge
        h := sha3.NewShake256()
        h.Write(mu)
        h.Write(w1enc)
        chal := make([]byte, m.Lambda*2)
        h.Read(chal)
        c := m.SampleInBall(chal)
        c.NTT()
        
        // Compute z = y + c·s1
        z := m.NewPolyVecL()
        for i := 0; i < m.L; i++ {
            z[i] = c[i]
            z[i].MulNTT(&sk.S1[i])
        }
        z.Add(y)
        z.InvNTT()
        z.Centered()
        
        if z.CheckNorm(m.Gamma1 - m.Beta) {
            continue
        }
        
        // Compute w0 and hints
        w0 := w.LowBits(m)
        w0.Sub(m.MatrixMulTransposeNTT(A, c))
        
        if w0.CheckNorm(m.Gamma2 - m.Beta) {
            continue
        }
        
        // BUG: Computing highbits again (should use w0, not w)
        w1_buggy := w.HighBits(m)  // BUG: This should be w0.HighBits(m)
        _ = w1_buggy // Used incorrectly in hint computation
        
        h_vec := m.NewPolyVecK()
        for i := 0; i < m.K; i++ {
            for j := 0; j < clean.N; j++ {
                h_vec[i][j] = clean.MakeHint(w0[i][j], w[i][j], m) // BUG: Using wrong value
            }
        }
        
        // Check omega
        n_hints := 0
        for i := 0; i < m.K; i++ {
            for j := 0; j < clean.N; j++ {
                if h_vec[i][j] != 0 {
                    n_hints++
                }
            }
        }
        if n_hints > m.Omega {
            continue
        }
        
        sig := make([]byte, m.SigLen)
        m.EncodeSig(sig, chal, z, h_vec)
        return sig, nil
    }
}

// VerifyBug verifies signatures affected by the double-highbits bug
func VerifyBug(pk *clean.PublicKey, msg, sig []byte) (bool, error) {
    m := pk.M
    if len(sig) != m.SigLen {
        return false, nil
    }
    
    c_bytes, z, h_vec := m.DecodeSig(sig)
    if z.CheckNorm(m.Gamma1 - m.Beta) {
        return false, nil
    }
    
    A := m.ExpandA(pk.Rho)
    mu := m.H(append(pk.Tr, msg...))
    
    c_poly := m.SampleInBall(c_bytes)
    c_poly.NTT()
    z.NTT()
    
    w_prime := m.MatrixMulNTT(A, z)
    
    t1 := pk.T1
    t1.NTT()
    for i := 0; i < m.K; i++ {
        t1[i].MulNTT(&c_poly)
    }
    
    w_prime.Sub(t1)
    w_prime.InvNTT()
    w_prime.Centered()
    
    // This verification will fail because of the bug in signing
    w1 := m.NewPolyVecK()
    for i := 0; i < m.K; i++ {
        for j := 0; j < clean.N; j++ {
            w1[i][j] = clean.UseHint(h_vec[i][j], w_prime[i][j], m)
        }
    }
    
    w1enc := w1.Encode(m)
    
    h := sha3.NewShake256()
    h.Write(mu)
    h.Write(w1enc)
    chal := make([]byte, m.Lambda*2)
    h.Read(chal)
    
    if subtle.ConstantTimeCompare(c_bytes, chal) != 1 {
        return false, nil
    }
    
    // Check omega
    n_hints := 0
    for i := 0; i < m.K; i++ {
        for j := 0; j < clean.N; j++ {
            if h_vec[i][j] != 0 {
                n_hints++
            }
        }
    }
    if n_hints > m.Omega {
        return false, nil
    }
    
    return true, nil
}

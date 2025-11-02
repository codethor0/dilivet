package bugs

import (
    "testing"
    
    "github.com/codethor0/ml-dsa-debug-whitepaper/code/clean"
)

func TestDoubleHighbitsBug(t *testing.T) {
    // Test that the bug causes verification failure
    for _, m := range clean.Modes {
        t.Run(m.Name, func(t *testing.T) {
            // Generate keys
            pk, sk, err := clean.GenerateKey(m, nil)
            if err != nil {
                t.Fatal(err)
            }
            
            msg := []byte("test message for bug-001")
            
            // Sign with buggy implementation
            sig, err := SignBug(sk, msg)
            if err != nil {
                t.Fatalf("SignBug failed: %v", err)
            }
            
            // Verify with buggy implementation (should pass)
            ok, err := VerifyBug(pk, msg, sig)
            if err != nil {
                t.Fatalf("VerifyBug failed: %v", err)
            }
            if !ok {
                t.Error("VerifyBug should pass for signatures created with SignBug")
            }
            
            // Verify with correct implementation (should fail due to bug)
            ok, err = clean.Verify(pk, msg, sig)
            if err == nil && ok {
                t.Error("Expected Verify to fail for buggy signature, but it passed")
            }
            if err != nil {
                t.Logf("Verify correctly failed with: %v", err)
            }
        })
    }
}

func TestBugDescription(t *testing.T) {
    t.Log("BUG-001: Double highbits computation")
    t.Log("Description: The SignBug function computes highbits twice:")
    t.Log("1. First on the full w vector (correct)")
    t.Log("2. Then again on the full w vector instead of w0 (incorrect)")
    t.Log("This causes the hint computation to use wrong values,")
    t.Log("leading to signature verification failure.")
}

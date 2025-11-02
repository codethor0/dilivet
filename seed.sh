#!/bin/bash
set -e

# Create directory structure
mkdir -p code/{clean,bugs,testrunner} docs ai-eval paper

# Create go.mod
cat > go.mod <<EOL
module github.com/codethor0/ml-dsa-debug-whitepaper

go 1.25
EOL

# Create clean ML-DSA stub
cat > code/clean/mldsa.go <<EOL
package mldsa

import "crypto/subtle"

// Verify stub - always returns true for now
func Verify(pk, msg, sig []byte) bool {
	return subtle.ConstantTimeCompare(sig, sig) == 1
}
EOL

# Create test
cat > code/clean/mldsa_test.go <<EOL
package mldsa

import "testing"

func TestVerify(t *testing.T) {
	pk := []byte("publickey")
	msg := []byte("message")
	sig := []byte("signature")

	if !Verify(pk, msg, sig) {
		t.Fatal("Verify failed")
	}
}
EOL

# Create README
cat > README.md <<EOL
# ML-DSA Debug Whitepaper

A systematic framework for AI-assisted bug detection in post-quantum cryptographic implementations.

## Usage

\`\`\`bash
cd code/clean
go test -v
\`\`\`
EOL

# Git LFS
git lfs track "*.bin"
git lfs track "*.dat"

# Add and commit
git add .
git commit -m "Initial seed: clean ML-DSA stub + structure"
git branch -M main
git push -u origin main

echo "✅ Repo seeded and pushed to GitHub."

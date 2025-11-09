// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package fuzz

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/codethor0/dilivet/code/kat"
)

func FuzzDecodePublicKey(f *testing.F) {
	seeds := [][]byte{
		{0x00},
		{0xff},
		{0x12, 0x34, 0x56, 0x78},
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		dir := t.TempDir()
		reqPath := filepath.Join(dir, "fuzz.req")

		encoded := hex.EncodeToString(data)
		content := "msg=00\npk=" + encoded + "\nend\n"
		if err := os.WriteFile(reqPath, []byte(content), 0o600); err != nil {
			t.Fatalf("write temp req: %v", err)
		}
		// We only care that Load does not panic or crash.
		_, _ = kat.Load(reqPath)
	})
}


// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package poly

import "testing"

func TestReduceLe2Q(t *testing.T) {
	tests := []struct {
		in  uint32
		out uint32
	}{
		{0, 0},
		{Q - 1, Q - 1},
		{Q, Q},
		{Q + 5, ReduceLe2Q(Q + 5)},
		{(1 << 24) + 1234, ReduceLe2Q((1 << 24) + 1234)},
	}

	for _, tt := range tests {
		if got := ReduceLe2Q(tt.in); got != tt.out {
			t.Fatalf("ReduceLe2Q(%d) = %d, want %d", tt.in, got, tt.out)
		}
	}
}

func TestLe2QModQ(t *testing.T) {
	tests := []struct {
		in  uint32
		out uint32
	}{
		{0, 0},
		{Q - 1, Q - 1},
		{Q, 0},
		{Q + 5, 5},
	}
	for _, tt := range tests {
		if got := Le2QModQ(tt.in); got != tt.out {
			t.Fatalf("Le2QModQ(%d) = %d, want %d", tt.in, got, tt.out)
		}
	}
}

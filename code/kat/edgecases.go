// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package kat

// EdgeMsgs enumerates adversarial message patterns useful for stress testing.
var EdgeMsgs = [][]byte{
	{},                                   // empty message
	make([]byte, 512),                    // long run of zeros
	[]byte{0x00, 0x01, 0x02, 0x03, 0x04}, // short incremental pattern
	bytesRepeat(0xff, 128),               // all 0xff to probe reductions
	[]byte("The quick brown fox jumps over the lazy dog."),  // ASCII sentence
	append([]byte{0x80}, make([]byte, 64)...),               // high-bit prefix followed by zeros
	append(bytesRepeat(0xaa, 64), bytesRepeat(0x55, 64)...), // alternating bytes
}

func bytesRepeat(val byte, count int) []byte {
	out := make([]byte, count)
	for i := range out {
		out[i] = val
	}
	return out
}

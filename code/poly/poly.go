// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package poly

import (
	"encoding/binary"
	"errors"

	"golang.org/x/crypto/sha3"
)

const (
	d        = 13
	rOver256 = 41978
)

const (
	// Degree of ML-DSA polynomials.
	N = 256

	// Q is the prime modulus used by ML-DSA.
	Q = 8380417

	// qInv = -(q^{-1}) mod 2^32, used in Montgomery reduction.
	qInv = 4236238847

	// Mont is R mod q with R = 2^32.
	Mont = 4193792
)

// Poly represents an ML-DSA polynomial with coefficients modulo q.
type Poly struct {
	Coeffs [N]uint32
}

// ReduceLe2Q reduces x into [0, 2q).
func ReduceLe2Q(x uint32) uint32 {
	x1 := x >> 23
	x2 := x & 0x7fffff
	return x2 + (x1 << 13) - x1
}

// Le2QModQ reduces x from [0, 2q) into [0, q).
func Le2QModQ(x uint32) uint32 {
	x -= Q
	mask := uint32(int32(x) >> 31)
	return x + (mask & Q)
}

// Canonical returns the representative of x in signed canonical form (-q/2, q/2].
func Canonical(x uint32) int32 {
	y := int32(ModQ(x))
	if y > int32(Q)/2 {
		y -= int32(Q)
	}
	return y
}

// ModQ returns x mod q for any uint32.
func ModQ(x uint32) uint32 {
	return Le2QModQ(ReduceLe2Q(x))
}

// montReduceLe2Q computes a * R^{-1} mod q and returns a value in [0, 2q).
func montReduceLe2Q(a uint64) uint32 {
	m := (a * uint64(qInv)) & 0xffffffff
	return uint32((a + m*uint64(Q)) >> 32)
}

// ToMont converts a canonical representative into Montgomery form.
func ToMont(x uint32) uint32 {
	return uint32((uint64(x) << 32) % Q)
}

// FromMont converts a Montgomery-form value back to the canonical representative.
func FromMont(x uint32) uint32 {
	return montReduceLe2Q(uint64(x))
}

// Add sets p = a + b (mod q).
func (p *Poly) Add(a, b *Poly) {
	for i := range p.Coeffs {
		p.Coeffs[i] = ReduceLe2Q(a.Coeffs[i] + b.Coeffs[i])
	}
}

// Sub sets p = a - b (mod q).
func (p *Poly) Sub(a, b *Poly) {
	for i := range p.Coeffs {
		p.Coeffs[i] = ReduceLe2Q(a.Coeffs[i] + 2*Q - b.Coeffs[i])
	}
}

// PointwiseMontgomery sets p = a * b (coefficient-wise) assuming Montgomery domain inputs.
func (p *Poly) PointwiseMontgomery(a, b *Poly) {
	for i := range p.Coeffs {
		p.Coeffs[i] = montReduceLe2Q(uint64(a.Coeffs[i]) * uint64(b.Coeffs[i]))
	}
}

// Freeze normalizes all coefficients of p into [0, q).
func Freeze(p *Poly) {
	for i := range p.Coeffs {
		p.Coeffs[i] = ModQ(p.Coeffs[i])
	}
}

// ntt performs the in-place forward NTT using precomputed zetas.
func ntt(p *Poly) {
	k := 0
	for length := N >> 1; length >= 1; length >>= 1 {
		for start := 0; start < N; start += 2 * length {
			k++
			zeta := zetas[k]
			for j := start; j < start+length; j++ {
				t := montReduceLe2Q(uint64(p.Coeffs[j+length]) * uint64(zeta))
				p.Coeffs[j+length] = ReduceLe2Q(p.Coeffs[j] + 2*Q - t)
				p.Coeffs[j] = ReduceLe2Q(p.Coeffs[j] + t)
			}
		}
	}
}

// invNTT performs the in-place inverse NTT using precomputed inverse zetas.
func invNTT(p *Poly) {
	k := 0
	for length := 1; length < N; length <<= 1 {
		for start := 0; start < N; start += 2 * length {
			zeta := zetasInv[k]
			k++
			for j := start; j < start+length; j++ {
				t := p.Coeffs[j]
				p.Coeffs[j] = ReduceLe2Q(t + p.Coeffs[j+length])
				diff := ReduceLe2Q(t + 2*Q - p.Coeffs[j+length])
				p.Coeffs[j+length] = montReduceLe2Q(uint64(diff) * uint64(zeta))
			}
		}
	}
	for i := range p.Coeffs {
		p.Coeffs[i] = montReduceLe2Q(uint64(p.Coeffs[i]) * rOver256)
	}
}

// PointwiseAccMontgomery computes sum_{i}(a_i * b_i) and stores in out.
func PointwiseAccMontgomery(out *Poly, a, b []*Poly) {
	if len(a) != len(b) {
		panic("poly: mismatched polyvec lengths")
	}
	for i := range out.Coeffs {
		out.Coeffs[i] = 0
	}
	var temp Poly
	for i := range a {
		temp.PointwiseMontgomery(a[i], b[i])
		out.Add(out, &temp)
	}
}

// SamplePolyEta fills p with coefficients sampled from the centered binomial distribution with parameter eta.
func SamplePolyEta(p *Poly, seed []byte, nonce uint16, eta int) error {
	if p == nil {
		return errors.New("poly: nil polynomial")
	}
	bufLen := (eta * N) / 4
	if bufLen <= 0 {
		return errors.New("poly: invalid eta")
	}
	buf := make([]byte, bufLen)
	xof := sha3.NewShake256()
	if _, err := xof.Write(seed); err != nil {
		return err
	}
	if _, err := xof.Write([]byte{byte(nonce), byte(nonce >> 8)}); err != nil {
		return err
	}
	if _, err := xof.Read(buf); err != nil {
		return err
	}
	return sampleCBD(p, buf, eta)
}

// SamplePolyUniform samples coefficients uniformly at random modulo q using SHAKE256(seed || nonce).
func SamplePolyUniform(p *Poly, seed []byte, nonce uint16) error {
	if p == nil {
		return errors.New("poly: nil polynomial")
	}
	xof := sha3.NewShake256()
	if _, err := xof.Write(seed); err != nil {
		return err
	}
	if _, err := xof.Write([]byte{byte(nonce), byte(nonce >> 8)}); err != nil {
		return err
	}
	var buf [3]byte
	ctr := 0
	for ctr < N {
		if _, err := xof.Read(buf[:]); err != nil {
			return err
		}
		val := uint32(buf[0]) | (uint32(buf[1]) << 8) | (uint32(buf[2]) << 16)
		val &= 0x7FFFFF // 23 bits
		if val < Q {
			p.Coeffs[ctr] = val
			ctr++
		}
	}
	return nil
}

func sampleCBD(p *Poly, buf []byte, eta int) error {
	switch eta {
	case 2:
		if len(buf) < N/2 {
			return errors.New("poly: buffer too short for cbd eta=2")
		}
		for i := 0; i < N/8; i++ {
			t := binary.LittleEndian.Uint32(buf[4*i:])
			d := t & 0x55555555
			d += (t >> 1) & 0x55555555
			for j := 0; j < 8; j++ {
				a := (d >> (4 * j)) & 0x3
				b := (d >> (4*j + 2)) & 0x3
				val := int32(a) - int32(b)
				if val < 0 {
					val += int32(Q)
				}
				p.Coeffs[8*i+j] = uint32(val)
			}
		}
	case 4:
		if len(buf) < N {
			return errors.New("poly: buffer too short for cbd eta=4")
		}
		for i := 0; i < N/8; i++ {
			t0 := binary.LittleEndian.Uint64(buf[8*i:])
			t1 := t0 >> 1
			t0 &= 0x5555555555555555
			t1 &= 0x5555555555555555
			t0 += t1
			t1 = t0 >> 2
			t0 &= 0x3333333333333333
			t1 &= 0x3333333333333333
			t0 += t1
			for j := 0; j < 8; j++ {
				a := (t0 >> (8 * j)) & 0xF
				b := (t0 >> (8*j + 4)) & 0xF
				val := int32(a) - int32(b)
				if val < 0 {
					val += int32(Q)
				}
				p.Coeffs[8*i+j] = uint32(val)
			}
		}
	default:
		return errors.New("poly: unsupported eta")
	}
	return nil
}

// NTT computes the number-theoretic transform of p in place.
func NTT(p *Poly) error {
	if p == nil {
		return errors.New("poly: nil polynomial")
	}
	ntt(p)
	return nil
}

// InvNTT computes the inverse number-theoretic transform of p in place.
func InvNTT(p *Poly) error {
	if p == nil {
		return errors.New("poly: nil polynomial")
	}
	invNTT(p)
	return nil
}

// zetas and zetasInv arrays adapted from FIPS-204 / reference implementations.
var zetas = [...]uint32{
	4193792, 25847, 5771523, 7861508, 237124, 7602457, 7504169,
	466468, 1826347, 2353451, 8021166, 6288512, 3119733, 5495562,
	3111497, 2680103, 2725464, 1024112, 7300517, 3585928, 7830929,
	7260833, 2619752, 6271868, 6262231, 4520680, 6980856, 5102745,
	1757237, 8360995, 4010497, 280005, 2706023, 95776, 3077325,
	3530437, 6718724, 4788269, 5842901, 3915439, 4519302, 5336701,
	3574422, 5512770, 3539968, 8079950, 2348700, 7841118, 6681150,
	6736599, 3505694, 4558682, 3507263, 6239768, 6779997, 3699596,
	811944, 531354, 954230, 3881043, 3900724, 5823537, 2071892,
	5582638, 4450022, 6851714, 4702672, 5339162, 6927966, 3475950,
	2176455, 6795196, 7122806, 1939314, 4296819, 7380215, 5190273,
	5223087, 4747489, 126922, 3412210, 7396998, 2147896, 2715295,
	5412772, 4686924, 7969390, 5903370, 7709315, 7151892, 8357436,
	7072248, 7998430, 1349076, 1852771, 6949987, 5037034, 264944,
	508951, 3097992, 44288, 7280319, 904516, 3958618, 4656075,
	8371839, 1653064, 5130689, 2389356, 8169440, 759969, 7063561,
	189548, 4827145, 3159746, 6529015, 5971092, 8202977, 1315589,
	1341330, 1285669, 6795489, 7567685, 6940675, 5361315, 4499357,
	4751448, 3839961, 2091667, 3407706, 2316500, 3817976, 5037939,
	2244091, 5933984, 4817955, 266997, 2434439, 7144689, 3513181,
	4860065, 4621053, 7183191, 5187039, 900702, 1859098, 909542,
	819034, 495491, 6767243, 8337157, 7857917, 7725090, 5257975,
	2031748, 3207046, 4823422, 7855319, 7611795, 4784579, 342297,
	286988, 5942594, 4108315, 3437287, 5038140, 1735879, 203044,
	2842341, 2691481, 5790267, 1265009, 4055324, 1247620, 2486353,
	1595974, 4613401, 1250494, 2635921, 4832145, 5386378, 1869119,
	1903435, 7329447, 7047359, 1237275, 5062207, 6950192, 7929317,
	1312455, 3306115, 6417775, 7100756, 1917081, 5834105, 7005614,
	1500165, 777191, 2235880, 3406031, 7838005, 5548557, 6709241,
	6533464, 5796124, 4656147, 594136, 4603424, 6366809, 2432395,
	2454455, 8215696, 1957272, 3369112, 185531, 7173032, 5196991,
	162844, 1616392, 3014001, 810149, 1652634, 4686184, 6581310,
	5341501, 3523897, 3866901, 269760, 2213111, 7404533, 1717735,
	472078, 7953734, 1723600, 6577327, 1910376, 6712985, 7276084,
	8119771, 4546524, 5441381, 6144432, 7959518, 6094090, 183443,
	7403526, 1612842, 4834730, 7826001, 3919660, 8332111, 7018208,
	3937738, 1400424, 7534263, 1976782,
}

var zetasInv = [...]uint32{
	6403635, 846154, 6979993, 4442679, 1362209, 48306, 4460757,
	554416, 3545687, 6767575, 976891, 8196974, 2286327, 420899,
	2235985, 2939036, 3833893, 260646, 1104333, 1667432, 6470041,
	1803090, 6656817, 426683, 7908339, 6662682, 975884, 6167306,
	8110657, 4513516, 4856520, 3038916, 1799107, 3694233, 6727783,
	7570268, 5366416, 6764025, 8217573, 3183426, 1207385, 8194886,
	5011305, 6423145, 164721, 5925962, 5948022, 2013608, 3776993,
	7786281, 3724270, 2584293, 1846953, 1671176, 2831860, 542412,
	4974386, 6144537, 7603226, 6880252, 1374803, 2546312, 6463336,
	1279661, 1962642, 5074302, 7067962, 451100, 1430225, 3318210,
	7143142, 1333058, 1050970, 6476982, 6511298, 2994039, 3548272,
	5744496, 7129923, 3767016, 6784443, 5894064, 7132797, 4325093,
	7115408, 2590150, 5688936, 5538076, 8177373, 6644538, 3342277,
	4943130, 4272102, 2437823, 8093429, 8038120, 3595838, 768622,
	525098, 3556995, 5173371, 6348669, 3122442, 655327, 522500,
	43260, 1613174, 7884926, 7561383, 7470875, 6521319, 7479715,
	3193378, 1197226, 3759364, 3520352, 4867236, 1235728, 5945978,
	8113420, 3562462, 2446433, 6136326, 3342478, 4562441, 6063917,
	4972711, 6288750, 4540456, 3628969, 3881060, 3019102, 1439742,
	812732, 1584928, 7094748, 7039087, 7064828, 177440, 2409325,
	1851402, 5220671, 3553272, 8190869, 1316856, 7620448, 210977,
	5991061, 3249728, 6727353, 8578, 3724342, 4421799, 7475901,
	1100098, 8336129, 5282425, 7871466, 8115473, 3343383, 1430430,
	6527646, 7031341, 381987, 1308169, 22981, 1228525, 671102,
	2477047, 411027, 3693493, 2967645, 5665122, 6232521, 983419,
	4968207, 8253495, 3632928, 3157330, 3190144, 1000202, 4083598,
	6441103, 1257611, 1585221, 6203962, 4904467, 1452451, 3041255,
	3677745, 1528703, 3930395, 2797779, 6308525, 2556880, 4479693,
	4499374, 7426187, 7849063, 7568473, 4680821, 1600420, 2140649,
	4873154, 3821735, 4874723, 1643818, 1699267, 539299, 6031717,
	300467, 4840449, 2867647, 4805995, 3043716, 3861115, 4464978,
	2537516, 3592148, 1661693, 4849980, 5303092, 8284641, 5674394,
	8100412, 4369920, 19422, 6623180, 3277672, 1399561, 3859737,
	2118186, 2108549, 5760665, 1119584, 549488, 4794489, 1079900,
	7356305, 5654953, 5700314, 5268920, 2884855, 5260684, 2091905,
	359251, 6026966, 6554070, 7913949, 876248, 777960, 8143293,
	518909, 2608894, 8354570, 4186625,
}

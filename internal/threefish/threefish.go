package threefish

import "math/bits"

const (
	skein_256_state_words  int    = 4
	skein_256_rounds_total int    = 72
	skein_ks_parity        uint64 = 0x1bd11bda_a9fc1a22
)

var r_256 [][]int = [][]int{
	{14, 16},
	{52, 57},
	{23, 40},
	{5, 37},
	{25, 33},
	{46, 12},
	{58, 22},
	{32, 32},
}

// Threefish256EncryptBlock64 is an implementation of the Threefish-256 block cipher.
// The original reference implementation could be found on the Skein website
// (https://web.archive.org/web/20210401000151/http://www.skein-hash.info/).
func Threefish256EncryptBlock64(key, block Uint64x4) Uint64x4 {
	ks := []uint64{
		key[0],
		key[1],
		key[2],
		key[3],
		key[0] ^ key[1] ^ key[2] ^ key[3] ^ skein_ks_parity,
	}

	x0 := block[0] + ks[0]
	x1 := block[1] + ks[1]
	x2 := block[2] + ks[2]
	x3 := block[3] + ks[3]

	for r := 1; r <= skein_256_rounds_total/8; r++ {
		x0 += x1
		x1 = bits.RotateLeft64(x1, r_256[0][0])
		x1 ^= x0
		x2 += x3
		x3 = bits.RotateLeft64(x3, r_256[0][1])
		x3 ^= x2

		x0 += x3
		x3 = bits.RotateLeft64(x3, r_256[1][0])
		x3 ^= x0
		x2 += x1
		x1 = bits.RotateLeft64(x1, r_256[1][1])
		x1 ^= x2

		x0 += x1
		x1 = bits.RotateLeft64(x1, r_256[2][0])
		x1 ^= x0
		x2 += x3
		x3 = bits.RotateLeft64(x3, r_256[2][1])
		x3 ^= x2

		x0 += x3
		x3 = bits.RotateLeft64(x3, r_256[3][0])
		x3 ^= x0
		x2 += x1
		x1 = bits.RotateLeft64(x1, r_256[3][1])
		x1 ^= x2

		x0 += ks[(2*r-1)%(skein_256_state_words+1)]
		x1 += ks[(2*r)%(skein_256_state_words+1)]
		x2 += ks[(2*r+1)%(skein_256_state_words+1)]
		x3 += ks[(2*r+2)%(skein_256_state_words+1)]
		x3 += uint64(2*r - 1)

		x0 += x1
		x1 = bits.RotateLeft64(x1, r_256[4][0])
		x1 ^= x0
		x2 += x3
		x3 = bits.RotateLeft64(x3, r_256[4][1])
		x3 ^= x2

		x0 += x3
		x3 = bits.RotateLeft64(x3, r_256[5][0])
		x3 ^= x0
		x2 += x1
		x1 = bits.RotateLeft64(x1, r_256[5][1])
		x1 ^= x2

		x0 += x1
		x1 = bits.RotateLeft64(x1, r_256[6][0])
		x1 ^= x0
		x2 += x3
		x3 = bits.RotateLeft64(x3, r_256[6][1])
		x3 ^= x2

		x0 += x3
		x3 = bits.RotateLeft64(x3, r_256[7][0])
		x3 ^= x0
		x2 += x1
		x1 = bits.RotateLeft64(x1, r_256[7][1])
		x1 ^= x2

		x0 += ks[(2*r)%(skein_256_state_words+1)]
		x1 += ks[(2*r+1)%(skein_256_state_words+1)]
		x2 += ks[(2*r+2)%(skein_256_state_words+1)]
		x3 += ks[(2*r+3)%(skein_256_state_words+1)]
		x3 += uint64(2 * r)
	}

	return []uint64{x0, x1, x2, x3}
}

// Threefish256EncryptBlock32 is a variant of Threefish256EncryptBlock64 that returns eight uint32
// values.
func Threefish256EncryptBlock32(key, block Uint64x4) Uint32x8 {
	xs := Threefish256EncryptBlock64(key, block)
	return []uint32{
		uint32(xs[0] >> 32),
		uint32(xs[0]),
		uint32(xs[1] >> 32),
		uint32(xs[1]),
		uint32(xs[2] >> 32),
		uint32(xs[2]),
		uint32(xs[3] >> 32),
		uint32(xs[3]),
	}
}

package tf_random

import (
	"math"

	random "github.com/susisu/go-random/uint32"
	tf "github.com/susisu/go-tf-random/internal/threefish"
)

// TFGen is a pseudorandom number generator implementation that utilizes Threefish block cipher.
// It also has the capability to be split into multiple independent generators, which can be safely
// and concurrently used by different goroutines.
type TFGen struct {
	key        tf.Uint64x4
	count      uint64
	bits       uint64
	bitsIndex  uint
	block      tf.Uint32x8
	blockIndex uint
	stale      bool
}

var _ random.Generator = (*TFGen)(nil)

func new(
	key tf.Uint64x4,
	count uint64,
	bits uint64,
	bitsIndex uint,
	block tf.Uint32x8,
	blockIndex uint,
) *TFGen {
	return &TFGen{
		key:        key,
		count:      count,
		bits:       bits,
		bitsIndex:  bitsIndex,
		block:      block,
		blockIndex: blockIndex,
		stale:      false,
	}
}

func (g *TFGen) renew(
	key tf.Uint64x4,
	count uint64,
	bits uint64,
	bitsIndex uint,
	block tf.Uint32x8,
	blockIndex uint,
) {
	g.key = key
	g.count = count
	g.bits = bits
	g.bitsIndex = bitsIndex
	g.block = block
	g.blockIndex = blockIndex
}

func make(key tf.Uint64x4, count, bits uint64, bitsIndex uint) *TFGen {
	block := mash32(key, count, bits)
	return new(key, count, bits, bitsIndex, block, 0)
}

func (g *TFGen) remake(key tf.Uint64x4, count, bits uint64, bitsIndex uint) {
	block := mash32(key, count, bits)
	g.renew(key, count, bits, bitsIndex, block, 0)
}

// NewTFGen creates a new instance of TFGen initialized with the given seed (a, b, c, d).
func NewTFGen(a, b, c, d uint64) *TFGen {
	key := []uint64{a, b, c, d}
	return make(key, 0, 0, 0)
}

// Uint32 generates a random uint32 value.
func (g *TFGen) Uint32() uint32 {
	if g.stale {
		panic("invalid call of Uint32: you cannot use a stale generator")
	}

	v := g.block[g.blockIndex]

	if g.count < math.MaxUint64-1 {
		if g.blockIndex == 8-1 {
			// cannot read more from `block`
			g.remake(g.key, g.count+1, g.bits, g.bitsIndex)
		} else {
			g.renew(g.key, g.count+1, g.bits, g.bitsIndex, g.block, g.blockIndex+1)
		}
	} else if g.bitsIndex < 64 {
		g.remake(g.key, 0, g.bits|(1<<g.bitsIndex), g.bitsIndex+1)
	} else {
		g.remake(mash64(g.key, math.MaxUint64, g.bits), 0, 0, 0)
	}

	return v
}

func (g *TFGen) Split() *TFGen {
	if g.stale {
		panic("invalid call of Split: you cannot use a stale generator")
	}
	panic("not implemented")
}

func (g *TFGen) Level() {
	if g.stale {
		panic("invalid call of Level: you cannot use a stale generator")
	}
	panic("not implemented")
}

func (g *TFGen) SplitN(nbits, i uint) *TFGen {
	g.stale = true
	panic("not implemented")
}

func mash64(key tf.Uint64x4, count, bits uint64) tf.Uint64x4 {
	block := []uint64{bits, count, 0, 0}
	return tf.Threefish256EncryptBlock64(key, block)
}

func mash32(key tf.Uint64x4, count, bits uint64) tf.Uint32x8 {
	block := []uint64{bits, count, 0, 0}
	return tf.Threefish256EncryptBlock32(key, block)
}

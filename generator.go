package tf_random

import (
	"math"

	random "github.com/susisu/go-random/uint32"
	tf "github.com/susisu/go-tf-random/internal/threefish"
)

// TFGen is a pseudorandom number generator implementation that utilizes Threefish block cipher.
// It also has the capability to be split into multiple independent generators, which can be safely
// passed to other goroutines.
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

// Clone creates a copy of the generator.
func (g *TFGen) Clone() *TFGen {
	return &TFGen{
		key:        g.key,
		count:      g.count,
		bits:       g.bits,
		bitsIndex:  g.bitsIndex,
		block:      g.block,
		blockIndex: g.blockIndex,
		stale:      g.stale,
	}
}

// Uint32 generates a random uint32 value.
func (g *TFGen) Uint32() uint32 {
	if g.stale {
		panic("invalid call of Uint32: you cannot use a stale generator")
	}

	v := g.block[g.blockIndex]

	if g.count < math.MaxUint64-1 {
		if g.blockIndex == 8-1 {
			// cannot read from `block` anymore
			g.remake(g.key, g.count+1, g.bits, g.bitsIndex)
		} else {
			g.renew(g.key, g.count+1, g.bits, g.bitsIndex, g.block, g.blockIndex+1)
		}
	} else if g.bitsIndex < 64 {
		g.remake(g.key, 0, g.bits|(1<<g.bitsIndex), g.bitsIndex+1)
	} else {
		newKey := mash64(g.key, math.MaxUint64, g.bits)
		g.remake(newKey, 0, 0, 0)
	}

	return v
}

// Split creates a new generator that is independent from the original one.
// The new generator can be safely passed to another goroutine.
func (g *TFGen) Split() *TFGen {
	if g.stale {
		panic("invalid call of Split: you cannot use a stale generator")
	}

	if g.bitsIndex == 64 {
		newKey := mash64(g.key, g.count, g.bits)
		sg := make(newKey, 0, 1, 1)
		g.remake(newKey, 0, 0, 1)
		return sg
	} else {
		sg := make(g.key, g.count, g.bits|(1<<g.bitsIndex), g.bitsIndex+1)
		g.remake(g.key, g.count, g.bits, g.bitsIndex+1)
		return sg
	}
}

// Level flushes the internal state of the generator.
// Calling this method before performing multiple `splitn` operations may reduce the total number of
// computations.
func (g *TFGen) Level() {
	if g.stale {
		panic("invalid call of Level: you cannot use a stale generator")
	}

	if g.bitsIndex+40 > 64 {
		newKey := mash64(g.key, g.count, g.bits)
		g.remake(newKey, 0, 0, 0)
	}
}

// SplitN splits the generator into 2**nbits new generators and returns the i-th of them.
// After calling this method, the original generator is marked as "stale" and can no longer be used,
// except for calling SplitN to obtain other child generators.
// It panics if nbits is greater than 32.
func (g *TFGen) SplitN(nbits, i uint) *TFGen {
	if nbits > 32 {
		panic("invalid argument to SplitN: nbits must be less than or equal to 32")
	}

	g.stale = true

	b := uint64((math.MaxUint32 >> (32 - nbits)) & i)
	if g.bitsIndex+nbits > 64 {
		newKey := mash64(g.key, g.count, g.bits|(b<<g.bitsIndex))
		return make(newKey, 0, b>>(64-g.bitsIndex), nbits-(64-g.bitsIndex))
	} else {
		return make(g.key, g.count, g.bits|(b<<g.bitsIndex), g.bitsIndex+nbits)
	}
}

func mash64(key tf.Uint64x4, count, bits uint64) tf.Uint64x4 {
	block := []uint64{bits, count, 0, 0}
	return tf.Threefish256EncryptBlock64(key, block)
}

func mash32(key tf.Uint64x4, count, bits uint64) tf.Uint32x8 {
	block := []uint64{bits, count, 0, 0}
	return tf.Threefish256EncryptBlock32(key, block)
}

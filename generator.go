package tf_random

import (
	random "github.com/susisu/go-random/uint32"
	tf "github.com/susisu/go-tf-random/internal/threefish"
)

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

func NewTFGen(a, b, c, d uint64) *TFGen {
	return &TFGen{}
}

func (g *TFGen) Uint32() uint32 {
	if g.stale {
		panic("invalid call of Uint32: you cannot use a stale generator")
	}
	panic("not implemented")
}

func (g *TFGen) Split() *TFGen {
	if g.stale {
		panic("invalid call of Split: you cannot use a stale generator")
	}
	panic("not implemented")
}
func (g *TFGen) Flush() {
	if g.stale {
		panic("invalid call of Flush: you cannot use a stale generator")
	}
	panic("not implemented")
}

func (g *TFGen) SplitN(nbits, i uint) *TFGen {
	g.stale = true
	panic("not implemented")
}

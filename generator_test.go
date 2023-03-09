package tf_random_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	tf_random "github.com/susisu/go-tf-random"
)

func initTFGen() *tf_random.TFGen {
	g := tf_random.NewTFGen(
		0x00000000_00000000,
		0x89abcdef_01234567,
		0x01234567_89abcdef,
		0xffffffff_ffffffff,
	)
	initCount := 100
	for i := 0; i < initCount; i++ {
		g.Uint32()
	}
	return g
}

func TestTFGen_Uint32(t *testing.T) {
	g := initTFGen()
	sampleCount := 100
	seq := make([]uint32, 0, sampleCount)
	for i := 0; i < sampleCount; i++ {
		seq = append(seq, g.Uint32())
	}
	snaps.MatchSnapshot(t, seq)
}

func TestTFGen_Split(t *testing.T) {
	g1 := initTFGen()
	g2 := g1.Split()
	sampleCount := 100

	seq1 := make([]uint32, 0, sampleCount)
	for i := 0; i < sampleCount; i++ {
		seq1 = append(seq1, g1.Uint32())
	}
	snaps.MatchSnapshot(t, seq1)

	seq2 := make([]uint32, 0, sampleCount)
	for i := 0; i < sampleCount; i++ {
		seq2 = append(seq2, g2.Uint32())
	}
	snaps.MatchSnapshot(t, seq2)
}

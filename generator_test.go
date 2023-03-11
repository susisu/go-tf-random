package tf_random_test

import (
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
	tf_random "github.com/susisu/go-tf-random"
)

func TestMain(t *testing.M) {
	v := t.Run()
	snaps.Clean(t)
	os.Exit(v)
}

func initTFGen() *tf_random.TFGen {
	g := tf_random.NewTFGen(
		0x00000000_00000000,
		0x89abcdef_01234567,
		0x01234567_89abcdef,
		0xffffffff_ffffffff,
	)
	return g
}

func progressTFGen(g *tf_random.TFGen, n int) {
	for i := 0; i < n; i++ {
		g.Uint32()
	}
}

func testSnapshot(t *testing.T, g *tf_random.TFGen) {
	numSamples := 100
	seq := make([]uint32, 0, numSamples)
	for i := 0; i < numSamples; i++ {
		seq = append(seq, g.Uint32())
	}
	snaps.MatchSnapshot(t, seq)
}

func TestTFGen_Clone(t *testing.T) {
	g1 := initTFGen()
	progressTFGen(g1, 100)
	g2 := g1.Clone()
	assert.Equal(t, g1.Uint32(), g2.Uint32())
}

func TestTFGen_Uint32_snapshot(t *testing.T) {
	g := initTFGen()
	progressTFGen(g, 100)
	testSnapshot(t, g)
}

func TestTFGen_Split_snapshot(t *testing.T) {
	g1 := initTFGen()
	progressTFGen(g1, 100)
	g2 := g1.Split()
	testSnapshot(t, g1)
	testSnapshot(t, g2)
}

func TestTFGen_Level_snapshot(t *testing.T) {
	g1 := initTFGen()
	progressTFGen(g1, 100)
	g2 := g1.SplitN(32, 0xff)
	g2.Level()
	testSnapshot(t, g2)
}

func TestTFGen_SplitN_snapshot(t *testing.T) {
	g1 := initTFGen()
	progressTFGen(g1, 100)
	g2 := g1.SplitN(32, 0xff)
	testSnapshot(t, g2)
}

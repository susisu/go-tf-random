package tf_random_test

import (
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

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

func TestTFGen_Uint32(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		g := initTFGen()
		progressTFGen(g, 100)
		testSnapshot(t, g)
	})

	t.Run("distribution", func(t *testing.T) {
		testRng := rand.New(rand.NewSource(time.Now().UnixNano()))
		seed := []uint64{testRng.Uint64(), testRng.Uint64(), testRng.Uint64(), testRng.Uint64()}
		g := tf_random.NewTFGen(seed[0], seed[1], seed[2], seed[3])
		numBins := 8
		numSamplesPerBin := 2000
		numSamples := numBins * numSamplesPerBin

		histogram := make([]int, numBins)
		for i := 0; i < numSamples; i++ {
			v := g.Uint32()

			nv := float64(v) / float64(math.MaxUint32)
			i := int(math.Floor(nv * float64(numBins)))
			if i == numBins {
				i = numBins - 1
			}
			histogram[i]++
		}

		delta := 4 * math.Sqrt(float64(numSamplesPerBin)*(1.0-1.0/float64(numBins)))
		for i, c := range histogram {
			assert.InDeltaf(t, numSamplesPerBin, c, delta,
				"histogram(%d) = %d should be close to %d, (seed = %v)", i, c, numSamplesPerBin, seed)
		}
	})
}

func BenchmarkTFGen_Uint32(b *testing.B) {
	g := initTFGen()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Uint32()
	}
}

func BenchmarkMathRand(b *testing.B) {
	s := rand.NewSource(0xc0ffee)
	r := rand.New(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Uint32()
	}
}

func TestTFGen_Split(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		g1 := initTFGen()
		progressTFGen(g1, 100)
		g2 := g1.Split()
		testSnapshot(t, g1)
		testSnapshot(t, g2)
	})
}

func TestTFGen_Level(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		g1 := initTFGen()
		progressTFGen(g1, 100)
		g2 := g1.SplitN(32, 256)
		g2.Level()
		testSnapshot(t, g2)
	})
}

func TestTFGen_SplitN(t *testing.T) {
	t.Run("marks the generator as stale", func(t *testing.T) {
		g := initTFGen()
		_ = g.SplitN(32, 256)
		assert.Panics(t, func() { g.Uint32() })
		assert.Panics(t, func() { _ = g.Split() })
		assert.Panics(t, func() { g.Level() })
		assert.NotPanics(t, func() { _ = g.SplitN(32, 512) })
	})

	t.Run("panics if i >= 2**nbits", func(t *testing.T) {
		g := initTFGen()
		assert.Panics(t, func() { _ = g.SplitN(32, 1<<32) })
	})

	t.Run("snapshot", func(t *testing.T) {
		g1 := initTFGen()
		progressTFGen(g1, 100)
		g2 := g1.SplitN(32, 256)
		testSnapshot(t, g2)
	})
}

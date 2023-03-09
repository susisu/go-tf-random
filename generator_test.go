package tf_random_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	tf_random "github.com/susisu/go-tf-random"
)

func TestTFGen_Uint32(t *testing.T) {
	g := tf_random.NewTFGen(
		0x00000000_00000000,
		0x89abcdef_01234567,
		0x01234567_89abcdef,
		0xffffffff_ffffffff,
	)
	sampleCount := 100
	seq := make([]uint32, 0, sampleCount)
	for i := 0; i < sampleCount; i++ {
		seq = append(seq, g.Uint32())
	}
	snaps.MatchSnapshot(t, seq)
}

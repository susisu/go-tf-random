package threefish_test

import (
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	tf "github.com/susisu/go-tf-random/internal/threefish"
)

func TestMain(t *testing.M) {
	v := t.Run()
	snaps.Clean(t)
	os.Exit(v)
}

func TestThreefish256EncryptBlock64(t *testing.T) {
	key := []uint64{
		0x00000000_00000000,
		0x89abcdef_01234567,
		0x01234567_89abcdef,
		0xffffffff_ffffffff,
	}
	block := []uint64{
		0x00000000_00000000,
		0x89abcdef_01234567,
		0x01234567_89abcdef,
		0xffffffff_ffffffff,
	}
	out := tf.Threefish256EncryptBlock64(key, block)
	snaps.MatchSnapshot(t, out)
}

func TestThreefish256EncryptBlock32(t *testing.T) {
	key := []uint64{
		0x00000000_00000000,
		0x89abcdef_01234567,
		0x01234567_89abcdef,
		0xffffffff_ffffffff,
	}
	block := []uint64{
		0x00000000_00000000,
		0x89abcdef_01234567,
		0x01234567_89abcdef,
		0xffffffff_ffffffff,
	}
	out := tf.Threefish256EncryptBlock32(key, block)
	snaps.MatchSnapshot(t, out)
}

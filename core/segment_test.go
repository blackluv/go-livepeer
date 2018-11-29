package core

import (
	"bytes"
	"testing"

	"github.com/livepeer/lpms/ffmpeg"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func TestSegmentFlatten(t *testing.T) {
	s := SegmentMetadata{
		ManifestID: ManifestID("abcdef"),
		Seq:        1234,
		Hash:       ethcommon.BytesToHash(ethcommon.RightPadBytes([]byte("browns"), 32)),
		Profiles:   []ffmpeg.VideoProfile{ffmpeg.P144p30fps16x9, ffmpeg.P240p30fps16x9},
	}
	sHash := ethcommon.FromHex("e97461de03dcb5bf7f2e95c4ca9c99db2d049fb18a6df67dd9d557b2c05f6473")
	if !bytes.Equal(ethcrypto.Keccak256(s.Flatten()), sHash) {
		t.Error("Flattened segment + hash did not match expected hash")
	}
}

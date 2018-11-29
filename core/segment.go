package core

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/livepeer/go-livepeer/common"
	"github.com/livepeer/go-livepeer/net"

	"github.com/livepeer/lpms/ffmpeg"
)

type SegmentMetadata struct {
	ManifestID ManifestID
	Seq        int64
	Hash       ethcommon.Hash
	Profiles   []ffmpeg.VideoProfile
	OS         *net.OSInfo
}

func (s *SegmentMetadata) Flatten() []byte {
	profiles := common.ProfilesToHex(s.Profiles)
	seq := big.NewInt(s.Seq).Bytes()
	buf := make([]byte, len(s.ManifestID)+32+len(s.Hash.Bytes())+len(profiles))
	i := copy(buf[0:], []byte(s.ManifestID))
	i += copy(buf[i:], ethcommon.LeftPadBytes(seq, 32))
	i += copy(buf[i:], s.Hash.Bytes())
	i += copy(buf[i:], []byte(profiles))
	// i += copy(buf[i:], []byte(s.OS))
	return buf
}

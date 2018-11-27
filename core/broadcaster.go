package core

import (
	"errors"
	"fmt"

	"github.com/livepeer/go-livepeer/drivers"
	"github.com/livepeer/go-livepeer/net"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var ErrNotFound = errors.New("ErrNotFound")

// Broadcaster RPC interface implementation

type broadcaster struct {
	node  *LivepeerNode
	jobId string // ANGIE - DO WE GET RID OF JOBS HERE AS WELL?
	tinfo *net.TranscoderInfo
	ios   drivers.OSSession
	oos   drivers.OSSession
}

func (bcast *broadcaster) SetOrchestratorOS(ios drivers.OSSession) {
	bcast.ios = ios
}
func (bcast *broadcaster) GetOrchestratorOS() drivers.OSSession {
	return bcast.ios
}
func (bcast *broadcaster) SetBroadcasterOS(oos drivers.OSSession) {
	bcast.oos = oos
}
func (bcast *broadcaster) GetBroadcasterOS() drivers.OSSession {
	return bcast.oos
}
func (bcast *broadcaster) Sign(msg []byte) ([]byte, error) {
	if bcast.node == nil || bcast.node.Eth == nil {
		return []byte{}, fmt.Errorf("Cannot sign; missing eth client")
	}
	return bcast.node.Eth.Sign(crypto.Keccak256(msg))
}
func (bcast *broadcaster) JobId() string {
	return bcast.jobId
}
func (bcast *broadcaster) GetTranscoderInfo() *net.TranscoderInfo {
	return bcast.tinfo
}
func (bcast *broadcaster) SetTranscoderInfo(t *net.TranscoderInfo) {
	bcast.tinfo = t
}
func (bcast *broadcaster) Address() ethcommon.Address {
	if bcast.node == nil || bcast.node.Eth == nil {
		return ethcommon.Address{}
	}
	return bcast.node.Eth.Account().Address
}
func NewBroadcaster(node *LivepeerNode, jobId string) *broadcaster {
	return &broadcaster{
		node:  node,
		jobId: jobId,
	}
}

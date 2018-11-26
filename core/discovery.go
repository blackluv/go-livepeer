package core

import (
	"net/url"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/livepeer/go-livepeer/eth"
	lpTypes "github.com/livepeer/go-livepeer/eth/types"
	"github.com/livepeer/go-livepeer/net"
)

type OrchestratorSelector interface {
	GetOrchestrator(int, string) ([]*net.TranscoderInfo, error)
}

type offchainOrchestrator struct {
	address *url.URL
	client  eth.LivepeerEthClient
}

func (o *offchainOrchestrator) Address() *url.URL {
	return o.address
}

func NewOffchainOrchestrator(address *url.URL) *offchainOrchestrator {
	return &offchainOrchestrator{address: address}
}

func (o *offchainOrchestrator) GetOrchestrator(numOrchestrators int, orchAddr string) (*lpTypes.Transcoder, error) {
	address := ethcommon.HexToAddress(orchAddr)
	t, err := o.client.GetTranscoder(address)
	if err != nil {
		return nil, err
	}
	return t, nil
}

package core

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/livepeer/go-livepeer/eth"
	lpTypes "github.com/livepeer/go-livepeer/eth/types"
)

type OrchestratorSelector interface {
	GetOrchestrator(int, string) ([]*lpTypes.Transcoder, error)
}

type offchainOrchestrator struct {
	address string
	client  eth.LivepeerEthClient
}

func (o *offchainOrchestrator) Address() string {
	return o.address
}

func NewOffchainOrchestrator(address string) *offchainOrchestrator {
	return &offchainOrchestrator{address: address}
}

func (o *offchainOrchestrator) GetOrchestrator(numOrchestrators int) (*lpTypes.Transcoder, error) {
	address := ethcommon.HexToAddress(o.Address())
	t, err := o.client.GetTranscoder(address)
	if err != nil {
		return nil, err
	}
	return t, nil
}

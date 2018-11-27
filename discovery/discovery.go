package discovery

import (
	"net/url"

	"github.com/livepeer/go-livepeer/core"
	"github.com/livepeer/go-livepeer/net"
	"github.com/livepeer/go-livepeer/server"

	"github.com/golang/glog"
)

type offchainOrchestrator struct {
	uri  *url.URL
	node *core.LivepeerNode
}

func NewOffchainOrchestrator(node *core.LivepeerNode, address string) *offchainOrchestrator {
	uri, err := url.Parse(address)
	if err != nil {
		glog.Error("Could not parse orchestrator URI: ", err)
		return nil
	}
	return &offchainOrchestrator{node: node, uri: uri}
}

func (o *offchainOrchestrator) GetOrchestrators(numOrchestrators int) ([]*net.TranscoderInfo, error) {
	bcast := core.NewBroadcaster(o.node, "remove this param")
	tinfo, err := server.GetOrchestratorInfo(bcast, o.uri)
	return []*net.TranscoderInfo{tinfo}, err
}

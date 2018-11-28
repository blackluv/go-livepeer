package server

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net/url"
	"testing"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/livepeer/go-livepeer/core"
	"github.com/livepeer/go-livepeer/net"
)

type stubOrchestrator struct {
	priv    *ecdsa.PrivateKey
	block   *big.Int
	jobId   string
	signErr error
}

func StubJob() string {
	return "iamajobstring"
}

func (r *stubOrchestrator) ServiceURI() *url.URL {
	url, _ := url.Parse("http://localhost:1234")
	return url
}

func (r *stubOrchestrator) CurrentBlock() *big.Int {
	return r.block
}

func (r *stubOrchestrator) Sign(msg []byte) ([]byte, error) {
	if r.signErr != nil {
		return nil, r.signErr
	}
	hash := ethcrypto.Keccak256(msg)
	ethMsg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", 32, hash)
	return ethcrypto.Sign(ethcrypto.Keccak256([]byte(ethMsg)), r.priv)
}
func (r *stubOrchestrator) Address() ethcommon.Address {
	return ethcrypto.PubkeyToAddress(r.priv.PublicKey)
}
func (r *stubOrchestrator) TranscodeSeg(jobId int64, seg *core.SignedSegment) (*core.TranscodeResult, error) {
	return nil, nil
}
func (r *stubOrchestrator) StreamIDs(jobId string) ([]core.StreamID, error) {
	return []core.StreamID{}, nil
}

func StubOrchestrator() *stubOrchestrator {
	pk, err := ethcrypto.GenerateKey()
	if err != nil {
		return &stubOrchestrator{}
	}
	return &stubOrchestrator{priv: pk, block: big.NewInt(5), jobId: StubJob()}
}

func (r *stubOrchestrator) GetTranscoderInfo() *net.TranscoderInfo {
	return nil
}
func (r *stubOrchestrator) SetTranscoderInfo(ti *net.TranscoderInfo) {
}
func (r *stubOrchestrator) ServeTranscoder(stream net.Transcoder_RegisterTranscoderServer) {
}
func (r *stubOrchestrator) TranscoderResults(job int64, res *core.RemoteTranscoderResult) {
}
func (r *stubOrchestrator) TranscoderSecret() string {
	return ""
}
func StubBroadcaster2() *stubOrchestrator {
	return StubOrchestrator() // lazy; leverage subtyping for interface commonalities
}

func TestRPCTranscoderReq(t *testing.T) {

	o := StubOrchestrator()
	b := StubBroadcaster2()

	req, err := genTranscoderReq(b)
	if err != nil {
		t.Error("Unable to create transcoder req ", req)
	}
	if verifyTranscoderReq(o, req) != nil { // normal case
		t.Error("Unable to verify transcoder request")
	}

	// wrong broadcaster
	req.Address = ethcrypto.PubkeyToAddress(StubBroadcaster2().priv.PublicKey).Bytes()
	if verifyTranscoderReq(o, req) == nil {
		t.Error("Did not expect verification to pass; should mismatch broadcaster")
	}

	// invalid address
	req.Address = []byte("#non-hex address!")
	if verifyTranscoderReq(o, req) == nil {
		t.Error("Did not expect verification to pass; should mismatch broadcaster")
	}

	// error signing
	b.signErr = fmt.Errorf("Signing error")
	_, err = genTranscoderReq(b)
	if err == nil {
		t.Error("Did not expect to generate a transcoder request with invalid address")
	}
}

func TestRPCCreds(t *testing.T) {

	r := StubOrchestrator()
	jobId := StubJob()

	creds, err := genToken(r, jobId)
	if err != nil {
		t.Error("Unable to generate creds from req ", err)
	}
	if _, err := verifyToken(r, creds); err != nil {
		t.Error("Creds did not validate: ", err)
	}

	// // corrupt the creds
	// idx := len(creds) / 2
	// kreds := creds[:idx] + string(^creds[idx]) + creds[idx+1:]
	// if _, err := verifyToken(r, kreds); err == nil || err.Error() != "illegal base64 data at input byte 46" {
	// 	t.Error("Creds unexpectedly validated", err)
	// }

	// wrong orchestrator
	if _, err := verifyToken(StubOrchestrator(), creds); err == nil || err.Error() != "Token sig check failed" {
		t.Error("Orchestrator unexpectedly validated", err)
	}

	// // empty profiles
	// r.job.Profiles = []ffmpeg.VideoProfile{}
	// if _, err := verifyToken(r, creds); err.Error() != "Job out of range" {
	// 	t.Error("Unclaimable job unexpectedly validated", err)
	// }

	// // reset to sanity check once again
	// r.job = StubJob()
	// r.block = big.NewInt(0)
	// if _, err := verifyToken(r, creds); err != nil {
	// 	t.Error("Block did not validate", err)
	// }

}

func TestRPCSeg(t *testing.T) {
	b := StubBroadcaster2()
	o := StubOrchestrator()
	s := &BroadcastSession{
		Broadcaster: b,
	}

	baddr := ethcrypto.PubkeyToAddress(b.priv.PublicKey)

	jobId := StubJob()
	broadcasterAddress = baddr

	segData := &net.SegData{Seq: 4, Hash: ethcommon.RightPadBytes([]byte("browns"), 32)}

	creds, err := genSegCreds(s, segData)
	if err != nil {
		t.Error("Unable to generate seg creds ", err)
		return
	}
	if _, err := verifySegCreds(o, jobId, creds); err != nil {
		t.Error("Unable to verify seg creds", err)
		return
	}

	// test invalid jobid
	// oldSid := StreamId
	// StreamId = StreamId + StreamId
	// if _, err := verifySegCreds(o, jobId, creds); err == nil || err.Error() != "Segment sig check failed" {
	// 	t.Error("Unexpectedly verified seg creds: invalid jobid", err)
	// 	return
	// }
	// StreamId = oldSid

	// test invalid bcast addr
	oldAddr := broadcasterAddress
	key, _ := ethcrypto.GenerateKey()
	broadcasterAddress = ethcrypto.PubkeyToAddress(key.PublicKey)
	if _, err := verifySegCreds(o, jobId, creds); err == nil || err.Error() != "Segment sig check failed" {
		t.Error("Unexpectedly verified seg creds: invalid bcast addr", err)
	}
	broadcasterAddress = oldAddr

	// sanity check
	if _, err := verifySegCreds(o, jobId, creds); err != nil {
		t.Error("Sanity check failed", err)
	}

	// test corrupt creds
	idx := len(creds) / 2
	kreds := creds[:idx] + string(^creds[idx]) + creds[idx+1:]
	if _, err := verifySegCreds(o, jobId, kreds); err == nil || err.Error() != "illegal base64 data at input byte 70" {
		t.Error("Unexpectedly verified bad creds", err)
	}
}

func TestPing(t *testing.T) {
	o := StubOrchestrator()

	tsSignature, _ := o.Sign([]byte(fmt.Sprintf("%v", time.Now())))
	pingSent := crypto.Keccak256(tsSignature)
	req := &net.PingPong{Value: pingSent}

	pong, err := ping(context.Background(), req, o)
	if err != nil {
		t.Error("Unable to send Ping request")
	}

	verified := verifyMsgSig(o.Address(), string(pingSent), pong.Value)

	if !verified {
		t.Error("Unable to verify response from ping request")
	}
}

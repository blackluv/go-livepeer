package core

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var ErrStreamID = errors.New("ErrStreamID")
var ErrManifestID = errors.New("ErrManifestID")

const (
	HashLength = 32
)

func RandomVideoID() []byte {
	rand.Seed(time.Now().UnixNano())
	x := make([]byte, HashLength, HashLength)
	for i := 0; i < len(x); i++ {
		x[i] = byte(rand.Uint32())
	}
	return x
}

//StreamID is VideoID|Rendition
type StreamID string

func MakeStreamID(id []byte, rendition string) (StreamID, error) {
	if len(id) != HashLength || rendition == "" {
		return "", ErrStreamID
	}
	return StreamID(fmt.Sprintf("%x%v", id, rendition)), nil
}

func (id *StreamID) GetVideoID() []byte {
	if len(*id) < 2*HashLength {
		return nil
	}
	vid, err := hex.DecodeString(string((*id)[:2*HashLength]))
	if err != nil {
		return nil
	}
	return vid
}

func (id *StreamID) GetRendition() string {
	// XXX add tests
	if len(*id) < 2*HashLength {
		return ""
	}
	return string((*id)[2*HashLength:])
}

func (id *StreamID) IsValid() bool {
	return len(*id) > 2*HashLength
}

func (id StreamID) String() string {
	return string(id)
}

//ManifestID is VideoID
type ManifestID string

func MakeManifestID(id []byte) (ManifestID, error) {
	if len(id) != HashLength {
		return ManifestID(""), ErrManifestID
	}
	return ManifestID(fmt.Sprintf("%x", id)), nil
}

func (id *ManifestID) GetVideoID() []byte {
	if len(*id) < 2*HashLength {
		return nil
	}
	vid, err := hex.DecodeString(string((*id)[:2*HashLength]))
	if err != nil {
		return nil
	}
	return vid
}

func (id StreamID) ManifestIDFromStreamID() (ManifestID, error) {
	if !id.IsValid() {
		return "", ErrManifestID
	}
	mid, err := MakeManifestID(id.GetVideoID())
	return mid, err
}

func (id *ManifestID) IsValid() bool {
	return len(*id) == 2*HashLength
}

func (id ManifestID) String() string {
	return string(id)
}

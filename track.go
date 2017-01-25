package midi

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	kMTrk = 0x4d54726b
)

type Track struct {
	Events []Event
}

func ReadTrack(from io.Reader) (*Track, error) {
	var MTrk uint32
	if err := binary.Read(from, binary.BigEndian, &MTrk); err != nil {
		return nil, err
	}
	if MTrk != kMTrk {
		return nil, fmt.Errorf("Invalid track header; expected %08x, but %08x", kMTrk, MTrk)
	}
	var MTrkLength uint32
	if err := binary.Read(from, binary.BigEndian, &MTrkLength); err != nil {
		return nil, err
	}

	proxy := new(io.LimitedReader)
	proxy.R = from
	proxy.N = int64(MTrkLength)

	track := new(Track)
	var status byte
	var tick uint64 = 0
	for proxy.N > 0 {
		delta, err := readDeltaTime(proxy)
		if err != nil {
			return nil, err
		}
		event, next, err := ReadEvent(proxy, status)
		if err != nil {
			return nil, err
		}
		status = next
		tick += uint64(delta)
		event.Tick = tick
		track.Events = append(track.Events, *event)
	}

	discardAll(proxy)

	return track, nil
}

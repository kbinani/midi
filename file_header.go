package midi

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	kMThd = 0x4d546864
)

type FileHeader struct {
	Format     uint16
	NumTrack   uint16
	TimeFormat uint16
}

func ReadFileHeader(from io.Reader) (*FileHeader, error) {
	var MThd uint32
	if err := binary.Read(from, binary.BigEndian, &MThd); err != nil {
		return nil, err
	}
	if MThd != kMThd {
		return nil, fmt.Errorf("Invalid header; expected %d, but %d", kMThd, MThd)
	}

	var MThdLength uint32
	if err := binary.Read(from, binary.BigEndian, &MThdLength); err != nil {
		return nil, err
	}
	proxy := io.LimitReader(from, int64(MThdLength))

	header := new(FileHeader)
	if err := binary.Read(proxy, binary.BigEndian, &header.Format); err != nil {
		return nil, err
	}
	if err := binary.Read(proxy, binary.BigEndian, &header.NumTrack); err != nil {
		return nil, err
	}
	if err := binary.Read(proxy, binary.BigEndian, &header.TimeFormat); err != nil {
		return nil, err
	}

	discardAll(proxy)

	return header, nil
}

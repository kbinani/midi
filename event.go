package midi

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

type Event struct {
	Tick     uint64
	Messages []byte
}

func ReadEvent(from io.Reader, status byte) (e *Event, nextStatus byte, err error) {
	var b byte
	if err := binary.Read(from, binary.BigEndian, &b); err != nil {
		return nil, 0, err
	}
	e = new(Event)
	if b < 0x80 {
		e.Messages = append(e.Messages, status)
	}
	e.Messages = append(e.Messages, b)
	status = e.Messages[0]
	control := status & 0xF0

	if control == 0x80 || control == 0x90 || control == 0xA0 || control == 0xB0 || control == 0xE0 || status == 0xF2 {
		for len(e.Messages) < 3 {
			if err := binary.Read(from, binary.BigEndian, &b); err != nil {
				return nil, 0, err
			}
			e.Messages = append(e.Messages, b)
		}
		return e, status, nil
	} else if control == 0xC0 || control == 0xD0 || status == 0xF1 || status == 0xF3 {
		for len(e.Messages) < 2 {
			if err := binary.Read(from, binary.BigEndian, &b); err != nil {
				return nil, 0, err
			}
			e.Messages = append(e.Messages, b)
		}
		return e, status, nil
	} else if status == 0xF6 {
		return e, status, nil
	} else if status == 0xFF {
		var metaEventType byte
		if err := binary.Read(from, binary.BigEndian, &metaEventType); err != nil {
			return nil, 0, err
		}
		e.Messages = append(e.Messages, metaEventType)
		length, err := readDeltaTime(from)
		if err != nil {
			return nil, 0, err
		}
		proxy := io.LimitReader(from, int64(length))
		all, err := ioutil.ReadAll(proxy)
		if err != nil {
			return nil, 0, err
		}
		e.Messages = append(e.Messages, all...)
		return e, status, nil
	} else if status == 0xF0 || status == 0xF7 {
		var length uint32
		if err := binary.Read(from, binary.BigEndian, &length); err != nil {
			return nil, 0, err
		}
		if status == 0xF0 {
			length += 1
		}
		proxy := io.LimitReader(from, int64(length))
		all, err := ioutil.ReadAll(proxy)
		if err != nil {
			return nil, 0, err
		}
		e.Messages = append(e.Messages, all...)
		return e, status, nil
	}
	return nil, 0, fmt.Errorf("Cannot handle status: 0x%x", status)
}

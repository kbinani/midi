package midi

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"unsafe"
)

func discardAll(r io.Reader) {
	io.Copy(ioutil.Discard, r)
}

func readDeltaTime(r io.Reader) (uint32, error) {
	var v uint32
	var shifted uintptr = 0
	for true {
		if shifted+7 >= unsafe.Sizeof(v)*8 {
			return 0, fmt.Errorf("variable length time-delta is too large")
		}
		var ch byte
		if err := binary.Read(r, binary.BigEndian, &ch); err != nil {
			return 0, err
		}
		v = (v << 7) | uint32(0x7f&ch)
		shifted += 7
		if (ch & 0x80) == 0 {
			break
		}
	}
	return v, nil
}

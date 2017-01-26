midi
====

[![Build Status](https://travis-ci.org/kbinani/midi.svg?branch=master)](https://travis-ci.org/kbinani/midi)
[![codecov.io](https://codecov.io/github/kbinani/midi/branch/master/graph/badge.svg)](https://codecov.io/github/kbinani/midi)
[![](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/kbinani/midi)
[![](https://img.shields.io/badge/license-MIT-428F7E.svg?style=flat)](https://github.com/kbinani/midi/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/kbinani/midi)](https://goreportcard.com/report/github.com/kbinani/midi)

* Go library to parse Standard MIDI file.

todo
====
- [ ] Extract tempo change table
- [ ] Extract time signature table
- [ ] Writing SMF file

example
=======

```go
package main

import (
	"fmt"
	"os"

	"github.com/kbinani/midi"
)

func main() {
	f, err := os.Open("test.mid")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	file, err := midi.Read(f)
	if err != nil {
		panic(err)
	}
	for i, track := range file.Tracks {
		fmt.Printf("track#%d: %5d events\n", i, len(track.Events))
	}
}
```

license
=======

MIT License

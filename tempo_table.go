package midi

import (
	"math/big"
	"sort"
)

const (
	kDefaultTempo MicroSecondPerBeat = 500000
)

type TempoTable struct {
	table sliceOfTempo
}

type sliceOfTempo []*Tempo

func (v sliceOfTempo) Len() int {
	return len(v)
}

func (v sliceOfTempo) Less(i, j int) bool {
	return v[i].tick < v[j].tick
}

func (v sliceOfTempo) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v *TempoTable) Get(i int) Tempo {
	t := v.table[i]
	return *t
}

func (v *TempoTable) Set(i int, tick Tick, tempo MicroSecondPerBeat) Tempo {
	t := v.table[i]
	t.tick = tick
	t.tempo = tempo
	v.update()
	return *t
}

func (v *TempoTable) Append(tick Tick, tempo MicroSecondPerBeat) Tempo {
	t := new(Tempo)
	t.tick = tick
	t.tempo = tempo
	v.table = append(v.table, t)
	v.update()
	return *t
}

func (v *TempoTable) Delete(i int) {
	v.table = append(v.table[0:i], v.table[(i+1):]...)
	v.update()
}

func (v *TempoTable) Size() int {
	return len(v.table)
}

func (v *TempoTable) update() {
	if len(v.table) == 0 {
		return
	}
	if len(v.table) > 1 {
		sort.Stable(v.table)
	}

	t0 := v.table[0]
	if t0.tick == 0 {
		t0.time = *big.NewRat(0, 1)
	} else {
		t0.time.Mul(big.NewRat(int64(t0.tick), 480), big.NewRat(int64(kDefaultTempo), 1000000))
	}

	for i := 1; i < len(v.table); i++ {
		t := v.table[i]
		deltaTick := t.tick - t0.tick
		deltaSec := new(big.Rat).Mul(big.NewRat(int64(deltaTick), 480), big.NewRat(int64(t0.tempo), 1000000))
		t.time.Add(t0.Time(), deltaSec)
		t0 = t
	}
}

func NewTempoTable(track Track) *TempoTable {
	t := new(TempoTable)
	for _, e := range track.Events {
		if len(e.Messages) < 5 {
			continue
		}
		if e.Messages[0] != 0xFF || e.Messages[1] != 0x51 {
			continue
		}
		tempo := (uint32(e.Messages[2]) << 16) | (uint32(e.Messages[3]) << 8) | (uint32(e.Messages[4]))
		t.Append(e.Tick, MicroSecondPerBeat(tempo))
	}
	return t
}

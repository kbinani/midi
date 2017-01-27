package midi

import (
	"fmt"
	"math/big"
	"sort"
)

const (
	kDefaultTempo MicroSecondPerBeat = 500000 // 120 bps
	kMega         int64              = 1000000
)

type TempoTable struct {
	table sliceOfTempo
	tpqn  int64 // tick per quater note
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

func (v *TempoTable) AppendAll(tickAndTempo map[Tick]MicroSecondPerBeat) {
	for tick, tempo := range tickAndTempo {
		t := new(Tempo)
		t.tick = tick
		t.tempo = tempo
		v.table = append(v.table, t)
	}
	v.update()
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
		t0.time = new(big.Rat)
	} else {
		t0.time = big.NewRat(int64(t0.tick), v.tpqn)
		t0.time.Mul(t0.time, big.NewRat(int64(kDefaultTempo), kMega))
	}

	for i := 1; i < len(v.table); i++ {
		t := v.table[i]
		deltaTick := t.tick - t0.tick
		deltaSec := big.NewRat(int64(deltaTick), v.tpqn)
		deltaSec.Mul(deltaSec, big.NewRat(int64(t0.tempo), kMega))
		t.time = new(big.Rat)
		t.time.Add(t.time, deltaSec)
		t0 = t
	}
}

func NewTempoTable(file *File, trackIndex int) *TempoTable {
	track := file.Tracks[trackIndex]

	t := new(TempoTable)
	t.tpqn = int64(file.Header.TimeFormat)
	tickAndTempo := make(map[Tick]MicroSecondPerBeat)
	for _, e := range track.Events {
		if len(e.Messages) < 5 {
			continue
		}
		if e.Messages[0] != 0xFF || e.Messages[1] != 0x51 {
			continue
		}
		tempo := (uint32(e.Messages[2]) << 16) | (uint32(e.Messages[3]) << 8) | (uint32(e.Messages[4]))
		tickAndTempo[e.Tick] = MicroSecondPerBeat(tempo)
		if len(tickAndTempo) > 256 {
			t.AppendAll(tickAndTempo)
			tickAndTempo = make(map[Tick]MicroSecondPerBeat)
		}
	}
	t.AppendAll(tickAndTempo)
	return t
}

func (v *TempoTable) SecFromTick(tick *big.Rat) *big.Rat {
	i := sort.Search(len(v.table), func(i int) bool {
		return tick.Cmp(big.NewRat(int64(v.table[i].tick), 1)) < 0
	})

	tempo := kDefaultTempo
	time := new(big.Rat)
	start := new(big.Rat)
	if i-1 >= 0 {
		t := v.table[i-1]
		tempo = t.tempo
		time = t.Sec()
		start = big.NewRat(int64(t.tick), 1)
	}

	deltaBeat := new(big.Rat)
	deltaBeat.Set(tick)
	deltaBeat.Sub(deltaBeat, start)
	deltaBeat.Mul(deltaBeat, big.NewRat(1, v.tpqn))

	vv := new(big.Rat)
	vv.Add(time, deltaBeat.Mul(deltaBeat, big.NewRat(int64(tempo), kMega)))

	return vv
}

func (v *TempoTable) TickFromSec(sec *big.Rat) *big.Rat {
	i := sort.Search(len(v.table), func(i int) bool {
		return sec.Cmp(v.table[i].Sec()) < 0
	})

	tempo := kDefaultTempo
	time := new(big.Rat)
	start := new(big.Rat)
	if i-1 >= 0 {
		t := v.table[i-1]
		tempo = t.tempo
		time = t.Sec()
		start = big.NewRat(int64(t.tick), 1)
	}

	deltaSec := new(big.Rat)
	deltaSec.Set(sec)
	deltaSec.Sub(deltaSec, time)

	deltaBeat := new(big.Rat)
	deltaBeat.Set(deltaSec)
	deltaBeat.Mul(deltaBeat, big.NewRat(kMega, int64(tempo)))

	deltaTick := new(big.Rat)
	deltaTick.Set(deltaBeat)
	deltaTick.Mul(deltaTick, big.NewRat(v.tpqn, 1))

	start.Add(start, deltaTick)

	return start
}

func (v *TempoTable) Debug() {
	fmt.Printf("------------------\n")
	for i, t := range v.table {
		fmt.Printf("#%3d %5d %s\n", i, t.tick, t.time.String())
	}
}

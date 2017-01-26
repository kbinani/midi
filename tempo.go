package midi

import "math/big"

type MicroSecondPerBeat uint32

type Tempo struct {
	tick  Tick
	tempo MicroSecondPerBeat
	time  big.Rat
}

func (v Tempo) Tick() Tick {
	return v.tick
}

func (v Tempo) MicroSecondPerBeat() MicroSecondPerBeat {
	return v.tempo
}

func (v Tempo) Time() *big.Rat {
	t := v.time
	return &t
}

func (v Tempo) TimeF() float64 {
	r, _ := v.time.Float64()
	return r
}

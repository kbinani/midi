package midi

import "math/big"

type MicroSecondPerBeat uint32

type Tempo struct {
	tick  Tick
	tempo MicroSecondPerBeat
	time  *big.Rat
}

func (v Tempo) Tick() Tick {
	return v.tick
}

func (v Tempo) MicroSecondPerBeat() MicroSecondPerBeat {
	return v.tempo
}

func (v Tempo) Sec() *big.Rat {
	r := new(big.Rat)
	r.Set(v.time)
	return r
}

func (v Tempo) FSec() float64 {
	f, _ := v.time.Float64()
	return f
}

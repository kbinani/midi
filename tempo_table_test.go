package midi

import (
	"math/big"
	"os"
	"testing"
)

func TestTempoTable(t *testing.T) {
	list := new(TempoTable)
	list.tpqn = 480

	if list.Size() != 0 {
		t.Errorf("Size() should be zero after 'new'")
	}
	list.Append(480, 525000)
	list.Append(0, 500000)

	if list.Size() != 2 {
		t.Errorf("Size should be 2, because Append called 2 times")
	}

	v0 := list.Get(0)
	v1 := list.Get(1)

	if v0.Tick() != 0 || v0.MicroSecondPerBeat() != 500000 {
		t.Errorf("#1 tick, tempo are unexpectedly changed")
	}
	if v1.Tick() != 480 || v1.MicroSecondPerBeat() != 525000 {
		t.Errorf("#2 tick, tempo are unexpectedly changed")
	}
	if v0.Sec().Cmp(big.NewRat(0, 1)) != 0 {
		t.Errorf("#1 time wrong: expected 0 but %v", v0.Sec().String())
	}
	if v1.Sec().Cmp(big.NewRat(1, 2)) != 0 {
		t.Errorf("#2 time wrong: expected 1/2 but %v", v1.Sec().String())
	}
	if v1.FSec() != 0.5 {
		t.Errorf("#2 time wrong: expected 0.5 but %f", v1.FSec())
	}

	if secAt480 := list.SecFromTick(big.NewRat(480, 1)); secAt480.String() != "1/2" {
		t.Errorf("time(sec) at tick=480 should be 1/2 but %v", secAt480.String())
	}
	if secAt960 := list.SecFromTick(big.NewRat(960, 1)); secAt960.String() != "41/40" { // 1.025
		t.Errorf("time(sec) at tick=960 should be 41/40 but %v", secAt960.String())
	}
	if tickAt1_2 := list.TickFromSec(big.NewRat(1, 2)); tickAt1_2.String() != "480/1" {
		t.Errorf("tick at t=0.5(sec) should be 480/1, but %v", tickAt1_2.String())
	}
	if tickAt41_40 := list.TickFromSec(big.NewRat(41, 40)); tickAt41_40.String() != "960/1" {
		t.Errorf("tick at t=1.025(sec) should be 960/1, but %v", tickAt41_40)
	}

	list.Set(1, 520, 600000)
	v1 = list.Get(1)
	if v1.Tick() != 520 || v1.MicroSecondPerBeat() != 600000 {
		t.Errorf("#2 tick, tempo was not changed")
	}
	if v1.Sec().Cmp(big.NewRat(13, 24)) != 0 {
		t.Errorf("#2 time wrong: expected 1/2 but %v", v1.Sec().String())
	}

	list.Set(1, 480, 525000)

	list.Delete(0)
	if list.Size() != 1 {
		t.Errorf("list.Size() didn't decrease after Delete")
	}

	v0 = list.Get(0)
	if v0.Tick() != 480 || v0.MicroSecondPerBeat() != 525000 {
		t.Errorf("#1 tick, tempo are wrong")
	}
	if t0 := v0.Sec(); t0.Cmp(big.NewRat(1, 2)) != 0 {
		t.Errorf("#1 time is wrong")
	}

	list.Delete(0)
	if list.Size() != 0 {
		t.Errorf("list.Size() didn't decrease after Delete")
	}
}

func TestNewTempoTable(t *testing.T) {
	fp, err := os.Open("testdata/tempo_table_test.mid")
	if err != nil {
		t.Errorf("cannot read test file: %v", err)
	}
	defer fp.Close()

	file, err := Read(fp)
	if err != nil {
		t.Errorf("cannot read testdata")
	}

	table := NewTempoTable(file, 0)
	if table.Size() != 2 {
		t.Errorf("number of tempo change: expected 2 but %d", table.Size())
	}
	t0 := table.Get(0)
	if t0.Tick() != 0 || t0.MicroSecondPerBeat() != 500000 || t0.FSec() != 0.0 {
		t.Errorf("#1 tempo tick, tempo, time mismatch")
	}
	t1 := table.Get(1)
	if t1.Tick() != 1920 || t1.MicroSecondPerBeat() != 250000 || t1.FSec() != 2.0 {
		t.Errorf("#2 tempo tick, tempo, time mismatch")
	}
}

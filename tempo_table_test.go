package midi

import (
	"math/big"
	"os"
	"testing"
)

func TestTempoTable(t *testing.T) {
	list := new(TempoTable)

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
	if v0.Time().Cmp(big.NewRat(0, 1)) != 0 {
		t.Errorf("#1 time wrong: expected 0 but %v", v0.Time().RatString())
	}
	if v1.Time().Cmp(big.NewRat(1, 2)) != 0 {
		t.Errorf("#2 time wrong: expected 1/2 but %v", v1.Time().RatString())
	}
	if v1.TimeF() != 0.5 {
		t.Errorf("#2 time wrong: expected 0.5 but %f", v1.TimeF())
	}

	list.Set(1, 520, 600000)
	v1 = list.Get(1)
	if v1.Tick() != 520 || v1.MicroSecondPerBeat() != 600000 {
		t.Errorf("#2 tick, tempo was not changed")
	}
	if v1.Time().Cmp(big.NewRat(13, 24)) != 0 {
		t.Errorf("#2 time wrong: expected 1/2 but %v", v1.Time().RatString())
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
	if t0 := v0.Time(); t0.Cmp(big.NewRat(1, 2)) != 0 {
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

	table := NewTempoTable(file.Tracks[0])
	if table.Size() != 2 {
		t.Errorf("number of tempo change: expected 2 but %d", table.Size())
	}
	t0 := table.Get(0)
	if t0.Tick() != 0 || t0.MicroSecondPerBeat() != 500000 || t0.TimeF() != 0.0 {
		t.Errorf("#1 tempo tick, tempo, time mismatch")
	}
	t1 := table.Get(1)
	if t1.Tick() != 1920 || t1.MicroSecondPerBeat() != 250000 || t1.TimeF() != 2.0 {
		t.Errorf("#2 tempo tick, tempo, time mismatch")
	}
}

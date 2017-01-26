package midi

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	fp, err := os.Open("file_test.mid")
	if err != nil {
		t.Errorf("cannot read test file: %v", err)
	}
	defer fp.Close()

	file, err := Read(fp)
	if err != nil {
		t.Errorf("cannot create File object: %v", err)
	}

	expectedNumEvents := []int{
		5, 149, 787, 27, 123,
	}

	if len(file.Tracks) != len(expectedNumEvents) {
		t.Errorf("track number mismatch: expected %d but %d", len(expectedNumEvents), len(file.Tracks))
	}

	for i, expectedNumEvent := range expectedNumEvents {
		track := file.Tracks[i]
		if len(track.Events) != expectedNumEvent {
			t.Errorf("event number mismatch: expected %d but %d", expectedNumEvent, len(track.Events))
		}
	}
}

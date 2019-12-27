package mdturbo

import (
	"fmt"
	"testing"
)

func TestDeserialize(t *testing.T) {
	var partitionChecks = []struct {
		number, start, length uint32
	}{
		{0, 256, 65535},
		{1, 65791, 65535},
		{2, 131326, 65535},
		{3, 196861, 65535},
		{4, 262396, 65535},
		{5, 327931, 65535},
		{6, 393466, 65535},
		{7, 459001, 65535},
	}

	partmap, err := Deserialize(testData)
	if err != nil {
		t.Errorf("correctly-sized partition table failed to deserialize")
	}
	if !partmap.Validate() {
		t.Errorf("valid partition table failed to validate")
	}
	if partmap.Cylinders != 995 {
		t.Errorf("cylinder count incorrect (got %d, wanted %d)", partmap.Cylinders, 995)
	}
	if partmap.Heads != 16 {
		t.Errorf("head count incorrect (got %d, wanted %d)", partmap.Heads, 16)
	}
	if partmap.Sectors != 63 {
		t.Errorf("sector count incorrect (got %d, wanted %d)", partmap.Sectors, 63)
	}
	if partmap.PartCount1 != 8 {
		t.Errorf("partition 1 count incorrect (got %d, wanted %d)", partmap.PartCount1, 8)
	}
	if partmap.PartCount2 != 0 {
		t.Errorf("partition 2 count incorrect (got %d, wanted %d)", partmap.PartCount2, 0)
	}
	if partmap.RomVersion != 3 {
		t.Errorf("ROM version incorrect (got %d wanted %d)", partmap.RomVersion, 3)
	}

	for _, tt := range partitionChecks {
		testname := fmt.Sprintf("partcheck-%d", tt.number)
		t.Run(testname, func(t *testing.T) {
			start := partmap.Partitions1[tt.number].Start
			if start != tt.start {
				t.Errorf("got %d, wanted %d", start, tt.start)
			}
			length := partmap.Partitions1[tt.number].rawLength
			if length != tt.length {
				t.Errorf("got %d, wanted %d", start, tt.length)
			}
		})
	}

}

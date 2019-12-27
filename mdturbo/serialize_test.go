package mdturbo

import (
	"testing"
)

func TestSerializeEmpty(t *testing.T) {
	var partmap MDTurbo
	data, err := Serialize(partmap)
	if err != nil {
		t.Errorf("encountered error serializing empty partition table: %v", err)
	}
	for _, i := range data {
		if i != 0 {
			t.Errorf("empty partition table encoded un-empty")
		}
	}
}

func TestSerializeRoundTrip(t *testing.T) {
	// Get a valid partmap from the existing testData
	partmap, err := Deserialize(testData)
	if err != nil {
		t.Errorf("correctly-sized partition table failed to deserialize")
	}
	data, err := Serialize(partmap)
	if err != nil {
		t.Errorf("failed to serialize testData partmap: %v", err)
	}
	if data != testData {
		t.Errorf("test data failed to match existing data")
	}
}

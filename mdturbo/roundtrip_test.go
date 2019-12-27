package mdturbo

import (
	"fmt"
	"math/rand"
	"testing"
)

func RandomRoundTrip(t *testing.T) {
	// Get a probably-invalid partmap from random data
	var chonk [PartitionBlkLen]byte

	buf := make([]byte, PartitionBlkLen)
	count, err := rand.Read(buf)
	if err != nil {
		t.Errorf("could not generate random data: %v", err)
	}
	if count != PartitionBlkLen {
		t.Errorf("short random data: %d", count)
	}
	copy(chonk[:], buf)

	partmap, err := Deserialize(chonk)
	if err != nil {
		t.Errorf("could not deserialize random data: %v", err)
	}
	data, err := Serialize(partmap)
	if err != nil {
		t.Errorf("failed to serialize randomized partmap: %v", err)
	}
	if data != chonk {
		t.Errorf("test data failed to match existing data")
	}
}

func TestRandomRoundTrip(t *testing.T) {
	for i := 0; i <= 10; i++ {
		testname := fmt.Sprintf("random-roundtrip-%d", i)
		t.Run(testname, RandomRoundTrip)
	}
}

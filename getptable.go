package main

import (
	"fmt"
	"os"
)

import (
	"github.com/disappearinjon/microdrive/mdturbo"
)

// GetPartitionTable returns an MDTurbo data structure and an error when
// provided a filename.
func GetPartitionTable(filename string) (ptable mdturbo.MDTurbo, err error) {

	// Open the file passed in for reading
	imagefile, err := os.Open(filename)
	defer imagefile.Close()
	if err != nil {
		return
	}

	// Get first disk sector, where the partition table sits
	var firstSector [mdturbo.SectorSize]byte
	buf := make([]byte, mdturbo.SectorSize)
	read, err := imagefile.Read(buf)
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Read %v bytes\n", read)
	copy(firstSector[:], buf)

	// Parse the sector
	ptable, err = mdturbo.Deserialize(firstSector)
	if err != nil {
		return
	}
	if !ptable.Validate() {
		fmt.Fprintf(os.Stderr, "WARNING: partition map appears invalid\n")
	}

	return
}

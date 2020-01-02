// Package mdturbo provides the MicroDrive/Turbo partition map format,
// along with serializer and deserializer functions.
//
// The format is AFAIK undocumented, but the CiderPress source at
// https://github.com/fadden/ciderpress/blob/master/diskimg/MicroDrive.cpp
// contains a partial description.
package mdturbo

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

// MaxPartitions is the maximum number of partitions in an image
const MaxPartitions = 8

// SectorSize is the number of bytes per sector
const SectorSize = 512

// PartitionBlkLen is the number of bytes for a partition block
const PartitionBlkLen = 512

// Partition is data for a single partition in a partition map
type Partition struct {
	Start uint32 // Offset in bytes of partition in sectors

	// RawLength may not reflect the "actual" partition length, but
	// we need it for correct round-tripping of arbitrary partition
	// tables.
	RawLength uint32 `json:"Length"` // Length of partition in sectors
}

// Length returns the actual partition length, with the for-dual-CF
// bytes masked out, as per the CiderPress documentation.
func (p Partition) Length() uint32 {
	return p.RawLength & 0x00ffffff
}

// End returns the last sector number of a partition
func (p Partition) End() uint32 {
	return p.Start + p.Length() - 1
}

// String returns a string representation of the partition details
// Format is start, end, size in kilobytes (tab-separated)
func (p Partition) String() string {
	return fmt.Sprintf("%d\t%d\t%d", p.Start, p.End(), p.Length()/2)
}

// MDTurbo is the data structure with what we know about a
// MicroDrive/Turbo partition map. The struct tag serialize is used to
// encode the offset of a field; length is calculated from the data type.
//
// Because the Partitions data structure doesn't map exactly, we do keep
// the offset but the prcessing is a very very special case.
type MDTurbo struct {
	Magic      uint16    `offset:"0x00"` // Drive type identifier
	Cylinders  uint16    `offset:"0x02"` // # of cylinders
	Unknown1   [2]uint8  `offset:"0x04"` // Unknown region 1
	Heads      uint16    `offset:"0x06"` // heads per cyl
	Sectors    uint16    `offset:"0x08"` // sectors per track
	Unknown2   [2]uint8  `offset:"0x0A"` // Unknown region 2
	PartCount1 uint8     `offset:"0x0C"` // # of partitions in first chunk
	PartCount2 uint8     `offset:"0x0D"` // # of partitions in second chunk
	Unknown3   [10]uint8 `offset:"0x0E"` // Unknown region 3
	RomVersion uint16    `offset:"0x18"` // IIgs Rom version (01 or 03)
	BootPart   uint16    `offset:"0x1A"` // Boot partition.
	Unknown4   [4]uint8  `offset:"0x1C"` // Unknown region 4

	// The partitions are actually represented as Start Sector
	// numbers (0x20, 4 bytes for each of the 8 = 32 bytes),
	// followed by lengths (another 32 bytes starting at 0x40).
	Partitions1 [MaxPartitions]Partition `offset:"0x20"`

	Unknown5 [32]uint8 `offset:"0x60"` // Unknown region 5

	// Same as Partitions1 but starting at 0x80 and 0xA0
	Partitions2 [MaxPartitions]Partition `offset:"0x80"`
	// Probably padding, per CiderPress
	Unknown6 [320]uint8 `offset:"0xC0"`
}

// Validate returns true if the partition tables appears to be valid,
// and False if it does not.
func (pt MDTurbo) Validate() bool {
	if pt.Magic != 52426 {
		return false
	}
	if pt.Partitions1[0].Start != 256 {
		return false
	}
	return true
}

// PartCount returns the total number of partitions in a partition
// table.
func (pt MDTurbo) PartCount() uint8 {
	return pt.PartCount1 + pt.PartCount2
}

// GetPartition returns the Partition data structure and an error
func (pt MDTurbo) GetPartition(partNum uint8) (Partition, error) {
	// Confirm the partition desired exists
	if partNum+1 > pt.PartCount() {
		return Partition{}, fmt.Errorf("partition %d does not exist (max %d)",
			partNum, pt.PartCount())
	}
	if partNum < MaxPartitions {
		return pt.Partitions1[partNum], nil
	} else {
		return pt.Partitions2[partNum-MaxPartitions], nil
	}
}

// Serialize is the struct-attached Serialize function, for convenience
func (pt MDTurbo) Serialize() ([PartitionBlkLen]uint8, error) {
	return Serialize(pt)
}

// String is the standard stringified partition printer
func (pt MDTurbo) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 6, 0, 4, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, PrettyPrint(pt))
	w.Flush()
	return buf.String()
}

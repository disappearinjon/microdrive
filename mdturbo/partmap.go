// Package mdturbo provides the MicroDrive/Turbo partition map format,
// along with serializer and deserializer functions.
//
// The format is AFAIK undocumented, but the CiderPress source at
// https://github.com/fadden/ciderpress/blob/master/diskimg/MicroDrive.cpp
// contains a partial description.
package mdturbo

// MaxPartitions is the maximum number of partitions in an image
const MaxPartitions = 8

// SectorSize is the number of bytes per sector
const SectorSize = 512

// PartitionBlkLen is the number of bytes for a partition block
const PartitionBlkLen = 512

// Partition is data for a single partition in a partition map
type Partition struct {
	Start uint32 // Offset in bytes of partition in sectors

	// rawLength may not reflect the "actual" partition length, but
	// we need it for correct round-tripping of arbitrary partition
	// tables.
	rawLength uint32 // Length of partition in sectors
}

func (p Partition) Length() uint32 {
	return p.rawLength & 0x00ffffff
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
	Unknown4   [6]uint8  `offset:"0x1A"` // Unknown region 4

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

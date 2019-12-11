// Package mdturbo provides the MicroDrive/Turbo partition map format,
// along with serializer and deserializer functions.
//
// The format is AFAIK undocumented, but the CiderPress source at
// https://github.com/fadden/ciderpress/blob/master/diskimg/MicroDrive.cpp
// contains a partial description.
package mdturbo

import "encoding/binary"
import "fmt"

// MaxPartitions is the maximum number of partitions in an image
const MaxPartitions = 8

// SectorSize is the number of bytes per sector
const SectorSize = 512

// PartitionBlkLen is the number of bytes for a partition block
const PartitionBlkLen = 256

// fieldNames is an enumerated type for the fields in the partition map
// data structure. See the descriptions in the MDTurbo struct definition
// for more detail
type fieldNames int

const (
	magic fieldNames = iota
	cylinders
	heads
	unknown1
	sectors
	unknown2
	partCount1
	partCount2
	unknown3
	romVersion
	unknown4
	partStart1
	partLen1
	unknown5
	partStart2
	partLen2
	unknown6
)

func (f fieldNames) String() string {
	return [...]string{"Magic", "Cylinders", "Heads", "Unknown1",
		"Sectors", "Unknown2", "PartCount1", "PartCount2",
		"Unknown3", "RomVersion", "Unknown4", "PartStart1",
		"PartLen1", "Unnknown5", "PartStart2", "PartLen2",
		"Unknown6"}[f]
}

// Partition is data for a single partition in a partition map
type Partition struct {
	Start  uint32 // Offset in bytes of partition in sectors
	Length uint32 // Length of partition in sectors
}

// MDTurbo is the data structure with what we know about a
// MicroDrive/Turbo partition map
type MDTurbo struct {
	Magic      uint16 // Drive type identifier
	Cylinders  uint16 // # of cylinders
	Heads      uint16 // heads per cyl
	Sectors    uint16 // sectors per track
	PartCount1 uint8  // # of partitions in first chunk
	PartCount2 uint8  // # of partitions in second chunk
	RomVersion uint16 // IIgs Rom version (01 or 03)

	Partitions1 [MaxPartitions]Partition // Partitions in first chunk
	Partitions2 [MaxPartitions]Partition // Partitions in second chunk

	// Unknown and Undocumented ranges
	Unknown1 []uint8 // Unknown region 1
	Unknown2 []uint8 // Unknown region 2
	Unknown3 []uint8 // Unknown region 3
	Unknown4 []uint8 // Unknown region 4
	Unknown5 []uint8 // Unknown region 5
	Unknown6 []uint8 // Probably padding, per CiderPress
}

// Field is the offset map for bytes in a sector
type Field struct {
	Start  uint16 // Byte offset
	Length uint16 // Length of field
}

// offsetMap is the byte offset for locations in the partition map
var offsetMap = map[fieldNames]Field{
	magic:      {0x00, 2},
	cylinders:  {0x02, 2},
	unknown1:   {0x04, 2},
	heads:      {0x06, 2},
	sectors:    {0x08, 2},
	unknown2:   {0x0A, 2},
	partCount1: {0x0C, 1},
	partCount2: {0x0D, 1},
	unknown3:   {0x0E, 10},
	romVersion: {0x18, 2},
	unknown4:   {0x1A, 6},
	partStart1: {0x20, 4}, // * number of MaxPartitions (i.e., 32 bytes total)
	partLen1:   {0x40, 4}, // * number of MaxPartitions
	unknown5:   {0x60, 32},
	partStart2: {0x80, 4}, // * number of MaxPartitions
	partLen2:   {0xa0, 4}, // * number of MaxPartitions
	unknown6:   {0xc0, 320},
}

// Deserialize converts from a disk sector into an MDTurbo struct
func Deserialize(data []byte) (MDTurbo, error) {
	var partmap MDTurbo

	// Basic sanity checks
	if len(data) < PartitionBlkLen {
		return partmap, fmt.Errorf("not an MDTurbo partition map: too short (%v bytes, expected %v)", len(data), PartitionBlkLen)
	}

	// Simple field deserialization
	partmap.Magic = binary.LittleEndian.Uint16(data[offsetMap[magic].Start : offsetMap[magic].Start+offsetMap[magic].Length])
	partmap.Cylinders = binary.LittleEndian.Uint16(data[offsetMap[cylinders].Start : offsetMap[cylinders].Start+offsetMap[cylinders].Length])
	partmap.Unknown1 = data[offsetMap[unknown1].Start : offsetMap[unknown1].Start+offsetMap[unknown1].Length]
	partmap.Heads = binary.LittleEndian.Uint16(data[offsetMap[heads].Start : offsetMap[heads].Start+offsetMap[heads].Length])
	partmap.Sectors = binary.LittleEndian.Uint16(data[offsetMap[sectors].Start : offsetMap[sectors].Start+offsetMap[sectors].Length])
	partmap.Unknown2 = data[offsetMap[unknown2].Start : offsetMap[unknown2].Start+offsetMap[unknown2].Length]
	partmap.PartCount1 = data[offsetMap[partCount1].Start]
	partmap.PartCount2 = data[offsetMap[partCount2].Start]
	partmap.Unknown3 = data[offsetMap[unknown3].Start : offsetMap[unknown3].Start+offsetMap[unknown3].Length]
	partmap.RomVersion = binary.LittleEndian.Uint16(data[offsetMap[romVersion].Start : offsetMap[romVersion].Start+offsetMap[romVersion].Length])
	partmap.Unknown4 = data[offsetMap[unknown4].Start : offsetMap[unknown4].Start+offsetMap[unknown4].Length]
	partmap.Unknown5 = data[offsetMap[unknown5].Start : offsetMap[unknown5].Start+offsetMap[unknown5].Length]
	partmap.Unknown6 = data[offsetMap[unknown6].Start : offsetMap[unknown6].Start+offsetMap[unknown6].Length]

	for partNum := 0; partNum < int(partmap.PartCount1); partNum++ {
		var length uint32
		startOffset := offsetMap[partStart1].Start + (uint16(partNum) * offsetMap[partStart1].Length)
		partmap.Partitions1[partNum].Start = binary.LittleEndian.Uint32(data[startOffset : startOffset+offsetMap[partStart1].Length])
		lengthOffset := offsetMap[partLen1].Start + (uint16(partNum) * offsetMap[partStart1].Length)
		length = binary.LittleEndian.Uint32(data[lengthOffset : lengthOffset+4])
		partmap.Partitions1[partNum].Length = length & 0x00ffffff
	}

	for partNum := 0; partNum < int(partmap.PartCount2); partNum++ {
		var length uint32
		startOffset := offsetMap[partStart2].Start + (uint16(partNum) * offsetMap[partStart2].Length)
		partmap.Partitions2[partNum].Start = binary.LittleEndian.Uint32(data[startOffset : startOffset+offsetMap[partStart2].Length])
		lengthOffset := offsetMap[partLen2].Start + (uint16(partNum) * offsetMap[partStart2].Length)
		length = binary.LittleEndian.Uint32(data[lengthOffset : lengthOffset+4])
		partmap.Partitions2[partNum].Length = length & 0x00ffffff
	}

	return partmap, nil
}

// Package mdturbo provides the MicroDrive/Turbo partition map format,
// along with serializer and deserializer functions.
//
// The format is AFAIK undocumented, but the CiderPress source at
// https://github.com/fadden/ciderpress/blob/master/diskimg/MicroDrive.cpp
// contains a partial description.
package mdturbo

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
)

// MaxPartitions is the maximum number of partitions in an image
const MaxPartitions = 8

// SectorSize is the number of bytes per sector
const SectorSize = 512

// PartitionBlkLen is the number of bytes for a partition block
const PartitionBlkLen = 512

// Partition is data for a single partition in a partition map
type Partition struct {
	Start     uint32 // Offset in bytes of partition in sectors
	RawLength uint32 // Length of partition in sectors
}

func (p Partition) Length() uint32 {
	return p.RawLength & 0x00ffffff
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

// Deserialize converts from a disk sector into an MDTurbo struct.
// Returns a partition table data structure, or an error if the
// structure cannot be parsed. It does *NOT* check for overall structure
// validity; use MDTurbo.Validate() for that.
func Deserialize(data []byte) (MDTurbo, error) {
	var length, offset uint16
	var partmap MDTurbo

	// Basic sanity check
	if len(data) < PartitionBlkLen {
		return partmap, fmt.Errorf("not an MDTurbo partition map: too short (%v bytes, expected %v)", len(data), PartitionBlkLen)
	}

	// Reflect-based field deserialization
	mdt := reflect.TypeOf(partmap)
	// Iterate over all available fields and read the tag value
	for i := 0; i < mdt.NumField(); i++ {
		field := mdt.Field(i) // https://golang.org/pkg/reflect/#StructField
		switch field.Name {
		case "Partitions1", "Partitions2":
			{
				// Get the offset
				bigOffset := tagToOffset(field)
				if bigOffset == -1 {
					return partmap, fmt.Errorf("field %s is not tagged with a valid offset", field.Name)
				}
				offset = uint16(bigOffset)

				// Extract the next 64 bytes into a
				// slice of uint8 and zip those
				// partitions into the data structure
				// [MaxPartitions]Partition
				end := offset + 64 // FIXME: I don't like the magic number...
				partArray := zipAllPartitions(data[offset:end])

				// stick our new data structure back
				// into the parent struct
				reflect.ValueOf(&partmap).Elem().Field(i).Set(reflect.ValueOf(partArray))
			}
		default:
			bigOffset := tagToOffset(field)
			if bigOffset == -1 {
				return partmap, fmt.Errorf("field %s is not tagged with a valid offset", field.Name)
			}
			offset = uint16(bigOffset)

			// Set our field - algorithm depends on our size in
			// bytes
			// if we have an array (presumably of bytes), this is easy
			if field.Type.Kind() == reflect.Array || field.Type.Kind() == reflect.Slice {
				// get the array length
				length := reflect.ValueOf(&partmap).Elem().Field(i).Len()
				end := offset + uint16(length)
				// this is annoying
				reflect.Copy(reflect.ValueOf(&partmap).Elem().Field(i), reflect.ValueOf(data[offset:end]))
				continue // this field is done, do the next one
			}

			// if length == 0, then it wasn't
			// tagged, but we will calculate
			length = uint16(field.Type.Size())
			if length > PartitionBlkLen {
				return partmap, fmt.Errorf("field %s reports too-large size %d", field.Name, length)
			}
			end := offset + length
			if end > PartitionBlkLen {
				return partmap, fmt.Errorf("field %s has invalid end %d", field.Name, end)
			}
			switch length {
			case 1:
				value := data[offset]
				reflect.ValueOf(&partmap).Elem().Field(i).Set(reflect.ValueOf(value))
			case 2:
				value := binary.LittleEndian.Uint16(data[offset:end])
				reflect.ValueOf(&partmap).Elem().Field(i).Set(reflect.ValueOf(value))
			case 4:
				value := binary.LittleEndian.Uint32(data[offset:end])
				reflect.ValueOf(&partmap).Elem().Field(i).Set(reflect.ValueOf(value))
			default:
				return partmap, fmt.Errorf("field %s has unexpected length %d", field.Name, length)
			}
		}
	}
	return partmap, nil
}

// The partitions are actually represented as Start Sector
// numbers (0x20, 4 bytes for each of the 8 = 32 bytes),
// followed by lengths (another 32 bytes starting at 0x40).
// Here we take that byte array and return a struct
func zipAllPartitions(byteBlock []uint8) [MaxPartitions]Partition {
	var start, length uint32
	var results [MaxPartitions]Partition

	// We always unzip MaxPartitions to ensure correct
	// round-tripping of arbitrary sectors
	for item := 0; item < MaxPartitions; item++ {
		startOffset := 4 * item // uint32 = 4 bytes
		lengthOffset := (4 * MaxPartitions) + startOffset
		start = binary.LittleEndian.Uint32(byteBlock[startOffset : startOffset+4])
		length = binary.LittleEndian.Uint32(byteBlock[lengthOffset : lengthOffset+4])
		results[item] = Partition{Start: start, RawLength: length}
	}
	return results
}

// tagToOffset extracts an offset tag from a field and returns it as an
// int
func tagToOffset(field reflect.StructField) int64 {
	rawOffset := field.Tag.Get("offset")
	offset, err := strconv.ParseInt(rawOffset, 0, 16)
	if err != nil {
		return -1
	}
	return offset
}

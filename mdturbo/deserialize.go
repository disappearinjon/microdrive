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
)

// Deserialize converts from a disk sector into an MDTurbo struct.
// Returns a partition table data structure, or an error if the
// structure cannot be parsed. It does *NOT* check for overall structure
// validity; use MDTurbo.Validate() for that.
func Deserialize(data [512]byte) (MDTurbo, error) {
	var length, offset uint16
	var partmap MDTurbo

	// Reflect-based field deserialization
	mdt := reflect.TypeOf(partmap)
	// Iterate over all available fields and read the tag value
	for i := 0; i < mdt.NumField(); i++ {
		field := mdt.Field(i) // https://golang.org/pkg/reflect/#StructField
		switch field.Name {
		case "Partitions1", "Partitions2":
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
			end := offset + PartChunkSize
			partArray := zipAllPartitions(data[offset:end])

			// stick our new data structure back
			// into the parent struct
			reflect.ValueOf(&partmap).Elem().Field(i).Set(reflect.ValueOf(partArray))
		default:
			bigOffset := tagToOffset(field)
			if bigOffset == -1 {
				return partmap, fmt.Errorf("field %s is not tagged with a valid offset", field.Name)
			}
			offset = uint16(bigOffset)

			// Find our length in bytes
			// if we have an array (presumably of bytes), this is easy
			if field.Type.Kind() == reflect.Array || field.Type.Kind() == reflect.Slice {
				// get the array length
				length := reflect.ValueOf(&partmap).Elem().Field(i).Len()
				end := offset + uint16(length)
				// this is annoying
				reflect.Copy(reflect.ValueOf(&partmap).Elem().Field(i), reflect.ValueOf(data[offset:end]))
				continue // this field is done, do the next one
			}

			// calculate non-array / non-slice field length
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

// Package mdturbo provides the MicroDrive/Turbo partition map format,
// along with serializer and deserializer functions.
//
// The format is AFAIK undocumented, but the CiderPress source at
// https://github.com/fadden/ciderpress/blob/master/diskimg/MicroDrive.cpp
// contains a partial description.
package mdturbo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// Serialize converts from a disk sector into an MDTurbo struct.
// Returns a 512-byte array, or an error if one occurs while parsing. It
// does not insist upon validity of the structure it encodes, so that we
// can successfully round-trip arbitrary sectors.
func Serialize(partmap MDTurbo) ([PartitionBlkLen]uint8, error) {
	var data [PartitionBlkLen]uint8
	var end uint16

	mdt := reflect.TypeOf(partmap)
	for i := 0; i < mdt.NumField(); i++ {
		field := mdt.Field(i)
		bigOffset := tagToOffset(field)
		if bigOffset == -1 {
			continue
		}
		offset := uint16(bigOffset)
		switch field.Name {
		case "Partitions1", "Partitions2":
			pt := reflect.ValueOf(&partmap).Elem().Field(i).Interface()
			partTable, ok := pt.([MaxPartitions]Partition)
			if !ok {
				return data, fmt.Errorf("could not assert partition map type on %s", field.Name)
			}
			pmBytes, err := unzipAllPartitions(partTable)
			if err != nil {
				return data, fmt.Errorf("could not zip partition map %s: %v", field.Name, err)
			}
			copy(data[offset:offset+uint16(len(pmBytes))], pmBytes)
		default:
			// Find our our length in bytes
			// if we have an array (presumably of bytes), this is easy
			if field.Type.Kind() == reflect.Array || field.Type.Kind() == reflect.Slice {
				// Copy from the array into the bytes
				// There's got to be a better way,
				// doesn't there?
				length := reflect.ValueOf(&partmap).Elem().Field(i).Len()
				thing := reflect.ValueOf(&partmap).Elem().Field(i)
				for b := 0; b < length; b++ {
					data[offset+uint16(b)] = uint8(thing.Index(b).Uint())
				}
				continue // this field is done, do the next one
			}

			// Find non-array/slice field length
			length := uint16(field.Type.Size())
			if length > PartitionBlkLen {
				return data, fmt.Errorf("field %s reports too-large size %d", field.Name, length)
			}
			if length > 1 {
				end = offset + length
			} else {
				end = offset
			}
			if end > PartitionBlkLen {
				return data, fmt.Errorf("field %s has invalid end %d", field.Name, end)
			}

			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, reflect.ValueOf(&partmap).Elem().Field(i).Uint())
			if err != nil {
				return data, fmt.Errorf("failed to encode field %s: %v", field.Name, err)
			}
			for b := offset; b <= end; b++ {
				nextByte, err := buf.ReadByte()
				data[b] = nextByte
				if err != nil {
					return data, fmt.Errorf("failed while encoding field %s: %v", field.Name, err)
				}
			}
		}
	}
	return data, nil
}

// Return a block of bytes for the partition block, suitable for use
// being copied into the overall serialized sector at the appropriate
// offset.
func unzipAllPartitions(partitions [MaxPartitions]Partition) ([]uint8, error) {
	var results [PartChunkSize]uint8

	for item := 0; item < MaxPartitions; item++ {
		startOffset := 4 * item // uint32 = 4 bytes
		lengthOffset := (4 * MaxPartitions) + startOffset
		startBuf := new(bytes.Buffer)
		err := binary.Write(startBuf, binary.LittleEndian, partitions[item].Start)
		if err != nil {
			return results[:], fmt.Errorf("could not write partition %d start: %v", item, err)
		}
		for b := startOffset; b < startOffset+binary.Size(partitions[item].Start); b++ {
			nextByte, err := startBuf.ReadByte()
			results[b] = nextByte
			if err != nil {
				return results[:], fmt.Errorf("failed while writing partition %d start: %v", item, err)
			}
		}
		lengthBuf := new(bytes.Buffer)
		err = binary.Write(lengthBuf, binary.LittleEndian, partitions[item].RawLength)
		if err != nil {
			return results[:], fmt.Errorf("could not write partition %d length: %v", item, err)
		}
		for b := lengthOffset; b < lengthOffset+binary.Size(partitions[item].RawLength); b++ {
			nextByte, err := lengthBuf.ReadByte()
			results[b] = nextByte
			if err != nil {
				return results[:], fmt.Errorf("failed while writing partition %d length: %v", item, err)
			}
		}
	}
	return results[:], nil
}

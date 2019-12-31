// Package mdturbo provides the MicroDrive/Turbo partition map format,
// along with serializer and deserializer functions.
//
// The format is AFAIK undocumented, but the CiderPress source at
// https://github.com/fadden/ciderpress/blob/master/diskimg/MicroDrive.cpp
// contains a partial description.
package mdturbo

import (
	"fmt"
	"os"
	"strings"
)

// PrettyPrint formats a nice string representation of an MDTurbo struct.
// Returns a multi-line string.
//
// This is a very manual structure and would be improved by automation,
// but getting good output is harder here than with
// serialize/deserialize.
func PrettyPrint(partmap MDTurbo) string {
	var output strings.Builder
	var count uint8

	if partmap.Magic != 52426 {
		output.WriteString("WARNING: Not a Valid Partition Table: Magic Number Incorrect\n")
	}
	if partmap.Partitions1[0].Start != 256 {
		output.WriteString("WARNING: Not a Valid Partition Table: Wrong Start Sector\n")
	}

	output.WriteString("Cylinders\tHeads\tSectors\t\n")
	output.WriteString("---------\t-----\t-------\t\n")
	output.WriteString(fmt.Sprintf("%d\t%d\t%d\t\n\n", partmap.Cylinders,
		partmap.Heads, partmap.Sectors))
	output.WriteString("ROM Version\tBoot Partition\tPartition Count\t\n")
	output.WriteString("-----------\t--------------\t---------------\t\n")
	output.WriteString(fmt.Sprintf("%d\t%d\t%d\t\n", partmap.RomVersion,
		partmap.BootPart, partmap.PartCount1+partmap.PartCount2))

	output.WriteString("\nPartition\tStart\tEnd\tLength (KB)\t\n")
	for count = 0; count < partmap.PartCount1; count++ {
		output.WriteString(fmt.Sprintf("%d\t%s\t\n", count,
			partmap.Partitions1[count].String()))
	}

	for count = 0; count < partmap.PartCount2; count++ {
		output.WriteString(fmt.Sprintf("%d\t%s\t\n",
			count+partmap.PartCount1,
			partmap.Partitions2[count].String()))
	}

	output.WriteString("\nUnknown Regions (non-empty)\n")
	if !empty(partmap.Unknown1[:]) {
		output.WriteString(fmt.Sprintf("1\t%v\t\n", partmap.Unknown1))
	}
	if !empty(partmap.Unknown2[:]) {
		output.WriteString(fmt.Sprintf("2\t%v\t\n", partmap.Unknown2))
	}
	if !empty(partmap.Unknown3[:]) {
		output.WriteString(fmt.Sprintf("3\t%v\t\n", partmap.Unknown3))
	}
	if !empty(partmap.Unknown4[:]) {
		output.WriteString(fmt.Sprintf("4\t%v\t\n", partmap.Unknown4))
	}
	if !empty(partmap.Unknown5[:]) {
		output.WriteString(fmt.Sprintf("5\t%v\t\n", partmap.Unknown5))
	}
	if !empty(partmap.Unknown6[:]) {
		output.WriteString(fmt.Sprintf("6\t%v\t\n", partmap.Unknown6))
	}

	return output.String()
}

// empty returns true if a data structure is empty
func empty(thing []uint8) bool {
	for _, item := range thing {
		if item != 0 {
			return false
		}
	}
	return true
}

// GoPrint returns strings with a nice Go-style data structure
// representing the partition table
func GoPrint(partmap MDTurbo) string {
	var output strings.Builder
	data, err := partmap.Serialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not serialize parition map!\n")
		return ""
	}

	output.WriteString(fmt.Sprintf("var data = [%d]uint8{\n", len(data)))

	for count := 0; count < len(data); count++ {
		if count%8 == 0 { // First in line
			output.WriteString("\t")
		}
		output.WriteString(fmt.Sprintf("0x%02x,", data[count]))

		if count%8 == 7 || count == len(data)-1 { // Last in line
			output.WriteString("\n")
		} else {
			output.WriteString(" ")
		}

	}
	output.WriteString("}\n")
	return output.String()
}

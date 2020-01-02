package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

import "github.com/disappearinjon/microdrive/mdturbo"

// ImportCmd contains the CLI args and flags for Read command
type ImportCmd struct {
	Source    string `arg:"positional,required" help:"Hard Drive Image File"`
	Target    string `arg:"positional,required" help:"Microdrive/Turbo image file"`
	Partition uint8  `arg:required" help:"Partition number"`
	Type      string `arg:"-s"  help:"Source file type: auto, hdv" default:"auto"`
	Force     bool   `help:"Force write even in unsafe conditions" default:"false"`
}

func importPartition() error {
	var sourceLength int64 // Length of source file, minus headers

	// Open the source file passed in for reading
	source, err := os.Open(cli.Import.Source)
	defer source.Close()

	if err != nil {
		return err
	}

	// For the source image, if it's HDV, there's no seek required
	if cli.Import.Type == "auto" {
		cli.Import.Type = imageAutoDetect(cli.Import.Source)
	}
	switch strings.ToLower(cli.Import.Type) {
	case "hdv":
		fi, err := source.Stat()
		if err != nil {
			return fmt.Errorf("could not stat file %s", cli.Import.Source)
		}
		sourceLength = fi.Size()
	default:
		return fmt.Errorf("unknown image format %s", cli.Import.Type)
	}

	// Open the target file passed in for writing
	target, err := os.OpenFile(
		cli.Import.Target,
		os.O_RDWR, 0666)
	defer target.Close()

	// Get first disk sector, where the partition table sits
	var firstSector [mdturbo.SectorSize]byte
	buf := make([]byte, mdturbo.SectorSize)
	_, err = target.Read(buf)
	if err != nil {
		return err
	}
	copy(firstSector[:], buf)

	// Parse the partition table
	partMap, err := mdturbo.Deserialize(firstSector)
	if err != nil {
		return err
	}

	// Validate partition table format
	if !partMap.Validate() {
		if cli.Import.Force {
			fmt.Fprintf(os.Stderr, "WARNING: partition map appears invalid\n")
		} else {
			return fmt.Errorf("partition map on %s appears invalid",
				cli.Import.Target)
		}
	}

	// Get partition data
	partition, err := partMap.GetPartition(cli.Import.Partition)
	if err != nil {
		return fmt.Errorf("failed to get partition: %v", err)
	}

	// Fail if partition is smaller than the file to be read
	if sourceLength > int64(partition.Length())*mdturbo.SectorSize {
		return fmt.Errorf("source (%d) larger than target partition (%d)",
			sourceLength, partition.Length()*mdturbo.SectorSize)
	}

	// Seek to the beginning of the partition
	_, err = target.Seek(int64(partition.Start)*mdturbo.SectorSize, os.SEEK_SET)

	// Copy bytes
	bytesWritten, err := io.CopyN(target, source, sourceLength)
	if err != nil {
		return fmt.Errorf("import copy returned error: %v", err)
	}
	if bytesWritten != sourceLength {
		return fmt.Errorf("import expected %d bytes; copied %d",
			sourceLength, bytesWritten)
	}

	// And done
	return nil
}

// Pick output format based on filename and return as a string
func imageAutoDetect(filename string) string {
	var filetype string

	// Get file suffix
	fileparts := strings.Split(filename, ".")
	if len(fileparts) == 1 {
		return ""
	}
	suffix := fileparts[len(fileparts)-1]
	switch strings.ToLower(suffix) {
	default:
		filetype = suffix // this will work or fail independently
	}
	return strings.ToLower(filetype)
}

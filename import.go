package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

import (
	"github.com/disappearinjon/microdrive/h2mg"
	"github.com/disappearinjon/microdrive/mdturbo"
)

// ImportCmd contains the CLI args and flags for the import command
type ImportCmd struct {
	Source    string `arg:"positional,required" help:"Hard Drive Image File"`
	Target    string `arg:"positional,required" help:"Microdrive/Turbo image file"`
	Type      string `arg:"-s"  help:"Source file type: auto, 2mg, hdv, po" default:"auto"`
	Partition uint8  `arg:"required" help:"Partition number"`
	Force     bool   `help:"Force write even in unsafe conditions" default:"false"`
}

func importPartition() error {
	return importImage(cli.Import.Source, cli.Import.Target,
		cli.Import.Type, cli.Import.Partition, cli.Import.Force)
}

func importImage(sourceFile, targetFile, targetType string, partNum uint8, force bool) error {
	var sourceLength int64 // Length of source file, minus headers

	// Open the source file passed in for reading
	source, err := os.Open(sourceFile)
	defer source.Close()

	if err != nil {
		return err
	}

	sourceLength, err = getSourceLength(source, targetType)
	if err != nil {
		return fmt.Errorf("could not get source length for %s: %v", sourceFile, err)
	}

	target, partMap, err := getTarget(targetFile, force)
	defer target.Close()
	if err != nil {
		return fmt.Errorf("could not open target %s: %v", targetFile, err)
	}

	// Get partition data
	partition, err := partMap.GetPartition(partNum)
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

// getSourceLength gets the length of the source image, and sets the
// offset at the end of the header, if any.
func getSourceLength(source *os.File, targetType string) (length int64, err error) {
	fi, err := source.Stat()
	if err != nil {
		return -1, fmt.Errorf("could not stat source file")
	}
	fileName := fi.Name()

	if targetType == "auto" {
		targetType = imageAutoDetect(fileName)
	}
	switch strings.ToLower(targetType) {
	case "2mg":
		buf := make([]uint8, h2mg.HeaderSize)
		read, err := source.Read(buf)
		if err != nil {
			return -1, fmt.Errorf("could not read %s: %v", fileName, err)
		}
		if read != h2mg.HeaderSize {
			return -1, fmt.Errorf("%s: unexpected read length (expected %d, got %d)",
				fileName, h2mg.HeaderSize, read)
		}
		sourceHeader, err := h2mg.Parse2MG(buf)
		if err != nil {
			return -1, fmt.Errorf("could not parse %s header: %v", fileName, err)
		}
		err = sourceHeader.Validate()
		if err != nil {
			return -1, fmt.Errorf("%s: could not validate: %v", fileName, err)
		}
		// Move to beginning of data - we should already be
		// there but this doesn't hurt
		_, err = source.Seek(int64(sourceHeader.Offset), os.SEEK_SET)
		if err != nil {
			return length, fmt.Errorf("could not seek to data for %s: %v", fileName, err)
		}
		return int64(sourceHeader.Length), nil
	case "hdv", "po":
		return fi.Size(), nil
		// For the source image, if it's HDV, there's no seek
		// required because we haven't read anything
	default:
		return -1, fmt.Errorf("unknown image format %s", targetType)
	}
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

func getTarget(targetFile string, force bool) (target *os.File, partMap mdturbo.MDTurbo, err error) {
	// Open the target file passed in for writing
	target, err = os.OpenFile(
		targetFile,
		os.O_RDWR, 0666)
	if err != nil {
		return target, partMap, fmt.Errorf("could not open %s: %v", targetFile, err)
	}

	// Get first disk sector, where the partition table sits
	var firstSector [mdturbo.SectorSize]byte
	buf := make([]byte, mdturbo.SectorSize)
	_, err = target.Read(buf)
	if err != nil {
		return target, partMap, err
	}
	copy(firstSector[:], buf)

	// Parse the partition table
	partMap, err = mdturbo.Deserialize(firstSector)
	if err != nil {
		return target, partMap, err
	}

	// Validate partition table format
	if !partMap.Validate() {
		if force {
			fmt.Fprintf(os.Stderr, "WARNING: partition map appears invalid\n")
		} else {
			return target, partMap, fmt.Errorf("partition map on %s appears invalid",
				targetFile)
		}
	}

	return
}

package main

import (
	"fmt"
	"io"
	"os"
)

import (
	// "github.com/disappearinjon/microdrive/h2mg"
	"github.com/disappearinjon/microdrive/mdturbo"
)

// ExportCmd contains the CLI args and flags for the export command
type ExportCmd struct {
	Source    string `arg:"positional,required" help:"Microdrive/Turbo image file"`
	Target    string `arg:"positional,required" help:"Hard Drive Image File"`
	Type      string `arg:"-s"  help:"Target file type: auto, 2mg, hdv, po" default:"auto"`
	Partition uint8  `arg:"required" help:"Partition number"`
	Force     bool   `help:"Force overwrite of an existing disk" default:"false"`
}

func exportPartition() error {
	return exportImage(cli.Export.Source, cli.Export.Target,
		cli.Export.Type, cli.Export.Partition, cli.Export.Force)
}

func exportImage(sourceFile, targetFile, targetType string, partNum uint8, force bool) error {
	source, partMap, err := getSource(sourceFile, force)
	defer source.Close()
	if err != nil {
		return err
	}

	// Fail early if partition is unavailable
	if partNum >= partMap.PartCount() {
		return fmt.Errorf("could not extract partition %d - maximum partition # is %d",
			partNum, partMap.PartCount()-1)
	}

	// Fail early if can't write in the requested format
	if targetType == "auto" {
		targetType = imageAutoDetect(targetFile)
	}
	if targetType != "hdv" && targetType != "po" {
		return fmt.Errorf("no support for disk image type %s", targetType)
	}

	// Check if target file already exists - if so, and not force,
	// then fail
	_, err = os.Stat(targetFile)
	if (os.IsExist(err) || err == nil) && !force {
		return fmt.Errorf("target %s exists - will not overwrite", targetFile)
	}

	// OK, create & truncate it
	target, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("could not create target %s: %v", targetFile, err)
	}
	defer target.Close()

	// Get partition data
	partition, err := partMap.GetPartition(partNum)
	if err != nil {
		return fmt.Errorf("failed to get partition: %v", err)
	}

	// Seek to the beginning of the partition
	_, err = source.Seek(int64(partition.Start)*mdturbo.SectorSize, os.SEEK_SET)

	// Copy bytes
	length := int64(partition.Length()) * mdturbo.SectorSize
	bytesWritten, err := io.CopyN(target, source, length)
	if err != nil {
		return fmt.Errorf("export copy returned error: %v", err)
	}
	if bytesWritten != length {
		return fmt.Errorf("export expected %d bytes; copied %d",
			length, bytesWritten)
	}

	// And done
	return nil
}

func getSource(sourceFile string, force bool) (source *os.File, partMap mdturbo.MDTurbo, err error) {
	// Open the source file passed in for reading
	source, err = os.Open(sourceFile)
	if err != nil {
		err = fmt.Errorf("could not open %s: %v", sourceFile, err)
		return
	}

	// Get first disk sector, where the partition table sits
	var firstSector [mdturbo.SectorSize]byte
	buf := make([]byte, mdturbo.SectorSize)
	_, err = source.Read(buf)
	if err != nil {
		return
	}
	copy(firstSector[:], buf)

	// Parse the partition table
	partMap, err = mdturbo.Deserialize(firstSector)
	if err != nil {
		return
	}

	// Validate partition table format
	if !partMap.Validate() {
		if force {
			fmt.Fprintf(os.Stderr, "WARNING: partition map appears invalid\n")
		} else {
			err = fmt.Errorf("partition map on %s appears invalid", sourceFile)
			return
		}
	}

	return
}

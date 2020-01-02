package main

import (
	"fmt"
	"os"
)

import (
	"github.com/disappearinjon/microdrive/mdturbo"
)

// AppendCmd contains the CLI args and flags for the append command
type AppendCmd struct {
	Source string `arg:"positional,required" help:"Hard Drive Image File"`
	Target string `arg:"positional,required" help:"Microdrive/Turbo image file"`
	Type   string `arg:"-s"  help:"Source file type: auto, 2mg, hdv, po" default:"auto"`
	Force  bool   `help:"Force write even in unsafe conditions" default:"false"`
}

func appendPartition() error {
	// Get the size of our source volume, in blocks
	source, err := os.Open(cli.Append.Source)
	if err != nil {
		source.Close()
		return err
	}
	sourceLength, err := getSourceLength(source, cli.Append.Type)
	source.Close()
	if err != nil {
		return fmt.Errorf("could not get %s image length: %v", cli.Append.Target, err)
	}
	if sourceLength == -1 {
		return fmt.Errorf("recieved incorrect image size (-1) for %s", cli.Append.Target)
	}

	// Length is in bytes; convert to blocks
	blockCount := sourceLength / mdturbo.SectorSize
	if sourceLength%mdturbo.SectorSize != 0 {
		blockCount++
	}
	if blockCount == 0 {
		return fmt.Errorf("nonsense size for file %s", cli.Append.Source)
		return nil
	}

	// Open the target file
	target, partMap, err := getTarget(cli.Append.Target, cli.Append.Force)
	if err != nil {
		return fmt.Errorf("could not open %s: %v", cli.Append.Target, err)
	}

	// Add a new partition map
	partNum, err := partMap.AddPartition(uint32(blockCount))
	if err != nil {
		return fmt.Errorf("could not add partition to table on %s: %v", cli.Append.Target, err)
	}
	if partNum < 0 {
		return fmt.Errorf("could not get new partition number for %s", cli.Append.Target)
	}

	serialized, err := partMap.Serialize()
	if err != nil {
		return fmt.Errorf("could not serialize updated partition table for %s: %v", cli.Append.Target, err)
	}

	// Seek to 0 and write it out
	bytesWritten, err := target.WriteAt(serialized[:], 0)
	if err != nil {
		return fmt.Errorf("could not write updated partition table for %s: %v", cli.Append.Target, err)
	}
	if bytesWritten != len(serialized) {
		return fmt.Errorf("partition table write unexpected length (got %d, expected %d)",
			bytesWritten, len(serialized))
	}

	// Flush and close
	err = target.Sync()
	if err != nil {
		return err
	}
	err = target.Close()
	if err != nil {
		return err
	}

	return importImage(cli.Append.Source, cli.Append.Target, cli.Append.Type, uint8(partNum), cli.Append.Force)
}

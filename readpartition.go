package main

import "encoding/json"
import "fmt"
import "os"

import "github.com/disappearinjon/microdrive/mdturbo"

// ReadCmd contains the CLI args and flags for Read command
type ReadCmd struct {
	Image  string `arg:"positional,required" help:"Microdrive/Turbo image file"`
	File   string `arg:"-f" help:"Output filename. - for STDOUT" default:"-"`
	Output string `arg:"-o" help:"Output format: text, go, go-bin, json" default:"text"`
}

func readPartition() error {
	var output *os.File
	var err error

	// Open the file passed in for reading
	imagefile, err := os.Open(cli.Read.Image)
	defer imagefile.Close()
	if err != nil {
		return err
	}

	// Set output device
	switch cli.Read.File {
	case "-":
		output = os.Stdout
	default:
		output, err = os.Create(cli.Read.File)
		if err != nil {
			return err
		}
	}
	defer output.Close()

	// Get first disk sector, where the partition table sits
	var firstSector [mdturbo.SectorSize]byte
	buf := make([]byte, mdturbo.SectorSize)
	read, err := imagefile.Read(buf)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Read %v bytes\n", read)
	copy(firstSector[:], buf)

	// Parse the sector
	partMap, err := mdturbo.Deserialize(firstSector)
	if err != nil {
		return err
	}
	if !partMap.Validate() {
		fmt.Fprintf(os.Stderr, "WARNING: partition map appears invalid\n")
	}

	// Print it
	switch cli.Read.Output {
	case "go":
		fmt.Fprintf(output, "%#v\n", partMap)
	case "go-bin":
		fmt.Fprintf(output, mdturbo.GoPrint(partMap))
	case "text":
		fmt.Fprintf(output, partMap.String())
	case "json":
		marshaled, err := json.MarshalIndent(partMap, "", "\t")
		if err != nil {
			return err
		}
		fmt.Fprintf(output, "%v\n", string(marshaled))
	default:
		return fmt.Errorf("unknown output format %s", cli.Read.Output)
	}

	// And done
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

import (
	"github.com/disappearinjon/microdrive/mdturbo"
)

// CLI for Write command
type WriteCmd struct {
	Image  string `arg:"positional, required" help:"Microdrive/Turbo image file"`
	File   string `arg:"-f" help:"Input filename (JSON presumed). - for STDIN" default:"-"`
	Force  bool   `help:"Write partition table, even when considered invalid" default:"false"`
	Offset int64  `arg:"-o" help:"File byte offset at which to write the table. Override with caution!" default:"0"`
}

func writePartition() error {
	var input *os.File
	var err error

	// Open the file passed in for writing - create a new file if
	// nothing present
	imagefile, err := os.OpenFile(cli.Write.Image, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer imagefile.Close()

	// Get input file
	switch cli.Write.File {
	case "-":
		input = os.Stdin
	default:
		input, err = os.Open(cli.Write.File)
		if err != nil {
			return err
		}
	}
	defer input.Close()

	// Read in our input file and convert to an object
	var mdt mdturbo.MDTurbo
	buf, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Read %v bytes from %s\n", len(buf), input.Name())
	err = json.Unmarshal(buf, &mdt)
	if err != nil {
		return err
	}
	if !mdt.Validate() {
		if !cli.Write.Force {
			return fmt.Errorf("file %s did not represent a valid partition table", cli.Write.File)
		}
		fmt.Fprintf(os.Stderr, "%s did not validate, but --force is specified. Writing anyway.\n", cli.Write.File)
	}

	// Serialize the partition table and write it to disk
	pt, err := mdt.Serialize()
	if err != nil {
		return fmt.Errorf("could not serialize partition table: %v", err)
	}

	write, err := imagefile.WriteAt(pt[:], cli.Write.Offset)
	if err != nil {
		return fmt.Errorf("failed to write image! wrote %d bytes with error %v", write, err)
	}

	return nil
}

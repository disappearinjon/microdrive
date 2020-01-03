package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

import "github.com/disappearinjon/microdrive/mdturbo"

// ReadCmd contains the CLI args and flags for Read command
type ReadCmd struct {
	Image  string `arg:"positional,required" help:"Microdrive/Turbo image file"`
	File   string `arg:"-f" help:"Output filename. - for STDOUT" default:"-"`
	Output string `arg:"-o" help:"Output format: auto, text, go, go-bin, json" default:"auto"`
}

func readPartition() (err error) {
	var output *os.File

	// Set output device
	switch cli.Read.File {
	case "-":
		output = os.Stdout
		if strings.ToLower(cli.Read.Output) == "auto" {
			cli.Read.Output = "text"
		}
	default:
		output, err = os.Create(cli.Read.File)
		if err != nil {
			return
		}
	}
	defer output.Close()

	partMap, err := GetPartitionTable(cli.Read.Image)
	if err != nil {
		return
	}

	// Print it
	if cli.Read.Output == "auto" {
		cli.Read.Output = autoDetect(cli.Read.File)
	}
	switch strings.ToLower(cli.Read.Output) {
	case "go":
		fmt.Fprintf(output, "%#v\n", partMap)
	case "go-bin":
		fmt.Fprintf(output, mdturbo.GoPrint(partMap))
	case "text":
		fmt.Fprintf(output, partMap.String())
	case "json":
		var marshaled []byte
		marshaled, err = json.MarshalIndent(partMap, "", "\t")
		if err != nil {
			return
		}
		fmt.Fprintf(output, "%v\n", string(marshaled))
	default:
		return fmt.Errorf("unknown output format %s", cli.Read.Output)
	}

	// And done
	return
}

// Pick output format based on filename and return as a string
func autoDetect(filename string) string {
	var filetype string

	// Get file suffix
	fileparts := strings.Split(filename, ".")
	if len(fileparts) == 1 {
		return ""
	}
	suffix := fileparts[len(fileparts)-1]
	switch strings.ToLower(suffix) {
	case "txt":
		filetype = "text"
	case "jsn", "json":
		filetype = "json"
	default:
		filetype = suffix // this will work or fail independently
	}
	return filetype
}

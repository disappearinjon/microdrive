package main

import "fmt"
import "os"

import "github.com/alexflint/go-arg"

// CLI flags and values

type args struct {
	Import *ImportCmd `arg:"subcommand:import"`
	Read   *ReadCmd   `arg:"subcommand:read"`
	Write  *WriteCmd  `arg:"subcommand:write"`
}

func (args) Description() string {
	return "CLI utility for manipulating Microdrive/Turbo images"
}

var cli args

func main() {
	var err error
	parsed := arg.MustParse(&cli)
	subcommand := parsed.SubcommandNames()
	if len(subcommand) != 1 {
		parsed.Fail("Must specify a command")
	}
	switch subcommand[0] {
	case "import":
		err = importPartition()
	case "read":
		err = readPartition()
	case "write":
		err = writePartition()
	default:
		parsed.Fail("Unknown command")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

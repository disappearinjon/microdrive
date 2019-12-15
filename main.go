package main

import "fmt"
import "os"

import "github.com/alexflint/go-arg"

// CLI flags and values
type ReadCmd struct {
	Image string `arg:"positional,required" help:"Microdrive/Turbo image file"`
}

type args struct {
	Read *ReadCmd `arg:"subcommand:read"`
}

func (args) Description() string {
	return "CLI utility for manipulating Microdrive/Turbo images"
}

var cli args

// checkFatal is a quick hack to simplify error handling for fatal
// errors.
func checkFatal(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var err error
	parsed := arg.MustParse(&cli)
	subcommand := parsed.SubcommandNames()
	if len(subcommand) != 1 {
		parsed.Fail("Must specify a command")
	}
	switch subcommand[0] {
	case "read":
		err = readPartition()
	default:
		parsed.Fail("Unknown command")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

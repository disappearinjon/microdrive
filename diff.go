// diff command for CLI
package main

import "fmt"

import "gopkg.in/d4l3k/messagediff.v1"

// DiffCmd contains the CLI args and flags for Diff command
type DiffCmd struct {
	File1 string `arg:"positional,required" help:"First Microdrive/Turbo image file"`
	File2 string `arg:"positional, required"  help:"Second Microdrive/Turbo image file"`
}

func diffPartitions() (err error) {
	pt1, err := GetPartitionTable(cli.Diff.File1)
	if err != nil {
		return
	}
	pt2, err := GetPartitionTable(cli.Diff.File2)
	if err != nil {
		return
	}

	diff, equal := messagediff.PrettyDiff(pt1, pt2)
	if !equal {
		fmt.Println(diff)
	}

	return
}

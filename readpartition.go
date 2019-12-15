package main

import "fmt"
import "os"

import "github.com/disappearinjon/microdrive/mdturbo"

func readPartition() error {
	file, err := os.Open(cli.Read.Image)
	if err != nil {
		return err
	}
	firstSector := make([]byte, mdturbo.SectorSize)

	read, err := file.Read(firstSector)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Read %v bytes\n", read)

	partMap, err := mdturbo.Deserialize(firstSector)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", partMap)
	return nil
}

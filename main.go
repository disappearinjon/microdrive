package main

import "fmt"
import "os"

import "github.com/disappearinjon/microdrive/mdturbo"

// checkFatal is a quick hack to simplify error handling for fatal
// errors.
func checkFatal(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	file, err := os.Open("test.bin")
	checkFatal(err)

	firstSector := make([]byte, mdturbo.SectorSize)
	read, err := file.Read(firstSector)
	checkFatal(err)
	fmt.Printf("Read %v bytes\n", read)
	partMap, err := mdturbo.Deserialize(firstSector)
	checkFatal(err)
	fmt.Printf("%+v\n", partMap)
}

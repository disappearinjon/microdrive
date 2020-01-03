# To Do (semi-prioritized)
* CLI: edit partition table command - interactive mode?
* CLI: partition table diff
* CLI: Extract partition to separate file (2MG)
* CLI: support for .gz and .bz2 files
* CLI: unit tests
* CLI: provide abstraction layer for file actions?
* Cleanup: omit JSON byte fields containing only zeroes
* Documentation
* MDTurbo Library: abstract away split in partition sets from data
  structure?
* MDTurbo Library: add more unit tests (down from 85% to 50%)

# Done
* CLI: Fix bug where import defaulted to partition 0 overwrite
* CLI: Add support for .po disk files (same as HDV)
* CLI: Import HDV/2MG files to new partitions ("append")
* CLI: Import 2MG files to existing partitions
* Reversed order of Done list to put most recent on top
* CLI: Import separate file to partition - same size only
* CLI: autodetect file output types based on filename
* Library: pretty-print for text & go-format output
* Find boot partition field; add to partition map. (I'm guessing as to
  whether it's one byte or two; two feels safer for now.)
* BUG: Why are my omitempty-tagged JSON fields being included? (Answer:
  not actually a bug, see https://github.com/golang/go/issues/29310)
* CLI: write partition table command (from json file)
* Serialization: unit tests
* Serialization of partition table
* Rework deserialize to use struct tags & reflections
* Unit Tests for partition table
* Validate function for partition table structure
* End() function on MDTurbo struct to simplify other code
* Output formats for structure printing
* Basic CLI: filename, print command
* Add license
* To-Do list
* Deserialization of partition table
* CLI: Extract partition to separate file (HDV/PO)
* Cleanup: replace magic numbers with constants or calculated values

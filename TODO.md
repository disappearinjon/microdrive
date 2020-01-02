# To Do (semi-prioritized)
* CLI: edit partition table command - interactive mode?
* CLI: partition table diff
* CLI: Import HDV/2MG files to new partitions
* CLI: Import 2MG files to existing partitions
* CLI: Extract partition to separate file (HDV)
* CLI: Extract partition to separate file (2MG)
* CLI: unit tests
* Cleanup: remove magic numbers, replace with constants or calculated
  values
* Documentation

# Done
* Deserialization of partition table
* To-Do list
* Add license
* Basic CLI: filename, print command
* Output formats for structure printing
* End() function on MDTurbo struct to simplify other code
* Validate function for partition table structure
* Unit Tests for partition table
* Rework deserialize to use struct tags & reflections
* Serialization of partition table
* Serialization: unit tests
* CLI: write partition table command (from json file)
* BUG: Why are my omitempty-tagged JSON fields being included? (Answer:
  not actually a bug, see https://github.com/golang/go/issues/29310)
* Find boot partition field; add to partition map. (I'm guessing as to
  whether it's one byte or two; two feels safer for now.)
* Library: pretty-print for text & go-format output
* CLI: autodetect file output types based on filename
* CLI: Import separate file to partition - same size only

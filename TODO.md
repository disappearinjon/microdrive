# To Do (semi-prioritized)
* Library: pretty-print for text output
* Nit: file output types based on filename?
* Documentation
* CLI: unit tests
* CLI: edit partition table command - interactive mode?
* CLI: partition table diff
* Extract partition to separate file (2MG?)
* Import separate file to partition - same size only

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

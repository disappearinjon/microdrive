# Microdrive
Golang tools for manipulating a MicroDrive/Turbo disk or image
by Jon Lasser <jon@lasser.org>

There's not much here yet. See my [To-Do List](./TODO.md) for more
details on what I have planned.

This project is under an MIT license. See [the license](./LICENSE.txt)
for specifics.

# Building and Installing

* [Download](https://golang.org/dl/) and install Go, if you have not
  already done so.
* Install dependencies: `go get -d ./...`
* Run tests: `go test -v ./...`
* Build and install: `go install`

This will install binaries in the directory defined by
[the rules in the Go documentation](https://golang.org/cmd/go/#hdr-Compile_and_install_packages_and_dependencies),
which is likely to be either `$GOPATH/bin` or `$HOME/go/bin`.

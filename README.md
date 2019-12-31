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

# Using the `microdrive` tool

Right now, the `microdrive` tool can read and write partition tables.
While an interactive editor is in the works, the best way to edit a
partition table today is to:

**When inserting the CF card in a Mac, you will be prompted to
Initialize, Eject, or Ignore the card. *Always* pick Ignore to avoid
erasing the card.**

1. Make a partition image using `dd` or similar tool. On my Mac's CF
   reader, I see the CF drive as `/dev/disk2` or `/dev/disk7` depending.
   (Use *Disk Utility* to determine the current disk name. Presuming
   it's /dev/disk2, the command line is something like
   `dd if=/dev/disk2 of=mydrive.mdt`.  You may need to use `sudo` to
   execute this command as root, adjusting the ownership of the file as
   necessary.
1. **Back up your partition image.** This is *very important*, as you
   may want or need to restore your original disk image at some point.
1. `microdrive read --output json --file partitions.json mydrive.mdt` to
   read the partition table from *mydrive.mdt* into a file named
   *partitions.json*.
1. Edit the JSON to reflect your desired partition table
1. `microdrive write --file partitions.json mydrive.mdt` to update the
   *mydrive.mdt* image.
1. Copy your updated partition to the compact flash:
  `dd if=mydrive.mdt of=/dev/disk2` or equivalent. Again, you may need
  to use `sudo` to work around permissions issues.

# Getting Help

The microdrive project is a labor of love--but I'd love to help you too!
If you need something,
[file an issue in Github](https://github.com/disappearinjon/microdrive/issues)
or send me e-mail at my address above. I'll do my best to help you!

# Helping me

If you've ever edited your Microdrive/Turbo partition table using the
on-Apple tool, I'd love to get a copy of your boot sector, and a text
description or screenshot of your configuration. This will help me
ensure that I'm properly addressing real-world configurations. You can
send me just the partition table (512 bytes), and I'll add it to my
catalog.

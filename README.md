# rpcs3-gameupdater

Fetches game updates for the RPCS3 emulator.

Still very much work in progress.

## Build status

All platforms [![Build Status](https://travis-ci.com/zom-ponks/rpcs3-gameupdater.svg?branch=master)](https://travis-ci.com/zom-ponks/rpcs3-gameupdater)


## TODO

* DL selection
* DL path
* Package versioning
* configuration / config file persistence
* CLI/command line arguments
* UI

## Building Instructions

### Windows

* Install [MSYS2](https://www.msys2.org/)
* Install [go](https://golang.org/dl/)
* * Restart your terminal and editors
* `go build` in the checkout directory

### Linux

* Install gcc/g++ and libraries for building (e.g. the *build-essentials* package on Debian/Ubuntu)
* Install go (package `golang` on Debian/Ubuntu)
* `go build` in the checkout directory

### BSD
* Install go (`pkg install go`)
* `go build` in the checkout directory


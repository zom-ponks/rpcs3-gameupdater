// +build windows
package main

import (
	"os"

	"golang.org/x/sys/windows"
)

func isTTY() bool {
	fd := os.Stdout.Fd()

	var st uint32
	err := windows.GetConsoleMode(windows.Handle(fd), &st)
	isTerm := (err == nil)

	return isTerm
}

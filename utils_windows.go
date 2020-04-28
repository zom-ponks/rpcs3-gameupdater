// +build windows
package main

import (
	"os"

	"golang.org/x/sys/windows"
)

func isTTY() bool {
	term := true
	fd := os.Stdout.Fd()

	var st uint32
	err := windows.GetConsoleMode(windows.Handle(fd), &st)
	term = err == nil

	return term
}
// +build linux freebsd netbsd openbsd

package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func isTTY() bool {
	fd := os.Stdout.Fd()
	fmt.Println("linux/bsd")
	isTerm := terminal.IsTerminal(int(fd))

	return isTerm
}

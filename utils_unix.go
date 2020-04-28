// +build linux freebsd netbsd openbsd

package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func isTTY() bool {
	term := true
	fd := os.Stdout.Fd()
	fmt.Println("linux/bsd")
	term = terminal.IsTerminal(int(fd))

	return term
}

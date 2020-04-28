/* rpcs3-gameupdater - logging functions */

package main

import (
	"fmt"

	"gopkg.in/gookit/color.v1"
)

/* standard print */

func print(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

/* prints debug messages */

func printDebug(format string, a ...interface{}) {
	if fetchConfig().verbosity < Debug {
		return
	}
	if fetchConfig().color {
		color.Gray.Printf(format+"\n", a...)
	} else {
		fmt.Printf(format+"\n", a...)
	}
}

/* prints error messages */

func printError(format string, a ...interface{}) {
	if fetchConfig().verbosity < Error {
		return
	}
	if fetchConfig().color {
		color.Red.Printf(format+"\n", a...)
	} else {
		fmt.Printf(format+"\n", a...)
	}
}

/* prints info messages */

func printInfo(format string, a ...interface{}) {
	if fetchConfig().verbosity < Info {
		return
	}
	if fetchConfig().color {
		color.Green.Printf(format+"\n", a...)
	} else {
		fmt.Printf(format+"\n", a...)
	}
}

/* prints warning messages */

func printWarning(format string, a ...interface{}) {
	if fetchConfig().verbosity < Warning {
		return
	}
	if fetchConfig().color {
		color.Yellow.Printf(format+"\n", a...)
	} else {
		fmt.Printf(format+"\n", a...)
	}
}

/* prints over the same line */

func sameLinePrint(format string, a ...interface{}) {
	fmt.Print("\033[G\033[K") // move the cursor left and clear the line
	fmt.Printf(format+"\n", a...)
	fmt.Print("\033[A") // move the cursor up
}

func stopSameLinePrint() {
	fmt.Print("\033[B") // move the cursor down
}

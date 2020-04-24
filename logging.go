/* rpcs3-gameupdater - logging functions */

package main

import (
	"fmt"

	"gopkg.in/gookit/color.v1"
)

/* prints debug messages */

func printDebug(format string, a ...interface{}) {
	if !fetchConfig().debug {
		return
	}
	// check variable for verbosity
	if true {
		color.Gray.Printf(format, a...)
	}
}

/* prints info messages */

func printInfo(format string, a ...interface{}) {
	if !fetchConfig().debug {
		return
	}
	if fetchConfig().color {
		color.Green.Printf("\n"+format+"\n", a...)
	} else {
		fmt.Printf("\n%s\n", format)
	}
}

/* prints error messages */

func printError(format string, a ...interface{}) {
	if !fetchConfig().debug {
		return
	}
	if fetchConfig().color {
		color.Red.Printf(format, a...)
	} else {
		fmt.Printf("\n%s\n", format)
	}
}

/* prints warning messages */

func printWarning(format string, a ...interface{}) {
	if !fetchConfig().debug {
		return
	}
	if fetchConfig().color {
		color.Yellow.Printf(format, a...)
	} else {
		fmt.Printf("\n%s\n", format)
	}
}

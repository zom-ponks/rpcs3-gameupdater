/* rpcs3-gameupdater - logging functions */

package main

import (
	"gopkg.in/gookit/color.v1"
)

/* prints debug messages */

func printDebug(format string, a ...interface{}) {
	// check variable for verbosity
	if true {
		color.Gray.Printf(format, a...)
	}
}

/* prints info messages */

func printInfo(format string, a ...interface{}) {
	color.Green.Printf("\n"+format+"\n", a...)
}

/* prints error messages */

func printError(format string, a ...interface{}) {
	color.Red.Printf(format, a...)
}

/* prints warning messages */

func printWarning(format string, a ...interface{}) {
	color.Yellow.Printf(format, a...)
}

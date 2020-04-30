/* rpcs3-gameupdater - logging functions */

package main

import (
	"fmt"
	"os"

	"gopkg.in/gookit/color.v1"
)

/* standard print */

func print(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format+"\n", a...)
}

/* prints debug messages */

func printDebug(format string, a ...interface{}) {
	if fetchConfig().verbosity < Debug {
		return
	}
	if fetchConfig().color {
		format = color.FgGray.Render(format)
	}
	fmt.Fprintf(os.Stdout, format+"\n", a...)

}

/* prints error messages */

func printError(format string, a ...interface{}) {
	if fetchConfig().verbosity < Error {
		return
	}
	if fetchConfig().color {
		format = color.FgRed.Render(format)
	}
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

/* prints info messages */

func printInfo(format string, a ...interface{}) {
	if fetchConfig().verbosity < Info {
		return
	}
	if fetchConfig().color {
		format = color.FgGreen.Render(format)
	}
	fmt.Fprintf(os.Stdout, format+"\n", a...)
}

/* prints warning messages */

func printWarning(format string, a ...interface{}) {
	if fetchConfig().verbosity < Warning {
		return
	}
	if fetchConfig().color {
		format = color.FgYellow.Render(format)
	}
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

/* prints over the same line */

func sameLinePrint(format string, a ...interface{}) {
	fmt.Fprint(os.Stdout, "\033[G\033[K") // move the cursor left and clear the line
	fmt.Fprintf(os.Stdout, format+"\n", a...)
	fmt.Fprint(os.Stdout, "\033[A") // move the cursor up
}

func stopSameLinePrint() {
	fmt.Fprint(os.Stdout, "\033[B") // move the cursor down
}

/* prints the fields of a struct */

func printStruct(a interface{}) {
	printDebug("%+v", a)
}

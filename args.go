/* rpcs3-gameupdater - command line argument parsing */

package main

import (
	"fmt"
)

const helpText string = `
rpcs3-downloader

-v, --version Display Version
-c <configuration file>, --conf <configuration file> Override default configuration file
`

// parse args and update config accordingly
func parseArguments() {
	fmt.Printf(helpText)

	//version string := appVersion
	//flag.IntVar(&version, "version", , "help message for flagname")

}

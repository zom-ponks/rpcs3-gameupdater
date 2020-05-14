/* rpcs3-gameupdater - command line argument parsing */

package main

import (
	"flag"
	"fmt"
	"os"
)

const helpText string = `
rpcs3-downloader

-v, --version Display Version
-c <configuration file>, --conf <configuration file> Override default configuration file
`

// parse args and update config accordingly
func parseArguments() {

	// populate with default values
	var parsedConf = fetchConfig()
	var conf string

	// NB. for now -h and --help are handled by flag itself, we might want to customize that

	var displayVersion bool

	flag.BoolVar(&displayVersion, "version", false, "Display application version")
	flag.BoolVar(&displayVersion, "v", false, "Display application version")

	//
	flag.StringVar(&conf, "conf", confFile, fmt.Sprintf("Override default configuration file (%v)", confFile))
	flag.StringVar(&conf, "c", confFile, fmt.Sprintf("Override default configuration file (%v)", confFile))

	// this needs to be implemented either with the flag.Value interface, or just straight string parsing
	//flag.StringVar(&parsedConf.verbosity, "verbosity", "Output verbosity, accepted values: error, warning, info, debug")

	flag.StringVar(&parsedConf.Rpcs3Path, "rpcs3path", parsedConf.Rpcs3Path, "Set RPCS3 path")
	flag.StringVar(&parsedConf.Rpcs3Path, "rpcs", parsedConf.Rpcs3Path, "Set RPCS3 path")

	flag.StringVar(&parsedConf.PkgDLPath, "dlpath", parsedConf.PkgDLPath, "Set download path")
	flag.StringVar(&parsedConf.PkgDLPath, "dl", parsedConf.PkgDLPath, "Set download path")

	flag.StringVar(&parsedConf.ConfigYMLPath, "configyml", parsedConf.ConfigYMLPath, "Set config.yml path")
	flag.StringVar(&parsedConf.ConfigYMLPath, "yml", parsedConf.ConfigYMLPath, "Set config.yml path")

	flag.StringVar(&parsedConf.XMLCachePath, "xmlcache", parsedConf.XMLCachePath, "XML cache location")
	flag.StringVar(&parsedConf.XMLCachePath, "xml", parsedConf.XMLCachePath, "XML cache location")

	// TODO: short ones for these as well?
	flag.IntVar(&parsedConf.DLTimeout, "timeout", parsedConf.DLTimeout, "Download timeout (in seconds)")
	flag.IntVar(&parsedConf.DLRetries, "retry", parsedConf.DLRetries, "Number of download retries")

	flag.BoolVar(&parsedConf.color, "nocolor", true, "Disable output color")

	flag.Parse()

	fmt.Printf("config file location: %v\n", conf)
	fmt.Printf("parsed args conf: %#v\n\n", parsedConf)

	// this is a special case, if version is requested we only want just that and then bail out
	if displayVersion {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	updateConfig(parsedConf)
}

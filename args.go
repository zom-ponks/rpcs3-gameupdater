/* rpcs3-gameupdater - command line argument parsing */

package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
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

	// NB. for now -h and --help are handled by pflag itself, we might want to customize that

	var displayVersion bool

	pflag.BoolVarP(&displayVersion, "version", "v", false, "Display application version")
	pflag.StringVarP(&conf, "conf", "c", confFile, fmt.Sprintf("Override default configuration file (%v)", confFile))

	// this needs to be implemented either with the pflag.Value interface, or just straight string parsing
	//pflag.StringVar(&parsedConf.verbosity, "verbosity", "Output verbosity, accepted values: error, warning, info, debug")

	pflag.StringVarP(&parsedConf.Rpcs3Path, "rpcs3path", "rpcs", parsedConf.Rpcs3Path, "Set RPCS3 path")
	pflag.StringVarP(&parsedConf.PkgDLPath, "dlpath", "dl", parsedConf.PkgDLPath, "Set download path")
	pflag.StringVarP(&parsedConf.ConfigYMLPath, "configyml", "yml", parsedConf.ConfigYMLPath, "Set config.yml path")
	pflag.StringVarP(&parsedConf.XMLCachePath, "xmlcache", "xml", parsedConf.XMLCachePath, "XML cache location")
	// TODO: short ones for these as well?
	pflag.IntVar(&parsedConf.DLTimeout, "timeout", parsedConf.DLTimeout, "Download timeout (in seconds)")
	pflag.IntVar(&parsedConf.DLRetries, "retry", parsedConf.DLRetries, "Number of download retries")
	pflag.BoolVar(&parsedConf.color, "nocolor", true, "Disable output color")

	pflag.Parse()

	fmt.Printf("config file location: %v\n", conf)
	fmt.Printf("parsed args conf: %#v\n\n", parsedConf)

	// this is a special case, if version is requested we only want just that and then bail out
	if displayVersion {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	updateConfig(parsedConf)
}

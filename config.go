/* rpcs3-gameupdater - maintains app wide configuration and persistence */

package main

import (
	"os"
	"runtime"

	"github.com/go-yaml/yaml"
)

// Verbosity defines debug level
type Verbosity int

// None, Error, Warning, Info, Debug define the debug levels
const (
	None Verbosity = iota
	Error
	Warning
	Info
	Debug
)

// Config is the app wide configuration structure
type Config struct {
	Rpcs3Path     string
	PkgDLPath     string
	ConfigYMLPath string
	DLTimeout     int
	DLRetries     int
	verbosity     Verbosity
	color         bool
}

var conf Config
var confPath string

func fetchConfig() Config {
	// fetch config here
	return conf
}

func persistConfig() {
	// store config here

	return
}

func updateConfig(newconf Config) {
	conf = newconf
}

// this sets up the initial configuration
func initConfig() {
	conf = Config{
		Rpcs3Path:     ".",
		PkgDLPath:     ".",
		ConfigYMLPath: "",
		DLTimeout:     30,
		DLRetries:     3,
		color:         true,
		verbosity:     Debug,
	}

	createConfFile()
	confFile := "/rpcs3/config.yml"
	goos := runtime.GOOS
	switch goos {
	case "freebsd":
		fallthrough
	case "linux":
		if home := os.Getenv("XDG_CONFIG_HOME"); home != "" {
			conf.ConfigYMLPath = home + confFile
		} else if home := os.Getenv("HOME"); home != "" {
			conf.ConfigYMLPath = home + "/.config" + confFile
		} else {
			conf.ConfigYMLPath = "~/.config" + confFile
		}
	case "windows":
		conf.ConfigYMLPath = os.Getenv("RPCS3_CONFIG_DIR") + confFile
	}
	printInfo("config.yml should be at: " + conf.ConfigYMLPath)
	confPath = "."
}

func parseConfFile() {
	//yaml.Unmarshal()
}

func createConfFile() {
	confYAML, err := yaml.Marshal(&conf)
	if isError(err) {
		printError("error: %v", err)
	}
	printInfo(string(confYAML))
	//fmt.printLn(string(confYAML[:len(confYAML)]))
}

// TODO: handle config changes gracefully, from UI and from command line args

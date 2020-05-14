/* rpcs3-gameupdater - maintains app wide configuration and persistence */

package main

import (
	"os"
	"runtime"

	"github.com/pelletier/go-toml"
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
	Rpcs3Path     string    `toml:"rpcs3path"`
	PkgDLPath     string    `toml:"pkgdlpath"`
	ConfigYMLPath string    `toml:"configymlpath"`
	XMLCachePath  string    `toml:"xmlcache"`
	DLTimeout     int       `toml:"timeout"`
	DLRetries     int       `toml:"retries"`
	verbosity     Verbosity `toml:"verbosity"`
	color         bool      `toml:"colorize"`
}

var conf Config
var confFile = "~/.config/rpcs3-gameupdater.toml"

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
		PkgDLPath:     "./Pkgs",
		XMLCachePath:  "./XMLs",
		ConfigYMLPath: "",
		DLTimeout:     30,
		DLRetries:     3,
		color:         true,
		verbosity:     Info,
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
	//confPath = "."
	err := os.MkdirAll(conf.XMLCachePath, 0755)
	if err != nil {
		printError("Error creating the XMLCache folder at %s (errorcode: %s)", conf.XMLCachePath, err)
	}
}

func parseConfFile() {
	//toml.Unmarshal()
}

func createConfFile() {

	/*
		config := Config{Postgres{User: "pelletier", Password: "mypassword", Database: "old_database"}}
		b, err := toml.Marshal(config)
		if err != nil {
		    log.Fatal(err)
		}
		fmt.Println(string(b)) */
	confTOML, err := toml.Marshal(conf)
	if err != nil {
		printError("error: %v", err)
	}
	printDebug(string(confTOML))
}

// TODO: handle config changes gracefully, from UI and from command line args

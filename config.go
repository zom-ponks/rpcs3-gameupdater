/* rpcs3-gameupdater - maintains app wide configuration and persistence */

package main

import (
	"fmt"
	"os"
	"runtime"
)

// Config is the app wide configuration structure
type Config struct {
	Rpcs3Path     string
	PkgDLPath     string
	ConfigYMLPath string
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

// this sets up the initial configuration
func initConfig() {
	conf = Config{
		Rpcs3Path:     ".",
		PkgDLPath:     ".",
		ConfigYMLPath: "",
	}
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
		fmt.Println(conf.ConfigYMLPath)
	case "windows":
		conf.ConfigYMLPath = os.Getenv("RPCS3_CONFIG_DIR") + confFile
	}

	confPath = "."
}

// TODO: handle config changes gracefully, from UI and from command line args

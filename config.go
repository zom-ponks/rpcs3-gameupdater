/* rpcs3-gameupdater - maintains app wide configuration and persistence */

package main

// Config is the app wide configuration structure
type Config struct {
	Rpcs3Path string
	PkgDLPath string
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
		Rpcs3Path: ".",
		PkgDLPath: ".",
	}
	confPath = "."
}

// TODO: handle config changes gracefully, from UI and from command line args

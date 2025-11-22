package main

import (
	"os"
	"strings"
	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	port int
	icecast string
	config gomn.Map
	library gomn.Map
	running = []int{}
	useExternalLib bool
	externalLib gomn.Map
	icecastDomain string
	logLevel = log.DebugLevel
	projectName = "[insert clever radio server name here]"
)

func init() {
	log.Infof("initializing %s...", projectName)
	log.SetLevel(log.DebugLevel)

	log.Debug("reading config...")
	readConf()

	//set log level
	log.SetLevel(logLevel)
}

func main() {
	//start stream server
	initStream()

	//start web server
	initWeb()

	//if the web server stops,
	//  the whole program stops
	log.Info("stopping...")
}

func validateConfig() bool {
	ok := true
	switch (icecast) {
	case "icecast://source:[password]@[ip]:[icecast port]", "", " ":
		log.Error("icecast interface url not set")
		log.Error("please set it in your 'config.gomn' file")
		log.Error("format:  'icecast://source:[password]@[ip]:[port]'")
		ok = false
	}
	
	switch (icecastDomain) {
	case "https://[your icecast domain]", "", " ":
		log.Error("icecast domain not set")
		log.Error("please set it in your 'config.gomn' file")
		log.Error("format:  'https://[your icecast domain]'")
		ok = false
	}

	if useExternalLib {
		log.Info("external libraries enabled")
	}
	
	if len(library) <= 0 {
		log.Error("library is empty")
		log.Error("please set your library in your 'config.gomn' file")
		ok = false
	}

	return ok
}

func readConf() {
	var err error  //these get around a
	var ok bool    //  bug in golang

	//read config file
	var configBytes []byte
	if configBytes, err = os.ReadFile("config.gomn"); err != nil {
		log.Fatalf("failed to read library:  %v", err)
	} else { log.Debug("success reading library") }
	
	//get the config file map
	log.Debug("parsing config") 
	if config, err = gomn.Parse(string(configBytes)); err != nil {
		log.Fatalf("failed to read parse config:  %v", err)
	} else { log.Debug("success parsing config") }
	
	//get the library map
	log.Debug("parsing library")
	if library, ok = config["library"].(gomn.Map); !ok {
		log.Fatal("failed parsing library")
	} else { log.Debug("failed parsing library") }
	
	//get the external library bool
	log.Debug("getting external radio bool")
	if externalLib, ok = config["external library"].(gomn.Map); !ok {
		log.Fatal("failed to assert web server integer")
	} else { log.Debug("set external libraries") }

	//replace the config file map with just config map
	log.Debug("separating config from library")
	if config, ok = config["config"].(gomn.Map); !ok {
		log.Fatal("failed to separate config from library")
	} else { log.Debug("library and config separated") }
	
	//get the log level
	log.Debug("getting log level")
	var logLvl string
	if logLvl, ok = config["log level"].(string); !ok {
		log.Fatal("failed to get log level")
	} else {
		switch strings.ToLower(logLvl) {
		case "info", "i":
			logLevel = log.InfoLevel
		case "debug", "d":
			logLevel = log.DebugLevel
			log.Debug("enabling log.SetReportCaller()")
			log.SetReportCaller(true)
		case "error", "e": 
			logLevel = log.ErrorLevel
		case "warn", "w":
			logLevel = log.WarnLevel
		case "fatal", "f":
			logLevel = log.FatalLevel
		default:
			logLevel = log.InfoLevel
		}

		log.SetLevel(logLevel)
		log.Info("log level set")
	}

	//get the icecast server interface
	log.Debug("setting icecast server interface url")
	if icecast, ok = config["icecast interface"].(string); !ok {
		log.Fatal("failed to assert icecast interface string")
	} else { log.Debug("icecast interface url set") }
	
	//get the icecast server domain
	log.Debug("setting icecast server domain")
	if icecastDomain, ok = config["icecast domain"].(string); !ok {
		log.Fatal("failed to assert icecast domain string")
	} else { log.Debug("icecast domain set") }

	//get the web server port
	log.Debug("getting web server port")
	if port, ok = config["web server port"].(int); !ok {
		log.Fatal("failed to assert web server integer")
	} else { log.Debug("web server port set") }

	//get the external library bool 
	log.Debug("getting external radio bool")
	if useExternalLib, ok = config["enable external radios"].(bool); !ok {
		log.Fatal("failed to assert web server integer")
	} else { log.Debug("external radio bool set") }

	//config checks
	if ok = validateConfig(); !ok {
		log.Fatal("invalid configuration")
	}
}

package main

import (
	"os"
	"fmt"
	"strings"
	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	config gomn.Map
	library gomn.Map
	externalLib gomn.Map
	icecastDomain string
	useExternalLib = true
)

func main() {
	var err error
	var ok bool
	var configRaw []byte
  if configRaw, err = os.ReadFile("config.gomn"); err != nil {
		log.Fatal(err)
	} else { log.Debug("success") }
	if config, err = gomn.Parse(string(configRaw)); err != nil {
		log.Fatal(err)
	} else { log.Debug("success") }
	
	if library, ok = config["library"].(gomn.Map); !ok {
		log.Fatal("failed to assert config[\"library\"] to gomn.Map")
	} else { log.Debug("success") }
	
	if externalLib, ok = config["external library"].(gomn.Map); !ok {
		log.Fatal("failed to assert config[\"library\"] to gomn.Map")
	} else { log.Debug("success") }


	if config, ok = config["config"].(gomn.Map); !ok {
		log.Fatal("failed to assert config[\"config\"] to gomn.Map")
	} else { log.Debug("success") }
	
	if icecastDomain, ok = config["icecast domain"].(string); !ok {
		log.Fatal("failed to assert config[\"icecast domain\"] to gomn.Map")
	} else { log.Debug("success") }

	fmt.Println(string(buildJSONsettings()))
}

func buildJSONsettings() []byte {
	res := []string{
		"[", "  \"config\": {",
	}
	for keyRaw, valRaw := range config {
		var key, valStr, object string
		var ok, valBool bool;
		var valInt int
		if key, ok = keyRaw.(string); !ok {
			log.Fatal("failed type assertion of key to string")
		} else { object += "    \""+key+"\": " }

		if valStr, ok = valRaw.(string); !ok {
			if valInt, ok = valRaw.(int); !ok {
				if valBool, ok = valRaw.(bool); !ok {
					log.Fatal("failed type assertion of value to bool")
				} else { object += fmt.Sprint(valBool)+"," }
			} else { object += fmt.Sprint(valInt)+"," }
		} else { object += "\""+valStr+"\"," }

		res = append(res, object)
	}; res = delLastCharInSlice(res)

	//close object
	res = append(res, "  },")

	//create library JSON object
	res = append(res, "  \"library\":  [")

	for nameRaw, pathRaw := range library {
		object := "    [\n" 
		var name, path string
		var ok bool
		if name, ok = nameRaw.(string); !ok {
			log.Fatal("failed to type assert name of station")
		}; object += "      \""+name+"\",\n"
	
		url := icecastDomain+"/"+name+".ogg"
		object += "      \""+url+"\",\n"

		if path, ok = pathRaw.(string); !ok {
			log.Fatal("failed to type assert path of station")
		}
		object += "      \""+path+"\",\n"

		object += "    ],"
		res = append(res, object)
	}
	for nameRaw, urlRaw := range externalLib {
		object := "    [\n" 
		var name, url string
		var ok bool
		if name, ok = nameRaw.(string); !ok {
			log.Fatal("failed to type assert name of station")
		}; object += "      \""+name+"\",\n"

		if url, ok = urlRaw.(string); !ok {
			log.Fatal("failed to type assert path of station")
		}
		object += "      \""+url+"\"\n"
		
		object += "      null\n"

		object += "    ],"
		res = append(res, object)
	}; delLastCharInSlice(res)
	
	res = append(res, "  ],")

	//remove last char from last line (a comma) 
	// and close json
	res = delLastCharInSlice(res)
	res = append(res, "]")

	return []byte(strings.Join(res, "\n"))
}

func delLastCharInSlice(slice []string) []string {
	lastItm := len(slice)-1
	lastItmLen := len(slice[lastItm])
	slice[lastItm] = slice[lastItm][:lastItmLen-1]
	return slice
}

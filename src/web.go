package main
	
import (
	"net"
	"time"
	"bytes"
	"strings"
	"strconv"
	"net/http"
	"io/ioutil"
	"github.com/charmbracelet/log"
)

func initWeb() {
  log.Info("starting web server")

  http.HandleFunc("/action", actionHandler)
  http.HandleFunc("/", webInterface)

	log.Debug("getting device addresses")
  ipAddressArray, err := net.InterfaceAddrs()
  if err != nil {
    log.Errorf("err detecting ip address:  %v", err)
  } else { log.Debug("device addresses found")
	
		log.Debug("filtering addresses")
		for _, ipAddress := range ipAddressArray {
	    if ipNet, ok := ipAddress.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
	      if ipNet.IP.To4() != nil {
	        log.Infof("listening on http:%s:%s", ipNet.IP, strconv.Itoa(port))
	      } 
	    }
	  }	
		log.Debug("starting web server func")
	  log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
	}
}

func webInterface(w http.ResponseWriter, r *http.Request) {
	var pageCont []byte;var err error
	curTime := time.Now()
  
	reqPage := r.URL.Path
	switch reqPage {
	case "/":
    reqPage = "web/index.html"
	case "/settings.json":
    pageCont = buildJSONsettings()
	case "/library.json":
		pageCont = buildJSONlibrary()
	default:
    reqPage = "web/" + reqPage[1:]
  }
	
	if pageCont == nil {
	  if pageCont, err = ioutil.ReadFile(reqPage); err != nil {
	    log.Errorf("err reading file for requested webpage:  %v", err)
	  }
	}
    
  log.Debugf("requestedPage:  %s", reqPage)

  http.ServeContent(w, r, reqPage, curTime, bytes.NewReader(pageCont))
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
  log.Debug("recieved action request")

  w.Write([]byte("recieved"));

  reqAct := r.Header.Get("action")
  todo := r.Header.Get("do")

  if reqAct == "settings" {
		log.Infof("requestedAction == \"%s\"\ndoThing == \"%s\"", reqAct, todo)
  } else { log.Warn("attempted to action does not exist.") }
}


//because I dislike the current
//  JSON interfacing for Go
func buildJSONlibrary() []byte {
	var ok bool

	jsonStr := []string{"["}

	for nameRaw, _ := range library {
		//start object
		jsonStr = append(jsonStr, "  [")

		var name string
		if name, ok = nameRaw.(string); !ok {
			log.Fatal("failed to assert name to string")
		} else {
			log.Debug("success asserting type")

			url := "    \""+icecastDomain+"/"+name+".ogg\""
			//add values
			jsonStr = append(jsonStr, "    \""+name+"\",")
			jsonStr = append(jsonStr, url)
		}

		//end object
		jsonStr = append(jsonStr, "  ],")
	}

	if useExternalLib {
		for nameRaw, urlRaw := range externalLib {
			//start object
			jsonStr = append(jsonStr, "  [")

			log.Debug("asserting external library name to string")
			var name, url string
			if name, ok = nameRaw.(string); !ok {
				log.Fatal("failed to assert external library name")
			} else { log.Debug("success asserting type") }

			log.Debug("asserting external library url to string")
			if url, ok = urlRaw.(string); !ok {
				log.Fatal("failed to assert external library url")
			} else { log.Debug("success asserting type") }

			//add values
			jsonStr = append(jsonStr, "    \""+name+"\",")
			jsonStr = append(jsonStr, "    \""+url+"\"")

			//end object
			jsonStr = append(jsonStr, "  ],")
		}
	}

	//remove last char from last line (a comma)
	lastItm := len(jsonStr)-1
	lastItmLen := len(jsonStr[lastItm])
	jsonStr[lastItm] = jsonStr[lastItm][:lastItmLen-1]

	//end json
	jsonStr = append(jsonStr, "]")

	//turn []string{} json to []byte
	jsonByte := []byte(strings.Join(jsonStr, "\n"))

	return jsonByte
}

func buildJSONsettings() []byte {
	res := []string{
		"{", "  \"config\": {",
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
				} else { object += strconv.FormatBool(valBool)+"," }
			} else { object += strconv.Itoa(valInt)+"," }
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
		object += "      \""+path+"\"\n"

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
		object += "      \""+url+"\",\n"
		
		object += "      null\n"

		object += "    ],"
		res = append(res, object)
	}; delLastCharInSlice(res)
	
	res = append(res, "  ],")

	//remove last char from last line (a comma) 
	// and close json
	res = delLastCharInSlice(res)
	res = append(res, "}")

	return []byte(strings.Join(res, "\n"))
}

func delLastCharInSlice(slice []string) []string {
	lastItm := len(slice)-1
	lastItmLen := len(slice[lastItm])
	slice[lastItm] = slice[lastItm][:lastItmLen-1]
	return slice
}

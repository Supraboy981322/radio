package main
	
import (
	"net"
	"time"
	"bytes"
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
	curTime := time.Now()
  
	reqPage := r.URL.Path
	switch reqPage {
	case "/":
    reqPage = "web/index.html"
	case "/settings.json":
    reqPage = "settings.json"
	case "web/":
    reqPage = "web/" + reqPage[1:]
	default:
		reqPage = "web/index.html"
  }

	var pageCont []byte;var err error
  if pageCont, err = ioutil.ReadFile(reqPage); err != nil {
    log.Errorf("err reading file for requested webpage:  %v", err)
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

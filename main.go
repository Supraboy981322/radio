package main

import (
    "fmt"
    "bytes"
    "time"
    "net/http"
    "io/ioutil"
    "net"
	"os/exec"
//	"os"
//	"encoding/json"
    //exclusively used for http.ListenAndServe so I can
    //  write one less if err != nil { ... }
    "log"
)

func webInterface(w http.ResponseWriter, r *http.Request) {
    var requestedPage string
    if r.URL.Path == "/" {
        requestedPage = "web/index.html"
    } else if r.URL.Path == "/settings.json" {
        requestedPage = "settings.json"
    } else {
        requestedPage = "web/" + r.URL.Path[1:]
    }

    webpageContent, err := ioutil.ReadFile(requestedPage)
    if err != nil {
        fmt.Errorf("err reading file for requested webpage:  \n", err)
    }
    
    fmt.Printf("requestedPage:  %s\n", requestedPage)

    http.ServeContent(w, r, requestedPage, time.Now(), bytes.NewReader(webpageContent))
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("recieved action request")
    w.Write([]byte("recieved"));
    requestedAction := r.Header.Get("action")
    doThing := r.Header.Get("do")
    if requestedAction == "settings" {
        fmt.Printf("requestedAction == \"%s\"\ndoThing == \"%s\"\n", requestedAction, doThing)
    } else {
        fmt.Errorf("attempted to action does not exist.")
    }
}

func main() {

    http.HandleFunc("/action", actionHandler)
    http.HandleFunc("/", webInterface)

    ipAddressArray, err := net.InterfaceAddrs()
    if err != nil {
        fmt.Errorf("err detecting ip address:  \n", err)
    }
    
    port := "4845"
    
    for _, ipAddress := range ipAddressArray {
        if ipNet, ok := ipAddress.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
            if ipNet.IP.To4() != nil {
                fmt.Printf("listening on http:%s:%s\n", ipNet.IP, port)
            }
        }
    }
	stream("redacted .ogg file path:", "[redacted icecast server address]")
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func stream(file string, url string) {
	cmdArgs := []string{
		"-re",
		"-i", file,
		"-c:a", "copy",
		"-content_type", "audio/ogg",
		"-f", "ogg",
		url,
	}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	
	err := cmd.Start()
	if err != nil {
		fmt.Errorf("failed to start FFmpeg: %v\n", err)
	}

	fmt.Printf("FFmpeg process started.\n")

	err = cmd.Wait()
	if err != nil {
		fmt.Errorf("FFmpeg process exited with error code: %v\n", err)
	}
	
	fmt.Printf("FFmpeg process finished.\n")
}

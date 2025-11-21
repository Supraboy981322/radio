package main

import (
//	"github.com/u2takey/ffmpeg-go"
	"os"
	"io"
	"strings"
	"os/exec"
	"math/rand"
	"path/filepath"
	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	projectName = "[insert clever radio server name here]"
	running = []int{}
	library gomn.Map
	config gomn.Map
	icecast string
	logLevel = log.DebugLevel
)


func quit() {	os.Exit(0) }

func init() {
	log.Infof("initializing %s...", projectName)
	log.SetLevel(log.DebugLevel)

	log.Debug("reading config...")
	readConf()
	log.SetLevel(logLevel)
}

func main() {
	var ok bool
	log.Debug("starting streams goroutines")
	for nameRaw, dirRaw := range library {
		var dir, name string
		if name, ok = nameRaw.(string); !ok {
			log.Fatal("failed to assert name to string")
		}; if dir, ok = dirRaw.(string); !ok {
			log.Fatal("failed to assert filename to string")
		}
		go stream(dir, icecast+"/"+name)
	}
	select{}
	log.Info("stopping...")
}

func stream(dir string, url string) {
	proc := len(running)
	running = append(running, proc)

	file := pickFile(dir)

	format := filepath.Ext(file)[1:]

	log.Debugf("file extention:  %s", format)

	log.Debug("stream() started")

	tranArgs := []string{
		"-loglevel", "panic",  //basically no log
		"-re", //real-time
		"-i", file, //filepath
		"-c:a", "libvorbis", //codec
		"-b:a", "128k", //bit-rate
		"-ar", "44100", //sample rate
		"-ac", "2", //audio channels
		"-f", format, //format (file ext.)
		"pipe:1", //pipe to stdout
	}

	stremArgs := []string{
		"-loglevel", "panic",  //basically no log
		"-i", "pipe:0", //read stdin
		"-content_type", "audio/ogg", //header
		"-f", "ogg", //filetype
		url+".ogg", //icecast url
	}

	log.Debug("creating transcoder cmd")
	tran := exec.Command("ffmpeg", tranArgs...)

	log.Debug("creating stream cmd")
	strem := exec.Command("ffmpeg", stremArgs...)

	log.Debug("creating pipe")
	reader, writer := io.Pipe()

	log.Debug("connecting pipe between transcoder and stream")
	tran.Stdout = writer
	strem.Stdin = reader
	strem.Stdout = os.Stdout
	strem.Stderr = os.Stderr
	tran.Stderr = os.Stderr

	log.Debug("starting transcoder...")
	if err := tran.Start(); err != nil {
		log.Errorf("failed to start FFmpeg transcoder:  %v\n", err)
	} else { log.Info("transcoder process started.") }

	log.Debug("starting stream...")
	if err := strem.Start(); err != nil {
		log.Errorf("failed to start FFmpeg stream:  %v\n", err)
	} else { log.Info("stream process started") }

	log.Debug("waiting for stream to complete")
	if err := strem.Wait(); err != nil {
		log.Errorf("FFmpeg process exited with error code: %v\n", err)
	} else { log.Info("FFmpeg process finished.\n") }

	log.Debug("closing pipe")
	writer.Close()

	log.Debugf("removing %d from running list", proc)
	running = append(running[:proc], running[proc+1:]...)

	log.Debugf("proc: %d ended", proc)

	go stream(dir, url) 
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

	//replace the config file map with just config map
	log.Debug("separating config from library")
	if config, ok = config["config"].(gomn.Map); !ok {
		log.Fatal("failed to separate config from library")
	} else { log.Debug("library and config separated") }
	
	log.Debug("getting log level")
	var logLvl string
	if logLvl, ok = config["log level"].(string); !ok {
		log.Fatal("failed to get log level")
	} else {
		switch strings.ToLower(logLvl) {
		case "info":
			logLevel = log.InfoLevel
		case "debug":
			logLevel = log.DebugLevel
		case "error": 
			logLevel = log.ErrorLevel
		case "warn":
			logLevel = log.WarnLevel
		case "fatal":
			logLevel = log.FatalLevel
		}

		log.SetLevel(logLevel)
		log.Info("log level set")
	}

	//get the icecast server
	log.Debug("setting icecast server url")
	if icecast, ok = config["icecast"].(string); !ok {
		log.Fatal("failed to assert icecast server string")
	} else { log.Debug("icecast server url set") }


	//config checks
	if ok = validateConfig(); !ok {
		log.Fatal("invalid configuration")
	}
}

func validateConfig() bool {
	ok := true
	switch (icecast) {
	case "icecast://source:[password]@[ip]:[port]", "", " ":
		log.Error("icecast url not set")
		log.Error("please set it in your 'config.gomn' file")
		log.Error("format:  'icecast://source:[password]@[ip]:[port]'")
		ok = false
	}
	
	if len(library) <= 0 {
		log.Error("library is empty")
		log.Error("please set your library in your 'config.gomn' file")
		ok = false
	}

	return ok
}

func pickFile(dir string) string {
	var file string
	if files, err := os.ReadDir(dir); err != nil {
		log.Fatalf("failed to read directory:  %v", err)
	} else {
		if len(files) <= 0 {
			log.Fatalf("dir empty:  %s", dir)
		}
		ranInt := rand.Intn(len(files))
		picked := files[ranInt]
		pickedPath := filepath.Join(dir, picked.Name())
		if picked.IsDir() {
			pickFile(pickedPath)
		} else {
			file = pickedPath
		}
	}
	return file
}

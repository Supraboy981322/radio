package main

import (
	"os"
	"io"
	"os/exec"
	"math/rand"
	"path/filepath"
	"github.com/charmbracelet/log"
)

func initStream() {
	var ok bool
	log.Debug("starting streams goroutines")
	for nameRaw, dirRaw := range library {
		var dir, name string
		if name, ok = nameRaw.(string); !ok {
			log.Fatal("failed to assert name to string")
		}; if dir, ok = dirRaw.(string); !ok {
			log.Fatal("failed to assert filename to string")
		}

		log.Debug("starting goroutine")
		//start goroutine
		go stream(dir, icecast+"/"+name)
	}

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

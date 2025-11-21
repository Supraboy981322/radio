package main

import (
	"github.com/charmbracelet/log"
)

func stream(file string, url string) {
	cmdArgs := []string{
		"-re",
		"-i", file,
		"-c:a", "copy",
		"-content_type", "audio/opus",
		"-f", "ogg",
		url,
	}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	
	err := cmd.Start()
	if err != nil {
		log.Errorf("failed to start FFmpeg: %v\n", err)
	}

	log.Infof("FFmpeg process started.\n")

	err = cmd.Wait()
	if err != nil {
		log.Errorf("FFmpeg process exited with error code: %v\n", err)
	}
	
	log.Info("FFmpeg process finished.\n")
}

package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"

	"github.com/ONSdigital/dp-zebedee-utils/content"
	"github.com/ONSdigital/log.go/log"
)

func main() {
	baseDir := flag.String("dir", "", "the zebedee master dir")
	targetType := flag.String("type", "", "the page type to count")
	flag.Parse()

	if !content.Exists(*baseDir) {
		errExit(errors.New("master dir does not exist"))
	}

	if *targetType == "" {
		errExit(errors.New("page type not specified"))
	}

	masterDir := filepath.Join(*baseDir, "master")

	log.Event(nil, "running count job for pageType: "+*targetType)
	if err := content.FilterAndProcess(masterDir, &CounterFiles{count: 0, pageType: *targetType}); err != nil {
		errExit(err)
	}
}

func errExit(err error) {
	log.Event(nil, "Filter and process script returned an error", log.Error(err))
	os.Exit(1)
}

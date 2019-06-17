package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"

	"github.com/ONSdigital/dp-zebedee-utils/collections"
	"github.com/ONSdigital/dp-zebedee-utils/content"
	"github.com/ONSdigital/log.go/log"
)

func main() {
	baseDir, collectionName, pageType := getFlags()

	if !content.Exists(baseDir) {
		errExit(errors.New("master dir does not exist"))
	}

	collectionsDir := filepath.Join(baseDir, "collections")
	masterDir := filepath.Join(baseDir, "master")

	if collectionName == "" {
		errExit(errors.New("no collection name provided"))
	}

	fixC := collections.New(collectionsDir, collectionName)
	if err := collections.Save(fixC); err != nil {
		errExit(err)
	}

	allCols, err := collections.GetCollections(collectionsDir)
	if err != nil {
		errExit(err)
	}

	job := &fixgsiEmails{
		Limit:     3300,
		FixCount:  0,
		FixLog:    make(map[string]int, 0),
		MasterDir: masterDir,
		AllCols:   allCols,
		FixC:      fixC,
		Blocked:   make([]string, 0),
		Type:      pageType,
	}

	if err = content.FilterAndProcess(masterDir, job); err != nil {
		errExit(err)
	}
}

func getFlags() (string, string, string) {
	baseDir := flag.String("dir", "", "the zebedee master dir")
	collectionName := flag.String("col", "", "the name of the collection to add the content to")
	pageType := flag.String("type", "", "the page type to filter by")
	flag.Parse()

	return *baseDir, *collectionName, *pageType
}

func errExit(err error) {
	log.Event(nil, "Filter and process script returned an error", log.Error(err))
	os.Exit(1)
}

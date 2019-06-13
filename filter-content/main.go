package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/ONSdigital/log.go/log"
)

type processFunc func(p *Params, json string, uri string) error

type filterFunc func(p *Params, path string, info os.FileInfo) ([]byte, error)

type doneFunc func(p *Params)

func main() {
	baseDir, collectionName, pageType, limit := getFlags()
	p := NewParams(baseDir, pageType, collectionName, limit)

	count := 0
	logResult := func(p *Params) {
		log.Event(nil, "completed", log.Data{"count": count})
	}

	err := fileWalker(p, countFilesByType(&count, p.TargetType), logResult)

	if err != nil {
		errExit(err)
	}
}

func getFlags() (string, string, string, int) {
	baseDir := flag.String("dir", "", "the zebedee master dir")
	collectionName := flag.String("col", "", "the name of the collection to add the content to")
	pageType := flag.String("type", "any", fmt.Sprintf("the page type to find and fix, use %q if non specific type filtering is required", "any"))
	limit := flag.Int("l", -1, "max size of the collection, -1 means there is no limit")
	flag.Parse()

	return *baseDir, *collectionName, *pageType, *limit
}

func fileWalker(p *Params, process processFunc, onComplete doneFunc) error {
	err := filepath.Walk(p.MasterDir, filterAndProcess(p, pageTypeFilterFunc(p.TargetType), process))
	if err != nil {
		if !reflect.DeepEqual(err, limitReached) {
			errExit(err)
		}
	}
	onComplete(p)
	return nil
}

func filterAndProcess(p *Params, filterFunc filterFunc, processFunc processFunc) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if p.limitReached() {
			return limitReached
		}

		jBytes, err := filterFunc(p, path, info)
		if err != nil {
			return err
		}

		if jBytes == nil {
			return nil
		}
		return processFunc(p, string(jBytes), path)
	}
}

func Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func errExit(err error) {
	log.Event(nil, "error happened", log.Error(err))
	os.Exit(1)
}

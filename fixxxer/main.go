package main

import (
	"flag"
	"github.com/ONSdigital/dp-zebedee-utils/collections"
	"github.com/ONSdigital/log.go/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	oldEmail = "@ons.gsi.gov.uk"
	newEmail = "@ons.gov.uk"
)

type Tracker struct {
	blocked []string
	fixed   int
	skipped int
}

type Err struct {
	Data log.Data
	Err  error
}

func (e Err) Error() string {
	return e.Err.Error()
}

func main() {
	log.Namespace = "fi-xxx-er"
	master, collectionsDir := getConfig()

	t, err := findAndReplace(master, collectionsDir)
	if err != nil {
		errExit(err)
	}

	log.Event(nil, "blocked", log.Data{
		"blocked": len(t.blocked),
		"fixed":   t.fixed,
		"skipped": t.skipped,
	})
}

func getConfig() (string, string) {
	master := flag.String("master", "", "the zebedee master dir")
	collectionsDir := flag.String("collections", "", "the zebedee collections dir")
	flag.Parse()

	if *master == "" {
		errExit(errors.New("master dir not specified"))
	}

	if !Exists(*master) {
		errExit(Err{Err: errors.New("master dir does not exist"), Data: log.Data{"master": *master}})
	}

	if *collectionsDir == "" {
		errExit(errors.New("collections dir not specified"))
	}

	if !Exists(*collectionsDir) {
		errExit(Err{Err: errors.New("collections dir does not exist"), Data: log.Data{"collectionsDir": *collectionsDir}})
	}

	return *master, *collectionsDir
}

func findAndReplace(masterDir string, collectionsDir string) (*Tracker, error) {
	cols, err := collections.GetCollections(collectionsDir)
	if err != nil {
		return nil, err
	}

	log.Event(nil, "scanner master dir for uses of target value", log.Data{"target_value": oldEmail})

	t := &Tracker{
		blocked: make([]string, 0),
		fixed:   0,
		skipped: 0,
	}
	err = filepath.Walk(masterDir, fileWalker(cols, masterDir, t))
	return t, err
}

func fileWalker(cols *collections.Collections, masterDir string, t *Tracker) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if ext := filepath.Ext(info.Name()); ext == ".json" {

			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			raw := string(b)

			if strings.Contains(raw, oldEmail) {

				if strings.Contains(path, "/previous/") ||
					strings.Contains(path, "/datasets/") ||
					strings.Contains(path, "/timeseries/") {
					t.skipped++
					return nil
				}

				uri, err := filepath.Rel(masterDir, path)
				if err != nil {
					return err
				}

				uri = "/" + uri
				for _, c := range cols.Collections {
					if c.Contains(uri) {
						t.blocked = append(t.blocked, uri)
						break
					}
				}
				t.fixed++
			}
		}
		return nil
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
	appErr, ok := err.(Err)
	if ok {
		log.Event(nil, "app error", log.Error(appErr.Err), appErr.Data)
	} else {
		log.Event(nil, "app error", log.Error(err))
	}
	os.Exit(1)
}

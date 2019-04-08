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
	total   int
	blocked []string
	fixed   int
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
		"total_found":           t.total,
		"fixes_applied":         t.fixed,
		"blocked_by_collection": len(t.blocked),
		"outstanding":           t.total - t.fixed,
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

	fixes := collections.New(collectionsDir, "GSIEmailFixes")
	if err := collections.Save(fixes); err != nil {
		return nil, err
	}

	log.Event(nil, "scanner master dir for uses of target value", log.Data{"target_value": oldEmail})

	t := &Tracker{
		blocked: make([]string, 0),
		fixed:   0,
		total:   0,
	}
	err = filepath.Walk(masterDir, fileWalker(cols, masterDir, t, fixes))
	return t, err
}

func fileWalker(cols *collections.Collections, masterDir string, t *Tracker, fixes *collections.Collection) func(path string, info os.FileInfo, err error) error {
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

				if strings.Contains(path, "/previous/") {
					return nil
				}

				t.total++

				if strings.Contains(path, "/datasets/") || strings.Contains(path, "/timeseries/") {
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

				raw = strings.Replace(raw, oldEmail, newEmail, -1)
				if err := fixes.AddReviewedContent(uri, []byte(raw)); err != nil {
					return err
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

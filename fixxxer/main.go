package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ONSdigital/dp-zebedee-utils/collections"
	"github.com/ONSdigital/log.go/log"
	"github.com/pkg/errors"
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
		"blocked_uris":          t.blocked,
		"outstanding":           t.total - t.fixed,
	})
}

func getConfig() (string, string) {
	zebedeeDir := flag.String("zeb", "", "the zebedee base dir")
	flag.Parse()

	if *zebedeeDir == "" {
		errExit(errors.New("zebedee dir not specified"))
	}

	masterDir := path.Join(*zebedeeDir, "master")
	collectionsDir := path.Join(*zebedeeDir, "collections")

	if !Exists(masterDir) {
		errExit(errors.New("master dir does not exist"))
	}

	if !Exists(collectionsDir) {
		errExit(errors.New("collections dir does not exist"))
	}

	return masterDir, collectionsDir
}

func findAndReplace(masterDir string, collectionsDir string) (*Tracker, error) {
	cols, err := collections.GetCollections(collectionsDir)
	if err != nil {
		return nil, err
	}

	fixes := collections.New(collectionsDir, "gsiemailfixes")
	if err := collections.Save(fixes); err != nil {
		return nil, err
	}

	log.Event(nil, "scanning master dir for uses of target value", log.Data{"target_value": oldEmail})

	t := &Tracker{
		blocked: make([]string, 0),
		fixed:   0,
		total:   0,
	}

	err = filepath.Walk(masterDir, fileWalker(cols, masterDir, t, fixes))
	return t, err
}

func fileWalker(collectionsList *collections.Collections, masterDir string, tracker *Tracker, fixCollection *collections.Collection) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if ext := filepath.Ext(info.Name()); ext == ".json" {
			uri, err := filepath.Rel(masterDir, path)
			uri = "/" + uri
			logD := log.Data{"uri": uri}

			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			contentJson := string(b)

			if strings.Contains(contentJson, oldEmail) {

				if strings.Contains(path, "/previous/") {
					log.Event(nil, "skipping previous version content", logD)
					return nil
				}

				tracker.total++

				for _, c := range collectionsList.Collections {
					if blocked := c.Contains(uri); blocked {
						logD["collection"] = c.Name
						log.Event(nil, "cannot fix content as it is contained in another collection", logD)
						tracker.blocked = append(tracker.blocked, uri)
						return nil
					}
				}

				contentJson = strings.Replace(contentJson, oldEmail, newEmail, -1)
				log.Event(nil, "applying content fix", logD)
				if err := fixCollection.AddToReviewed(uri, []byte(contentJson)); err != nil {
					return err
				}
				tracker.fixed++
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
	log.Event(nil, "app error", log.Error(err))
	os.Exit(1)
}

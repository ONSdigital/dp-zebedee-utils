package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/ONSdigital/dp-zebedee-utils/collections"
	"github.com/ONSdigital/log.go/log"
	"github.com/pkg/errors"
)

const (
	oldEmail = "@ons.gsi.gov.uk"
	newEmail = "@ons.gov.uk"
)

var currentDir = ""

type Tracker struct {
	total   int
	blocked []string
	fixed   int
}

func main() {
	log.Namespace = "gsi-email-fix"
	master, collectionsDir := getConfig()

	start := time.Now()
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
		"time":                  time.Now().Sub(start).Seconds(),
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

		uri, err := filepath.Rel(masterDir, path)
		uri = "/" + uri
		logD := log.Data{"uri": uri}

		if ext := filepath.Ext(info.Name()); ext == ".json" {
			if info.Name() == "data.json" || info.Name() == "data_cy.json" {
				b, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				contentJson := string(b)

				if strings.Contains(contentJson, oldEmail) {

					if strings.Contains(path, "/previous/") {
						return nil
					}

					tracker.total++

					for _, c := range collectionsList.Collections {
						if blocked := c.Contains(uri); blocked {
							logD["collection"] = c.Name
							log.Event(nil, "cannot fix content as it is contained in another collection", logD)
							tracker.blocked = append(tracker.blocked, fmt.Sprintf("%s:%s", c.Name, uri))
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
				return nil
			}
			log.Event(nil, "skipping non data.json/data_cy.json json file", logD)
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

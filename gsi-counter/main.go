package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ONSdigital/log.go/log"
)

const (
	oldEmail = "@ons.gsi.gov.uk"
	jsonExt  = ".json"
)

type Tracker struct {
	Datasets   int `json:"datasets"`
	Previous   int `json:"previous"`
	Timeseries int `json:"timeseries"`
	Content    int `json:"contents"`
}

func main() {
	masterDir := flag.String("master", "", "the zebedee master dir")
	flag.Parse()

	if !Exists(*masterDir) {
		errExit(errors.New("master dir does not exist"))
	}

	t := &Tracker{Datasets: 0, Previous: 0, Timeseries: 0, Content: 0}

	err := filepath.Walk(*masterDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if ext := filepath.Ext(info.Name()); ext == jsonExt {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			raw := string(b)
			fmt.Print(".")

			if strings.Contains(raw, oldEmail) {
				if strings.Contains(path, "/previous/") {
					t.Previous += 1
				} else if strings.Contains(path, "/datasets/") {
					t.Datasets += 1
				} else if strings.Contains(path, "/timeseries/") {
					t.Timeseries += 1
				} else {
					t.Content += 1
				}
			}
		}
		return nil
	})

	if err != nil {
		errExit(err)
	}

	log.Event(nil, "gsi scan completed", log.Data{"breakdown": t})
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

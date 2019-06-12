package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ONSdigital/log.go/log"
)

type Page struct {
	PageType string `json:"type"`
}

func main() {
	targetDir := flag.String("dir", "", "the zebedee master dir")
	flag.Parse()

	if !Exists(*targetDir) {
		errExit(errors.New("master dir does not exist"))
	}

	totalCount := 0
	pdfs := 0

	err := filepath.Walk(*targetDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if ext := filepath.Ext(info.Name()); ext != ".json" {
			return nil
		}

		totalCount += 1

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var pt Page
		if err := json.Unmarshal(b, &pt); err != nil {
			return err
		}

		switch pt.PageType {
		case "article", "bulletin", "compendium_landing_page", "compendium_chapter", "static_methodology":
			pdfs += 1
		}
		return nil
	})

	if err != nil {
		errExit(err)
	}

	log.Event(nil, "pdf generating pages", log.Data{"pdf_count": pdfs, "totalCount": totalCount})
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

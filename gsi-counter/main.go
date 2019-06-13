package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	totalFiles := 0
	tick := time.Now()

	pdfPages := map[string]int{
		"article":                 0,
		"bulletin":                0,
		"compendium_landing_page": 0,
		"compendium_chapter":      0,
		"static_methodology":      0,
	}

	err := filepath.Walk(*targetDir, func(path string, info os.FileInfo, err error) error {
		if time.Since(tick) >= time.Second*2 {
			fmt.Print(".")
			tick = time.Now()
		}

		if info.IsDir() {
			return nil
		}

		if ext := filepath.Ext(info.Name()); ext != ".json" {
			return nil
		}

		if info.Name() == "data.json" || info.Name() != "data_cy.json" {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			jsonStr := string(b)

			if !strings.Contains(jsonStr, "@ons.gsi.gov.uk") {
				return nil
			}

			totalFiles += 1

			var pt Page
			if err := json.Unmarshal(b, &pt); err != nil {
				return err
			}

			if count, ok := pdfPages[pt.PageType]; ok {
				pdfPages[pt.PageType] = count + 1
			}
		}
		return nil
	})

	if err != nil {
		errExit(err)
	}

	pdfCount := 0
	for _, val := range pdfPages {
		pdfCount += val
	}

	log.Event(nil, "pdf generating pages", log.Data{
		"pdf_pages":  pdfPages,
		"totalCount": totalFiles,
		"pdf_count":  pdfCount,
	})
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

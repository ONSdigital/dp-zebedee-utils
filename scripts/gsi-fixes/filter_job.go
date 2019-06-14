package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ONSdigital/dp-zebedee-utils/collections"
	"github.com/ONSdigital/dp-zebedee-utils/content"
	"github.com/ONSdigital/log.go/log"
)

var (
	oldEmail = "@ons.gsi.gov.uk"
	newEmail = "@ons.gov.uk"
)

type fixgsiEmails struct {
	MasterDir string
	AllCols   *collections.Collections
	FixC      *collections.Collection
	Limit     int
	FixCount  int
	FixLog    map[string]int
	Blocked   []string
}

func (f *fixgsiEmails) Filter(path string, info os.FileInfo) ([]byte, error) {
	if info.IsDir() {
		return nil, nil
	}

	if strings.Contains(path, "/previous/") {
		return nil, nil
	}

	if info.Name() != "data.json" && info.Name() != "data_cy.json" {
		return nil, nil
	}

	jBytes, err := content.ReadJson(path)
	if err != nil {
		return nil, err
	}

	pageType, err := content.GetPageType(jBytes)
	if err != nil {
		return nil, err
	}

	if pageType.Value != "timeseries" {
		return nil, nil
	}
	return jBytes, nil
}

func (f *fixgsiEmails) Process(jBytes []byte, path string) error {
	jsonStr := string(jBytes)
	uri, err := filepath.Rel(f.MasterDir, path)
	if err != nil {
		return err
	}

	if !strings.Contains(jsonStr, oldEmail) {
		return nil
	}

	uri = "/" + uri
	if blocked, name := f.AllCols.IsBlocked(uri); blocked {
		f.Blocked = append(f.Blocked, fmt.Sprintf("%s:%s", name, uri))
		return nil
	}

	jsonStr = strings.Replace(jsonStr, oldEmail, newEmail, -1)

	log.Event(nil, "applying content fix", log.Data{"uri": uri})
	if err := f.FixC.AddToReviewed(uri, []byte(jsonStr)); err != nil {
		return err
	}

	f.FixCount += 1

	pageType, err := content.GetPageType([]byte(jsonStr))
	if err != nil {
		return err
	}

	f.logFix(pageType)
	return nil
}

func (f *fixgsiEmails) OnComplete() error {
	log.Event(nil, "script fixing timeseries content completed successfully", log.Data{
		"stats":          f.FixLog,
		"fix_count":      f.FixCount,
		"fix_collection": f.FixC.Name,
		"blocked":        f.Blocked,
	})
	return nil
}

func (f *fixgsiEmails) LimitReached() bool {
	if f.Limit == -1 {
		return false
	}
	return f.FixCount >= f.Limit
}

func (f *fixgsiEmails) logFix(pageType *content.PageType) {
	if count, ok := f.FixLog[pageType.Value]; ok {
		f.FixLog[pageType.Value] = count + 1
	} else {
		f.FixLog[pageType.Value] = 1
	}
}

func isPDFPage(pageType *content.PageType) bool {
	switch pageType.Value {
	case "article", "bulletin", "compendium_landing_page", "compendium_chapter", "static_methodology":
		return true
	}
	return false
}

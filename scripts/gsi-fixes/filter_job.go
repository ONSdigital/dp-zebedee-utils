package main

import (
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

type FixNonPDFContent struct {
	MasterDir string
	AllCols   *collections.Collections
	FixC      *collections.Collection
	Limit     int
	FixCount  int
	FixLog    map[string]int
}

func (f *FixNonPDFContent) Filter(path string, info os.FileInfo) ([]byte, error) {
	if info.IsDir() {
		return nil, nil
	}

	if strings.Contains(path, "/previous/") || strings.Contains(path, "/timeseries/") || strings.Contains(path, "/datasets/") {
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

	switch pageType.Value {
	case "article", "bulletin", "compendium_landing_page", "compendium_chapter", "static_methodology":
		return nil, nil
	}
	return jBytes, nil
}

func (f *FixNonPDFContent) Process(jBytes []byte, path string) error {
	jsonStr := string(jBytes)
	uri, err := filepath.Rel(f.MasterDir, path)
	if err != nil {
		return err
	}

	if !strings.Contains(jsonStr, oldEmail) {
		return nil
	}

	uri = "/" + uri
	if f.AllCols.IsBlocked(uri) {
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

func (f *FixNonPDFContent) OnComplete() error {
	log.Event(nil, "script fixing non [previous, timeseries, datasets, article, bulletin, compendium_landing_page, compendium_chapter, static_methodology] content completed successfully", log.Data{
		"stats":          f.FixLog,
		"fix_count":      f.FixCount,
		"fix_collection": f.FixC.Name,
	})
	return nil
}

func (f *FixNonPDFContent) LimitReached() bool {
	if f.Limit == -1 {
		return false
	}
	return f.FixCount >= f.Limit
}

func (f *FixNonPDFContent) logFix(pageType *content.PageType) {
	if count, ok := f.FixLog[pageType.Value]; ok {
		f.FixLog[pageType.Value] = count + 1
	} else {
		f.FixLog[pageType.Value] = 1
	}
}

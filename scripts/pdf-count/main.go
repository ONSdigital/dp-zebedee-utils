package main

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/ONSdigital/dp-zebedee-utils/content"
	"github.com/ONSdigital/log.go/log"
)

type PDFCount struct {
	totalCount int
	typeCount  map[string]int
}

func main() {
	dir := flag.String("dir", "", "")
	flag.Parse()

	if *dir == "" {
		log.Event(nil, "Filter and process script returned an error", log.Error(errors.New("no dir flag specified")))
		os.Exit(1)
	}

	log.Event(nil, "counting PDF generating files")
	if err := content.FilterAndProcess(*dir, &PDFCount{typeCount: make(map[string]int, 0), totalCount: 0}); err != nil {
		log.Event(nil, "Filter and process script returned an error", log.Error(err))
		os.Exit(1)
	}
}

func (c *PDFCount) Filter(path string, info os.FileInfo) (bool, error) {
	if info.IsDir() {
		return false, nil
	}

	if strings.Contains(path, "/previous/") {
		return false, nil
	}

	if info.Name() != "data.json" && info.Name() != "data_cy.json" {
		return false, nil
	}

	jBytes, err := content.ReadJson(path)
	if err != nil {
		return false, err
	}

	pageType, err := content.GetPageType(jBytes)
	if err != nil {
		return false, err
	}

	switch pageType.Value {
	case "article", "bulletin", "compendium_landing_page", "compendium_chapter", "static_methodology":
		return strings.Contains(string(jBytes), "@ons.gsi.gov.uk"), nil
	}
	return false, nil
}

func (c *PDFCount) Process(path string) error {
	jBytes, err := content.ReadJson(path)
	if err != nil {
		return err
	}

	pageType, err := content.GetPageType(jBytes)
	if err != nil {
		return err
	}

	if count, ok := c.typeCount[pageType.Value]; ok {
		c.typeCount[pageType.Value] = count + 1
	} else {
		c.typeCount[pageType.Value] = 0
	}

	c.totalCount += 1
	return nil
}

func (c *PDFCount) OnComplete() error {
	log.Event(nil, "count pdf generating files", log.Data{
		"pageTypes": c.typeCount,
		"total":     c.totalCount,
	})
	return nil
}

func (c *PDFCount) LimitReached() bool {
	return false
}

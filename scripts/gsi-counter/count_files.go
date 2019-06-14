package main

import (
	"os"
	"strings"

	"github.com/ONSdigital/dp-zebedee-utils/content"
	"github.com/ONSdigital/log.go/log"
)

type CounterFiles struct {
	pageType string
	count    int
}

func (c *CounterFiles) Filter(path string, info os.FileInfo) (bool, error) {
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

	if pageType.Value != c.pageType {
		return false, nil
	}
	return strings.Contains(string(jBytes), "@ons.gsi.gov.uk"), nil
}

func (c *CounterFiles) Process(path string) error {
	c.count += 1
	return nil
}

func (c *CounterFiles) OnComplete() error {
	log.Event(nil, "count timeseries contain gsi emails complete", log.Data{
		"type":  c.pageType,
		"found": c.count,
	})
	return nil
}

func (c *CounterFiles) LimitReached() bool {
	return false
}

package main

import (
	"os"
	"strings"

	"github.com/ONSdigital/dp-zebedee-utils/content"
	"github.com/ONSdigital/log.go/log"
)

type Counter struct {
	pageType string
	count    int
}

func (c *Counter) Filter(path string, info os.FileInfo) ([]byte, error) {
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

	if pageType.Value != c.pageType {
		return nil, nil
	}
	return jBytes, nil
}

func (c *Counter) Process(jBytes []byte, path string) error {
	c.count += 1
	return nil
}

func (c *Counter) OnComplete() error {
	log.Event(nil, "count timeseries contain gsi emails complete", log.Data{
		"type":  c.pageType,
		"found": c.count,
	})
	return nil
}

func (c *Counter) LimitReached() bool {
	return false
}

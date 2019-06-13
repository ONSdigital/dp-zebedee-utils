package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type PageType struct {
	Value string `json:"type"`
}

func isIgnored(path string, info os.FileInfo) bool {
	if info.IsDir() {
		return true
	}
	if strings.Contains(path, "/previous/") {
		return true
	}

	if ext := filepath.Ext(info.Name()); ext != ".json" {
		return true
	}

	if info.Name() != "data.json" && info.Name() != "data_cy.json" {
		return true
	}
	return false
}

func pageTypeFilterFunc(targetType string) filterFunc {
	return func(p *Params, path string, info os.FileInfo) ([]byte, error) {
		if isIgnored(path, info) {
			return nil, nil
		}

		jBytes, err := readJson(path)
		if err != nil {
			return nil, err
		}

		pageType, err := getPageType(jBytes)
		if err != nil {
			return nil, err
		}

		if targetType == "any" {
			return jBytes, nil
		}

		if pageType.Value != targetType {
			return nil, nil
		}
		return jBytes, nil
	}
}

func readJson(path string) (b []byte, err error) {
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func getPageType(b []byte) (*PageType, error) {
	var page PageType
	if err := json.Unmarshal(b, &page); err != nil {
		return nil, err
	}
	return &page, nil

}

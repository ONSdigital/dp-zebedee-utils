package main

import (
	"path/filepath"
	"strings"

	"github.com/ONSdigital/log.go/log"
)

const (
	oldEmail = "@ons.gsi.gov.uk"
	newEmail = "@ons.gov.uk"
)

func countFilesByType(count *int, t string) processFunc {
	log.Event(nil, "counting files by type: "+t)
	return func(p *Params, json string, uri string) error {
		*count += 1
		return nil
	}
}

func gsiEmailFixerFunc() processFunc {
	return func(p *Params, jsonStr string, path string) error {
		uri, err := filepath.Rel(p.MasterDir, path)
		if err != nil {
			return err
		}

		if !strings.Contains(jsonStr, "@ons.gsi.gov.uk") {
			return nil
		}

		uri = "/" + uri
		if p.AllCols.IsBlocked(uri) {
			return nil
		}

		jsonStr = strings.Replace(jsonStr, oldEmail, newEmail, -1)

		log.Event(nil, "applying content fix", log.Data{"uri": uri})
		if err := p.FixC.AddToReviewed(uri, []byte(jsonStr)); err != nil {
			return err
		}

		p.FixCount += 1
		return nil
	}
}

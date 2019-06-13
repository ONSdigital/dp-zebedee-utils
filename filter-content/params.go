package main

import (
	"errors"
	"path/filepath"

	"github.com/ONSdigital/dp-zebedee-utils/collections"
)

var limitReached = errors.New("limit reached")

type Params struct {
	MasterDir  string
	TargetType string
	AllCols    *collections.Collections
	FixC       *collections.Collection
	Limit      int
	FixCount   int
}

func NewParams(baseDir string, pageType string, collectionName string, limit int) *Params {
	if !Exists(baseDir) {
		errExit(errors.New("master dir does not exist"))
	}

	collectionsDir := filepath.Join(baseDir, "collections")
	masterDir := filepath.Join(baseDir, "master")

	if collectionName != "" {
		fixC := collections.New(collectionsDir, collectionName)
		if err := collections.Save(fixC); err != nil {
			errExit(err)
		}
	}

	allCols, err := collections.GetCollections(collectionsDir)
	if err != nil {
		errExit(err)
	}

	return &Params{
		MasterDir:  masterDir,
		FixCount:   0,
		Limit:      limit,
		TargetType: pageType,
		AllCols:    allCols,
	}
}

func (p *Params) limitReached() bool {
	if p.Limit == -1 {
		return false
	}
	return p.FixCount >= p.Limit
}

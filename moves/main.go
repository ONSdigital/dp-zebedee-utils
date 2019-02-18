package main

import (
	"context"
	"errors"
	"flag"
	"github.com/ONSdigital/dp-zebedee-utils/collections"
	"github.com/ONSdigital/log.go/log"
	"os"
	"path"
)

const collectionsFolder = "collections"

type Args struct {
	zebRoot        string
	collectionName string
	action         string
}

func main() {
	log.Namespace = "content-mover"
	log.Event(context.Background(), "starting up")

	a, err := getArgs()
	if err != nil {
		log.Event(nil, "missing env var", log.Error(err), log.Data{"var": "ZEB_ROOT"})
		os.Exit(1)
	}

	if a.action == "del" {
		if err := collections.Delete(a.getCollectionsDir(), a.collectionName); err != nil {
			log.Event(nil, "delete collection failed", log.Error(err))
			os.Exit(1)
		}

	} else if a.action == "mk" {
		if err := collections.Create(a.getCollectionsDir(), a.collectionName); err != nil {
			log.Event(nil, "create collection failed", log.Error(err))
			os.Exit(1)
		}
	}
}

func getArgs() (*Args, error) {
	zebRoot := flag.String("zeb_root", "", "")
	collectionName := flag.String("collection", "", "")
	action := flag.String("action", "", "")
	flag.Parse()

	if *zebRoot == "" {
		return nil, errors.New("missing env arg")
	}

	if *collectionName == "" {
		return nil, errors.New("missing env arg")
	}

	return &Args{zebRoot: *zebRoot, collectionName: *collectionName, action: *action}, nil
}

func (a *Args) getCollectionsDir() string {
	return path.Join(a.zebRoot, collectionsFolder)
}

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

type Config struct {
	zebRoot        string
	collectionName string
	src            string
	dest           string
}

func main() {
	log.Namespace = "content-mover"
	log.Event(context.Background(), "starting up")

	a, err := getConfig()
	if err != nil {
		log.Event(nil, "missing env var", log.Error(err), log.Data{"var": "ZEB_ROOT"})
		os.Exit(1)
	}

	col := collections.New(a.getCollectionsDir(), a.collectionName)
	if err := collections.Save(col); err != nil {
		log.Event(nil, "create collection failed", log.Error(err))
		os.Exit(1)
	}

	if err := collections.Move(col, a.getSrc(), a.dest); err != nil {
		log.Event(nil, "move failed", log.Error(err), log.Data{"src": a.getSrc()})
		os.Exit(1)
	}
}

func getConfig() (*Config, error) {
	zebRoot := flag.String("zeb_root", "", "")
	collectionName := flag.String("collection", "", "")
	src := flag.String("src", "", "")
	dest := flag.String("dest", "", "")
	flag.Parse()

	if *zebRoot == "" {
		return nil, errors.New("missing env arg zebedee root")
	}

	if *collectionName == "" {
		return nil, errors.New("missing env arg collection name")
	}

	if *src == "" {
		return nil, errors.New("missing env arg src")
	}

	if *dest == "" {
		return nil, errors.New("missing env arg dest")
	}

	return &Config{
		zebRoot:        *zebRoot,
		collectionName: *collectionName,
		src:            *src,
		dest:           *dest,
	}, nil
}

func (a *Config) getCollectionsDir() string {
	return path.Join(a.zebRoot, "collections")
}

func (a *Config) getMasterDir() string {
	return path.Join(a.zebRoot, "master")
}

func (a *Config) getSrc() string {
	return path.Join(a.zebRoot, "master", a.src)
}

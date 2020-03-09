package config

import (
	"flag"
	"github.com/ONSdigital/dp-zebedee-utils/errs"
	"github.com/ONSdigital/log.go/log"
	"path"
)

type Args struct {
	zebRoot        string
	collectionName string
	reverseChanges bool
}

func (a *Args) GetCollectionsDir() string {
	return path.Join(a.zebRoot, "collections")
}

func (a *Args) GetMasterDir() string {
	return path.Join(a.zebRoot, "master")
}

func (a *Args) GetCollectionName() string {
	return a.collectionName
}

func (a *Args) GetZebedeeDir() string {
	return a.zebRoot
}

func (a *Args) ReverseChanges() bool {
	return a.reverseChanges
}

func GetArgs() (*Args, error) {
	zebRoot := flag.String("zeb_root", "", "The root zebedee directory")
	collectionName := flag.String("collection", "", "The name of the collection to use")
	reverseChanges := flag.Bool("reverse_changes", false, "True flag to create a collection, false to load the collection specified")
	flag.Parse()

	if *zebRoot == "" {
		return nil, errs.New("missing flag", nil, log.Data{"var": "zeb_root"})
	}
	if *collectionName == "" {
		return nil, errs.New("missing flag", nil, log.Data{"var": "collection"})
	}

	return &Args{
		zebRoot:        *zebRoot,
		collectionName: *collectionName,
		reverseChanges: *reverseChanges,
	}, nil
}

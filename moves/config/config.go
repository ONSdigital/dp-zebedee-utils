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
	src            string
	dest           string
	create         bool
}

func (a *Args) GetCollectionsDir() string {
	return path.Join(a.zebRoot, "collections")
}

func (a *Args) GetMasterDir() string {
	return path.Join(a.zebRoot, "master")
}

func (a *Args) GetAbsSrc() string {
	return path.Join(a.zebRoot, "master", a.src)
}

func (a *Args) GetRelSrc() string {
	return a.src
}

func (a *Args) GetDest() string {
	return a.dest
}

func (a *Args) GetCollectionName() string {
	return a.collectionName
}

func (a *Args) GetZebedeeDir() string {
	return a.zebRoot
}

func (a *Args) CreateCollection() bool {
	return a.create
}

func GetArgs() (*Args, error) {
	zebRoot := flag.String("zeb_root", "", "The root zebedee directory")
	collectionName := flag.String("collection", "", "The name of the collection to use")
	create := flag.Bool("create", false, "True flag to create a collection, false to load the collection specified")
	src := flag.String("src", "", "The source taxonomy uri of the content to move")
	dest := flag.String("dest", "", "The destination taxonomy uri to move the content to")
	flag.Parse()

	if *zebRoot == "" {
		return nil, errs.New("missing flag", nil, log.Data{"var": "zeb_root"})
	}

	if *collectionName == "" {
		return nil, errs.New("missing flag", nil, log.Data{"var": "collection"})
	}

	if *src == "" {
		return nil, errs.New("missing flag", nil, log.Data{"var": "src"})
	}

	if *dest == "" {
		return nil, errs.New("missing flag", nil, log.Data{"var": "dest"})
	}

	return &Args{
		zebRoot:        *zebRoot,
		collectionName: *collectionName,
		src:            *src,
		dest:           *dest,
		create:         *create,
	}, nil
}

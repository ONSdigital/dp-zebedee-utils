package config

import (
	"flag"
	"github.com/ONSdigital/log.go/log"
	"os"
	"path"
)

type Args struct {
	zebRoot        string
	collectionName string
	src            string
	dest           string
}

func (a *Args) GetCollectionsDir() string {
	return path.Join(a.zebRoot, "collections")
}

func (a *Args) GetMasterDir() string {
	return path.Join(a.zebRoot, "master")
}

func (a *Args) GetSrc() string {
	return path.Join(a.zebRoot, "master", a.src)
}

func (a *Args) GetDest() string {
	return a.dest
}

func (a *Args) GetCollectionName() string {
	return a.collectionName
}

func GetArgs() *Args {
	zebRoot := flag.String("zeb_root", "", "")
	collectionName := flag.String("collection", "", "")
	src := flag.String("src", "", "")
	dest := flag.String("dest", "", "")
	flag.Parse()

	if *zebRoot == "" {
		log.Event(nil, "missing flag", log.Data{"var": "zeb_root"})
		os.Exit(1)
	}

	if *collectionName == "" {
		log.Event(nil, "missing flag", log.Data{"var": "collection"})
		os.Exit(1)
	}

	if *src == "" {
		log.Event(nil, "missing flag", log.Data{"var": "src"})
		os.Exit(1)
	}

	if *dest == "" {
		log.Event(nil, "missing flag", log.Data{"var": "dest"})
		os.Exit(1)
	}

	return &Args{
		zebRoot:        *zebRoot,
		collectionName: *collectionName,
		src:            *src,
		dest:           *dest,
	}
}

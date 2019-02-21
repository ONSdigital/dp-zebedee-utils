package main

import (
	"fmt"
	"github.com/ONSdigital/dp-zebedee-utils/collections"
	"github.com/ONSdigital/dp-zebedee-utils/moves/config"
	"github.com/ONSdigital/log.go/log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	log.Namespace = "content-mover"
	args := config.GetArgs()

	log.Event(nil, "moving content", log.Data{
		"src":        args.GetSrc(),
		"dest":       args.GetDest(),
		"collection": args.GetCollectionName(),
	})

	cols, err := collections.LoadCollections(args.GetCollectionsDir())
	if err != nil {
		logAndExit(err)
	}

	for _, c := range cols {
		fmt.Println(c.Name)
		if c.Contains(args.GetSrc()) {
			log.Event(nil, "cannot complete move as src is contained in another collections")
			os.Exit(1)
		}
	}

	col := collections.New(args.GetCollectionsDir(), args.GetCollectionName())
	if err := collections.Save(col); err != nil {
		logAndExit(err)
	}

	if err := collections.MoveContent(col, args.GetSrc(), args.GetDest()); err != nil {
		logAndExit(err)
	}

	relPath, _ := filepath.Rel(args.GetMasterDir(), args.GetSrc())
	relPath = path.Join("/", relPath)
	uri, _ := filepath.Split(relPath)

	brokenLinks := collections.FindBrokenLinks(args.GetMasterDir(), filepath.Clean(uri))
	log.Event(nil, "links to fix", log.Data{"brokenLinks": brokenLinks})
}

func logAndExit(err error) {
	if colErr, ok := err.(collections.Error); ok {
		if colErr.OriginalErr != nil {
			log.Event(nil, colErr.Message, log.Error(colErr.OriginalErr), colErr.Data)
		} else {
			log.Event(nil, colErr.Message, colErr.Data)
		}
	} else {
		log.Event(nil, "unknown error", log.Error(err))
	}
	os.Exit(1)
}

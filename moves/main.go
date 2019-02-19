package main

import (
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

	cols := collections.LoadCollections(args.GetCollectionsDir())

	for _, c := range cols {
		if c.Contains(args.GetSrc()) {
			log.Event(nil, "cannot complete move as src is contained in another collections")
			os.Exit(1)
		}
	}

	col := collections.New(args.GetCollectionsDir(), args.GetCollectionName())
	collections.Save(col)

	collections.MoveContent(col, args.GetSrc(), args.GetDest())


	relPath, _ := filepath.Rel(args.GetMasterDir(), args.GetSrc())
	relPath = path.Join("/", relPath)
	uri, _ := filepath.Split(relPath)

	brokenLinks := collections.FindBrokenLinks(args.GetMasterDir(), filepath.Clean(uri))
	log.Event(nil, "links to fix", log.Data{"brokenLinks": brokenLinks})
}

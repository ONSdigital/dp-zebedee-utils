package main

import (
	"github.com/ONSdigital/dp-zebedee-utils/collections"
	"github.com/ONSdigital/dp-zebedee-utils/moves/config"
	"github.com/ONSdigital/log.go/log"
	"os"
	"path/filepath"
)

func main() {
	log.Namespace = "content-mover"
	args := config.GetArgs()

	log.Event(nil, "moving content", log.Data{
		"src":        args.GetRelSrc(),
		"dest":       args.GetDest(),
		"collection": args.GetCollectionName(),
	})

	if args.CreateCollection() {
		log.Event(nil, "creating new collection")

		col := collections.New(args.GetCollectionsDir(), args.GetCollectionName())
		if err := collections.Save(col); err != nil {
			logAndExit(err)
		}
	}

	doMove(args)
}

func doMove(args *config.Args) {
	plan := collections.MovePlan{
		MovingFromAbs: args.GetAbsSrc(),
		MovingFromRel: args.GetRelSrc(),
		MovingToRel:   args.GetDest(),
		MasterDir:     args.GetMasterDir(),
	}

	// find all the pages in master that contain the uri being moved.
	pagesContainingURI, err := collections.FindUsesOfUris(plan)
	if err != nil {
		logAndExit(err)
	}

	// load the existing collections.
	cols, err := collections.GetCollections(args.GetCollectionsDir())
	if err != nil {
		logAndExit(err)
	}

	plan.Collection, err = cols.GetByName(args.GetCollectionName())
	if err != nil {
		logAndExit(err)
	}

	// check that none of the affected files are in another collection
	for _, usage := range pagesContainingURI {
		relURI, err := filepath.Rel(plan.MasterDir, usage)
		if err != nil {
			logAndExit(err)
		}

		blockingCollection := collections.GetCollectionContaining(relURI, cols)
		if blockingCollection != nil && blockingCollection.Name != plan.Collection.Name {
			err := collections.NewErr(
				"cannot proceed with move as affected uri is contained in another collection",
				nil,
				log.Data{
					"collection": blockingCollection,
					"uri":        relURI,
				})
			logAndExit(err)
		}
	}

	// do the move.
	movedUris, err := collections.MoveContent(plan)
	if err != nil {
		logAndExit(err)
	}

	fixedLinks, err := collections.FixUris(plan, pagesContainingURI, movedUris)
	if err != nil {
		logAndExit(err)
	}

	log.Event(nil, "content move completed successfully", log.Data{
		"collection":    args.GetCollectionName(),
		"move_src":      args.GetRelSrc(),
		"move_dest":     args.GetDest(),
		"moved_content": movedUris,
		"link_fixes":    fixedLinks,
	})
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

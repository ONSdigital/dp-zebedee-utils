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
	cols, err := collections.LoadCollections(args.GetCollectionsDir())
	if err != nil {
		logAndExit(err)
	}

	// check that none of the affected files are in another collection
	for _, usage := range pagesContainingURI {
		relURI, err := filepath.Rel(plan.MasterDir, usage)
		if err != nil {
			logAndExit(err)
		}
		if err := collections.IsMoveBlocked(relURI, cols); err != nil {
			logAndExit(err)
		}
	}

	// move not blocked so create a new collection for the move
	col := collections.New(args.GetCollectionsDir(), args.GetCollectionName())
	if err := collections.Save(col); err != nil {
		logAndExit(err)
	}

	plan.Collection = col

	// do the move.
	movedUris, err := collections.MoveContent(plan)
	if err != nil {
		logAndExit(err)
	}

	fixedLinks, err := collections.FixUris(plan, pagesContainingURI, movedUris)
	if err != nil {
		logAndExit(err)
	}

	sanitisedMoves := make(map[string]string)
	for src, dest := range movedUris {
		srcRel, _ := filepath.Rel(args.GetMasterDir(), src)
		destRel, _ := filepath.Rel(col.GetInProgress(), dest)
		sanitisedMoves[srcRel] = destRel
	}

	log.Event(nil, "content move completed successfully", log.Data{
		"collection":    args.GetCollectionName(),
		"move_src":      args.GetRelSrc(),
		"move_dest":     args.GetDest(),
		"moved_content": humanizeMoveResults(movedUris, args.GetMasterDir(), col),
		"link_fixes":    humanizeLinkFixResults(fixedLinks, args.GetMasterDir()),
	})
}

func humanizeMoveResults(movedUris map[string]string, masterDir string, col *collections.Collection) map[string]string {
	sanitisedMoves := make(map[string]string)
	for src, dest := range movedUris {
		srcRel, _ := filepath.Rel(masterDir, src)
		destRel, _ := filepath.Rel(col.GetInProgress(), dest)
		sanitisedMoves[srcRel] = destRel
	}
	return sanitisedMoves
}

func humanizeLinkFixResults(linkFixes []string, masterDir string) []string {
	humanized := make([]string, 0)
	for _, raw := range linkFixes {
		s, _ := filepath.Rel(masterDir, raw)
		humanized = append(humanized, s)
	}
	return humanized
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

package collections

import (
	"github.com/ONSdigital/log.go/log"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func MoveContent(plan MovePlan) (map[string]string, error) {
	// from -> to
	completedMoves := make(map[string]string)

	err := filepath.Walk(plan.MovingFromAbs, func(srcFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if plan.MovingFromAbs == srcFilePath || info.IsDir() {
			// skip
			return nil
		}

		completedMoves[srcFilePath] = srcFilePath

		// relative to the root dir of the move.
		uri, _ := filepath.Rel(plan.MovingFromAbs, srcFilePath)
		collectionURI := plan.Collection.inProgressURI(path.Join(plan.MovingToRel, uri))

		dirs, _ := filepath.Split(collectionURI)
		// create any dirs that do not exist
		if err := os.MkdirAll(dirs, filePerm); err != nil {
			return err
		}

		// if not json file just copy to new home.
		if filepath.Ext(srcFilePath) != ".json" {
			err = moveFile(srcFilePath, collectionURI)
		} else {
			err = moveAndFixJson(srcFilePath, plan, collectionURI)
		}
		if err != nil {
			return err
		}
		completedMoves[srcFilePath] = collectionURI
		return nil
	})
	if err != nil {
		return nil, err
	}
	return completedMoves, nil
}

func FindUsesOfUris(p MovePlan) (map[string]string, error) {
	log.Event(nil, "Scanning master for uses of uri", log.Data{"uri": p.MovingFromRel})
	brokenUris := make(map[string]string)

	err := filepath.Walk(p.MasterDir, func(srcFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(srcFilePath) != ".json" {
			// skip.
			return nil
		}

		b, err := ioutil.ReadFile(srcFilePath)
		if err != nil {
			return err
		}

		fileStr := string(b)
		if strings.Contains(fileStr, p.MovingFromRel) {
			brokenUris[srcFilePath] = srcFilePath
		}
		return nil
	})
	return brokenUris, err
}

func FixUris(p MovePlan, affectedFiles map[string]string, completedMoves map[string]string) ([]string, error) {
	brokenLinks := make([]string, 0)
	for _, srcFilePath := range affectedFiles {
		_, alreadyMoved := completedMoves[srcFilePath]
		if alreadyMoved {
			continue
		}

		b, err := ioutil.ReadFile(srcFilePath)
		if err != nil {
			return nil, err
		}

		fileStr := string(b)
		fileStr = strings.Replace(fileStr, p.MovingFromRel, p.MovingToRel, -1)
		relPath, _ := filepath.Rel(p.MasterDir, srcFilePath)

		if err := WriteFileToCollection(p.Collection, relPath, []byte(fileStr)); err != nil {
			return nil, err
		}
		brokenLinks = append(brokenLinks, srcFilePath)
	}
	return brokenLinks, nil
}

func moveFile(srcFilePath string, collectionURI string) error {
	destFile, err := os.Create(collectionURI)
	if err != nil {
		return err
	}
	defer destFile.Close()

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	data := log.Data{"from": srcFilePath, "to": destFile.Name()}
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return newErr("failed to copy", err, data)
	}
	return nil
}

func moveAndFixJson(srcFilePath string, plan MovePlan, collectionURI string) error {
	b, err := ioutil.ReadFile(srcFilePath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(collectionURI, fixURIs(b, plan.MovingFromRel, plan.MovingToRel), filePerm)
}

func fixURIs(fileBytes []byte, movingFrom string, movingTo string) []byte {
	fileStr := string(fileBytes)
	if strings.Contains(fileStr, movingFrom) {
		fileStr = strings.Replace(fileStr, movingFrom, movingTo, -1)
	}
	return []byte(fileBytes)
}

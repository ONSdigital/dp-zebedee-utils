package collections

import (
	"github.com/ONSdigital/log.go/log"
	"github.com/pkg/errors"
	"os"
	"path"
)

const (
	inProgress = "inprogress"
	complete   = "complete"
	reviewed   = "reviewed"
)

func Create(rootPath string, name string) error {
	collection := path.Join(rootPath, name)

	if Exists(collection) {
		return errors.New("cannot create collection as a collection with this name already exists")
	}

	dirs := []string{
		collection,
		path.Join(collection, inProgress),
		path.Join(collection, complete),
		path.Join(collection, reviewed),
	}

	log.Event(nil, "creating collection", log.Data{"dirs": dirs})

	for _, d := range dirs {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			return err
		}
	}
	log.Event(nil, "collection created successfully", log.Data{"collection": collection})
	return nil
}

func Exists(collection string) bool {
	_, err := os.Stat(collection)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func Delete(rootPath string, name string) error {
	target := path.Join(rootPath, name)
	if !Exists(target) {
		return errors.New("cannot delete collection as it does not exist")
	}

	log.Event(nil, "deleting collection", log.Data{"collection": target})
	return os.RemoveAll(target)
}

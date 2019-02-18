package collections

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/log.go/log"
	"io"
	"os"
	"path"
	"path/filepath"
)

func Save(c *Collection) error {
	if Exists(c.Metadata.collectionRoot) {
		return errors.New("cannot create collection as a collection with this name already exists")
	}
	if err := createCollectionDirectories(c); err != nil {
		return err
	}

	if err := createCollectionJson(c); err != nil {
		return err
	}
	log.Event(nil, "collection created successfully", log.Data{"collection": c.Name})
	return nil
}

func Exists(filePath string) bool {
	_, err := os.Stat(filePath)
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

func Move(c *Collection, src string, dest string) error {
	var err error
	var srcFile *os.File = nil
	var destFile *os.File = nil

	if !Exists(src) {
		return errors.New("src does not exist")
	}

	defer func() {
		for _, c := range []io.Closer{srcFile, destFile} {
			if c != nil {
				if err := c.Close(); err != nil {
					panic(err)
				}
			}
		}
	}()

	srcFile, err = os.Open(src)
	if err != nil {
		return err
	}

	destFile, err = c.CreateFile(dest)
	if err != nil {
		return err
	}

	wr := bufio.NewWriter(destFile)
	_, err = io.Copy(wr, srcFile)
	wr.Flush()
	return err
}

func createCollectionDirectories(c *Collection) error {
	for _, d := range c.getDirs() {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			return err
		}
	}
	log.Event(nil, "collection directories created successfully")
	return nil
}

func createCollectionJson(c *Collection) error {
	f, err := os.Create(c.Metadata.collectionJSON)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Event(nil, "failed to close collection json file", log.Error(err))
			panic(err)
		}
	}()

	b, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	w.Write(b)
	w.Flush()
	return nil
}

func (c *Collection) CreateFile(dest string) (*os.File, error) {
	relPath := path.Join(c.Metadata.inProgress, dest)
	dir, _ := filepath.Split(relPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(relPath)
}

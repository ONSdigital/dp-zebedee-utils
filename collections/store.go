package collections

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/log.go/log"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Save(c *Collection) {
	if Exists(c.Metadata.CollectionRoot) {
		log.Event(nil, "cannot create collection as a collection with this name already exists", log.Data{"name": c.Name})
		os.Exit(1)
	}

	if err := createCollectionDirectories(c); err != nil {
		log.Event(nil, "error creating collection directories", log.Error(err), log.Data{"name": c.Name})
		os.Exit(1)
	}

	if err := createCollectionJson(c); err != nil {
		log.Event(nil, "error creating collection json", log.Error(err), log.Data{"name": c.Name})
		os.Exit(1)
	}
	log.Event(nil, "collection created successfully", log.Data{"collection": c.Name})
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

func MoveContent(c *Collection, src string, dest string) {
	var err error
	var srcFile *os.File = nil
	var destFile *os.File = nil

	if !Exists(src) {
		log.Event(nil, "failed to move content src does not exit", log.Data{"src": src})
		os.Exit(1)
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
		log.Event(nil, "failed to move content error opening src file", log.Error(err), log.Data{"src": src})
		os.Exit(1)
	}

	relPath := path.Join(c.Metadata.InProgress, dest)
	dir, _ := filepath.Split(relPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Event(nil, "failed to move content error creating dirs in collection", log.Error(err), log.Data{"src": src})
		os.Exit(1)
	}

	destFile, err = os.Create(relPath)
	if err != nil {
		log.Event(nil, "failed to move content error creating dest file", log.Error(err), log.Data{"dest": dest})
		os.Exit(1)
	}

	wr := bufio.NewWriter(destFile)
	_, err = io.Copy(wr, srcFile)
	if err != nil {
		log.Event(nil, "failed to move content error coping from src to dest", log.Error(err),
			log.Data{
				"src":  src,
				"dest": dest,
			})
		os.Exit(1)
	}
	wr.Flush()
}

func FindBrokenLinks(root string, src string) []string {
	filesToFix := make([]string, 0)
	log.Event(nil, "scanning published content for broken links", log.Data{"link": src})

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || info.Name() != "data.json" {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if strings.Contains(string(b), src) {
			filesToFix = append(filesToFix, path)
		}
		return nil
	})

	if err != nil {
		log.Event(nil, "failed to find broken links for content move", log.Error(err))
		os.Exit(1)
	}
	return filesToFix
}

func LoadCollections(collectionsRoot string) []*Collection {
	log.Event(nil, "loading existing collections")
	collectionFiles, err := ioutil.ReadDir(collectionsRoot)
	if err != nil {
		log.Event(nil, "failed to read collections dir", log.Error(err))
		os.Exit(1)
	}

	collections := make([]*Collection, 0)
	for _, f := range collectionFiles {
		if f.IsDir() {
			c, err := LoadCollection(collectionsRoot, f.Name())
			if err != nil {
				log.Event(nil, "failed to load collection", log.Error(err), log.Data{"collectionName": f.Name()})
				os.Exit(1)
			}

			collections = append(collections, c)
		}
	}
	return collections
}

func LoadCollection(collectionsRoot string, name string) (*Collection, error) {
	metadata := NewMetadata(collectionsRoot, name)
	b, err := ioutil.ReadFile(metadata.CollectionJSON)
	if err != nil {
		return nil, err
	}

	var col Collection
	if err := json.Unmarshal(b, &col); err != nil {
		return nil, err
	}
	col.Metadata = metadata
	return &col, nil
}

func createCollectionDirectories(c *Collection) error {
	for _, d := range c.getDirs() {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func createCollectionJson(c *Collection) error {
	f, err := os.Create(c.Metadata.CollectionJSON)
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

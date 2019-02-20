package collections

import (
	"bufio"
	"encoding/json"
	"github.com/ONSdigital/log.go/log"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Save(c *Collection) error {
	if Exists(c.Metadata.CollectionRoot) {
		return newErr("cannot create collection as a collection with this name already exists", nil, log.Data{"name": c.Name})
	}

	if err := createCollectionDirectories(c); err != nil {
		return newErr("error creating collection directories", err, log.Data{"name": c.Name})
	}

	if err := createCollectionJson(c); err != nil {
		return newErr("error creating collection json", err, log.Data{"name": c.Name})
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
		return newErr("cannot delete collection as it does not exist", nil, log.Data{"collection": name})
	}

	log.Event(nil, "deleting collection", log.Data{"collection": target})
	return os.RemoveAll(target)
}

func MoveContent(c *Collection, src string, dest string) error {
	if !Exists(src) {
		return newErr("failed to move content src does not exit", nil, log.Data{"src": src})
	}

	f, err := ioutil.ReadFile(src)
	if err != nil {
		return newErr("failed to move content error reading src file", err, log.Data{"src": src})
	}

	var page map[string]interface{}
	if err := json.Unmarshal(f, &page); err != nil {
		return newErr("failed to unmarshall src json file", err, log.Data{"src": src})
	}

	fullDestPath := path.Join(c.Metadata.InProgress, dest)
	destDir, _ := filepath.Split(fullDestPath)
	newURI, _ := filepath.Split(dest)

	if _, ok := page["uri"]; ok {
		page["uri"] = strings.TrimRight(newURI, "/")
	}

	modified, err := json.MarshalIndent(page, "", " ")
	if err != nil {
		return newErr("failed to marshall modified src json file", err, log.Data{"src": src})
	}

	if err := os.MkdirAll(destDir, 0666); err != nil {
		return newErr("failed to move content error creating dirs in collection", err, log.Data{"src": src})
	}

	_, err = os.Create(fullDestPath)
	if err != nil {
		return newErr("failed to move content error creating dest file", err, log.Data{"dest": dest})
	}

	if err := ioutil.WriteFile(fullDestPath, modified, 0666); err != nil {
		return newErr("failed to write modified content to file", err, log.Data{"dest": dest})
	}
	return nil
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

func LoadCollections(collectionsRoot string) ([]*Collection, error) {
	log.Event(nil, "loading existing collections")
	collectionFiles, err := ioutil.ReadDir(collectionsRoot)
	if err != nil {
		return nil, newErr("failed to read collections dir", err, nil)
	}

	collections := make([]*Collection, 0)
	for _, f := range collectionFiles {
		if f.IsDir() {
			c, err := LoadCollection(collectionsRoot, f.Name())
			if err != nil {
				return nil, newErr("failed to load collection", err, log.Data{"collectionName": f.Name()})
			}
			if c != nil {
				collections = append(collections, c)
			}
		}
	}
	return collections, nil
}

func LoadCollection(collectionsRoot string, name string) (*Collection, error) {
	metadata := NewMetadata(collectionsRoot, name)
	if !Exists(metadata.CollectionJSON) {
		log.Event(nil, "no collection json file exists for collection", log.Data{"collection": name})
		return nil, nil
	}

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
		if err := os.MkdirAll(d, 0666); err != nil {
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

package collections

import (
	"bufio"
	"encoding/json"
	"github.com/ONSdigital/log.go/log"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const filePerm = 0755

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

func MoveContent(c *Collection, currentTaxonomy string, newTaxonomy string) error {
	return filepath.Walk(currentTaxonomy, func(contentAbsolutePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Make the file path relative to the master directory
		// i.e /zebedee/content/master/aboutus/data.json -> /aboutus/data.json
		taxonomyURI, _ := filepath.Rel(currentTaxonomy, contentAbsolutePath)

		// Collection uri is the absolute path of the file in the collection's in progress dir
		// e.g. /zebedee/content/collections/test123/inprogress/aboutus/data.json
		collectionURI := c.inProgressURI(path.Join(newTaxonomy, taxonomyURI))

		// Create the taxonomy directory in the collection
		if info.IsDir() {
			return createContentDirInCollection(collectionURI)
		}

		// otherwise copy the file into the new location within the collection directory
		return createContentFileInCollection(contentAbsolutePath, collectionURI)
	})
	return nil
}

func createContentDirInCollection(collectionURI string) error {
	if err := os.MkdirAll(collectionURI, filePerm); err != nil {
		return err
	}
	return nil
}

func createContentFileInCollection(srcFilePath string, collectionURI string) error {
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


	data := log.Data{"from": srcFilePath, "to": collectionURI}
	log.Event(nil, "copying content", data)

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return newErr("failed to copy", err, data)
	}
	return nil
}

func fixURI(f []byte, newUri string) ([]byte, error) {
	var page map[string]interface{}
	if err := json.Unmarshal(f, &page); err != nil {
		return nil, err
	}

	if _, ok := page["uri"]; ok {
		page["uri"] = strings.TrimRight(newUri, "/")
	}

	modified, err := json.MarshalIndent(page, "", " ")
	if err != nil {
		return nil, err
	}
	return modified, nil
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
		if err := os.MkdirAll(d, filePerm); err != nil {
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

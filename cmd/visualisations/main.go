package main

import (
	"fmt"
	"github.com/ONSdigital/dp-zebedee-utils/cmd/visualisations/config"
	"github.com/ONSdigital/dp-zebedee-utils/collections"
	"github.com/ONSdigital/dp-zebedee-utils/errs"
	"github.com/ONSdigital/log.go/log"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var snippetReplacements = map[string]string{
	"var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;": "// var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;",

	"ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';": "// ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';",

	"var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);": "// var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);",

	"ga.src = ('https:' == document.location.protocol ? 'https://' : 'http://') + 'stats.g.doubleclick.net/dc.js';": "// ga.src = ('https:' == document.location.protocol ? 'https://' : 'http://') + 'stats.g.doubleclick.net/dc.js';",

	"var _gaq = _gaq || [];": "// var _gaq = _gaq || [];",

	"_gaq.push(['_setAccount', 'UA-37894017-1']);": "// _gaq.push(['_setAccount', 'UA-37894017-1']);",

	"_gaq.push(['_trackPageview']);": "// _gaq.push(['_trackPageview']);",

	"ga('create', 'GTM-MBCBVQS');": "// ga('create', 'GTM-MBCBVQS');",

	"ga('send', 'pageview');": "// ga('send', 'pageview');",

	"ga('create', 'GTM-MBCBVQS', {'name': 'section_linger' });": "// ga('create', 'GTM-MBCBVQS', {'name': 'section_linger' });",

	"ga('create', 'UA-42055132-1', {'name': 'shorthand' });": "// ga('create', 'UA-42055132-1', {'name': 'shorthand' });",

	"ga('create', 'UA-37894017-2', 'auto');": "// ga('create', 'UA-37894017-2', 'auto');",

	"ga('shorthand.send', 'pageview', {\n        'dimension1': 24846,\n        'dimension2': 67636,\n        'metric1':24846,\n        'metric2':67636\n      });": "// ga('shorthand.send', 'pageview', {\n        // 'dimension1': 24846,\n        // 'dimension2': 67636,\n        // 'metric1':24846,\n        // 'metric2':67636\n      // });",

	// some have windows line endings
	"ga('shorthand.send', 'pageview', {\r\n        'dimension1': 24846,\r\n        'dimension2': 67636,\r\n        'metric1':24846,\r\n        'metric2':67636\r\n      });": "// ga('shorthand.send', 'pageview', {\n        // 'dimension1': 24846,\n        // 'dimension2': 67636,\n        // 'metric1':24846,\n        // 'metric2':67636\n      // });",
}

type Tracker struct {
	numOfHtmlFiles     int
	dataJsonFilesMoved int
	filesFixed         []string
	blocked            []string
	snippetsReplaceed  map[string]int
}

func main() {

	args, err := config.GetArgs()
	if err != nil {
		logAndExit(err)
	}

	log.Event(nil, "Content move configuration", log.Data{
		"collection":     args.GetCollectionName(),
		"master dir":     args.GetMasterDir(),
		"reverseChanges": args.ReverseChanges(),
	})

	col := collections.New(args.GetCollectionsDir(), args.GetCollectionName())
	if err := collections.Save(col); err != nil {
		logAndExit(err)
	}

	cols, err := collections.GetCollections(args.GetCollectionsDir())
	if err != nil {
		logAndExit(err)
	}

	replaceCodeInVisualisations(args, cols, col)
}

func replaceCodeInVisualisations(args *config.Args, cols *collections.Collections, col *collections.Collection) {

	t := &Tracker{
		filesFixed:        make([]string, 0),
		snippetsReplaceed: make(map[string]int, 0),
		numOfHtmlFiles:    0,
	}

	visualisationDir := path.Join(args.GetMasterDir(), "visualisations")
	err := filepath.Walk(visualisationDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// only html files have the Google Analytics snippets we want to replace
		if ext := filepath.Ext(info.Name()); ext == ".html" {
			err := replaceCodeInHtmlFile(path, t, args, cols, col)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logAndExit(err)
	}

	addDataJsonFilesToCollection(t, args, col)

	log.Event(nil, "Finished", log.Data{
		"numOfHtmlFiles":                t.numOfHtmlFiles,
		"snippetsReplaced":              t.snippetsReplaceed,
		"filesFixed":                    t.filesFixed,
		"dataJsonFilesMoved":            t.dataJsonFilesMoved,
		"numOfFilesFixed":               len(t.filesFixed),
		"numOfFilesBlockedByCollection": len(t.blocked),
		"filesBlockedByCollection":      t.blocked,
	})
}

func addDataJsonFilesToCollection(t *Tracker, args *config.Args, col *collections.Collection) {
	// for each file fixed, add the data.json to the collection
	for _, htmlFile := range t.filesFixed {
		pathSegments := strings.Split(htmlFile, "/")

		dataJsonUri := path.Join(pathSegments[0], pathSegments[1], "data.json")
		dataJsonMasterUri := path.Join(args.GetMasterDir(), dataJsonUri)
		dataJsonReviewedUri := path.Join(col.GetReviewed(), dataJsonUri)

		if !fileExists(dataJsonMasterUri) {
			fmt.Println("Skipping data json not found: " + dataJsonUri)
			continue
		}

		if fileExists(dataJsonReviewedUri) {
			//fmt.Println("Skipping data json already added: " + dataJsonReviewedUri)
			// file already added to the collection
			continue
		}

		fmt.Println("Moving data json from : " + dataJsonMasterUri + " to: " + dataJsonReviewedUri)
		b, err := ioutil.ReadFile(dataJsonMasterUri)
		if err != nil {
			logAndExit(err)
		}

		err = collections.WriteContent(dataJsonReviewedUri, b)
		if err != nil {
			logAndExit(err)
		}

		t.dataJsonFilesMoved++
	}
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func replaceCodeInHtmlFile(masterPath string, t *Tracker, args *config.Args, cols *collections.Collections, col *collections.Collection) error {
	t.numOfHtmlFiles++

	b, err := ioutil.ReadFile(masterPath)
	if err != nil {
		return err
	}
	fileContents := string(b)

	existingCollectionsChecked := false
	fileUpdated := false

	uri, err := filepath.Rel(args.GetMasterDir(), masterPath)
	if err != nil {
		return err
	}

	fmt.Println("Checking file: " + masterPath)
	for k, v := range snippetReplacements {

		searchString := k
		replaceString := v

		if args.ReverseChanges() {
			searchString = v
			replaceString = k
		}

		if strings.Contains(fileContents, searchString) {

			checkIfFileAlreadyExistsInACollection(&existingCollectionsChecked, uri, cols, t)

			t.snippetsReplaceed[k] = t.snippetsReplaceed[k] + 1
			fileUpdated = true
			fmt.Println("   Replacing snippet: " + searchString)
			fileContents = strings.Replace(fileContents, searchString, replaceString, -1)
		}
	}

	if fileUpdated {
		t.filesFixed = append(t.filesFixed, uri)

		dataJsonReviewedUri := path.Join(col.GetReviewed(), uri)

		err = collections.WriteContent(dataJsonReviewedUri, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkIfFileAlreadyExistsInACollection(existingCollectionsChecked *bool, uri string, cols *collections.Collections, t *Tracker) {

	if !*existingCollectionsChecked {
		for _, c := range cols.Collections {
			if c.Contains(uri) {
				t.blocked = append(t.blocked, uri)
				break
			}
		}
		*existingCollectionsChecked = true
	}
}

func logAndExit(err error) {
	if colErr, ok := err.(errs.Error); ok {
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

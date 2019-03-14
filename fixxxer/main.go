package main

import (
	"flag"
	"fmt"
	"github.com/ONSdigital/log.go/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Err struct {
	Data log.Data
	Err  error
}

func (e Err) Error() string {
	return e.Err.Error()
}

func main() {
	master := getConfig()

	uses, err := findUses(master, "@ons.gsi.gov.uk")
	if err != nil {
		errExit(err)
	}

	log.Event(nil, "results", log.Data{"uses": len(uses)})
	log.Event(nil, "subset", log.Data{"uses": uses[:20]})
}

func getConfig() string {
	master := flag.String("master", "", "the zebedee master dir")
	flag.Parse()

	if *master == "" {
		errExit(errors.New("master dir not specified"))
	}

	if !Exists(*master) {
		errExit(Err{Err: errors.New("master dir does not exist"), Data: log.Data{"master": *master}})
	}
	return *master
}

func findUses(masterDir string, targetVal string) ([]string, error) {
	uses := make([]string, 0)
	ticker := func (c chan bool) {
		run := true
		for run {
			select{
				case <- c:
					run = false
			default:
				time.Sleep(time.Second * 3)
				fmt.Print(".")
			}
		}
	}

	log.Event(nil, "scanner master dir for uses of target value", log.Data{"target_value": targetVal})
	exitChan := make(chan bool, 0)
	defer func() {
		exitChan <- true
	}()
	go ticker(exitChan)

	err := filepath.Walk(masterDir, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		if ext := filepath.Ext(info.Name()); ext == ".json" {

			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			if strings.Contains(string(b), targetVal) {
				uses = append(uses, path)
			}
		}
		return nil
	})
	return uses, err
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

func errExit(err error) {
	appErr, ok := err.(Err)
	if ok {
		log.Event(nil, "app error", log.Error(appErr.Err), appErr.Data)
	} else {
		log.Event(nil, "app error", log.Error(err))
	}
	os.Exit(1)
}

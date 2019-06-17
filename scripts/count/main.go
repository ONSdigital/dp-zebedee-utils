package main

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/ONSdigital/dp-zebedee-utils/content"
	"github.com/ONSdigital/log.go/log"
)

func main() {
	dir := flag.String("dir", "", "the zebedee master dir")
	targetTypes := flag.String("types", "", "comma separated list of page types to count")
	flag.Parse()

	if !content.Exists(*dir) {
		errExit(errors.New("master dir does not exist"))
	}

	if *targetTypes == "" {
		errExit(errors.New("page type not specified"))
	}

	types := strings.Split(*targetTypes, ",")
	typeMap := map[string]int{}
	for _, val := range types {
		typeMap[strings.TrimSpace(val)] = 0
	}

	log.Event(nil, "running count job for pageTypes", log.Data{
		"types": typeMap,
		"dir":   *dir,
	})
	if err := content.FilterAndProcess(*dir, &Counter{total: 0, targetTypes: typeMap}); err != nil {
		errExit(err)
	}
}

func errExit(err error) {
	log.Event(nil, "Filter and process script returned an error", log.Error(err))
	os.Exit(1)
}

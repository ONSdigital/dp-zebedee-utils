package main

import (
	"flag"
	"github.com/ONSdigital/dp-zebedee-utils/setup/cms"
	"github.com/ONSdigital/dp-zebedee-utils/setup/log"
	"os"
)

func main() {
	root := flag.String("root", "", "the root directory in which to build zebedeeDir directory structure")
	flag.Parse()

	if *root == "" {
		log.Warn.Printf("please specify %q flag, use -h to see help menu\n", "root")
		os.Exit(1)
	}

	cms.Out = log.InfoHandler
	cms.OutErr = log.ErrorHandler

	cms, err := cms.New(*root)
	if err != nil {
		log.Error.Fatal(err)
		os.Exit(1)
	}

	err = cms.Initialize()
	if err != nil {
		log.Error.Fatal(err)
	}

	log.Info.Println("successfully generated zebedee file system")
}

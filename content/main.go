package main

import (
	"flag"
	"github.com/ONSdigital/dp-zebedee-utils/content/cms"
	"github.com/ONSdigital/dp-zebedee-utils/content/log"
	"github.com/ONSdigital/dp-zebedee-utils/content/scripts"
	"os"
)

func main() {
	root := flag.String("r", "", "the root directory in which to build zebedeeDir directory structure")
	isCMD := flag.Bool("cmd", false, "if true creates a CMD service account, default is false")
	flag.Parse()

	if *root == "" {
		log.Warn.Printf("please specify %q flag, use -h to see help menu\n", "root")
		os.Exit(1)
	}

	cms.Out = log.InfoHandler
	cms.OutErr = log.ErrorHandler

	builder, err := cms.New(*root, *isCMD)
	if err != nil {
		errorAndExit(err)
	}

	err = builder.Build()
	if err != nil {
		errorAndExit(err)
	}

	log.Info.Println("successfully generated zebedee file system")

	var file string
	file, err = scripts.GenerateCMSRunScript(builder)
	if err != nil {
		errorAndExit(err)
	}

	log.Info.Printf("a customized script for running zebedee cms has been generated under %q", file)
}

func errorAndExit(err error) {
	log.Error.Fatal(err)
	os.Exit(1)
}

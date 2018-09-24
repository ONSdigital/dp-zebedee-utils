package main

import (
	"flag"
	"github.com/ONSdigital/dp-zebedee-utils/content/cms"
	"github.com/ONSdigital/dp-zebedee-utils/content/log"
	"os"
)

func main() {
	root := flag.String("r", "", "the root directory in which to build zebedeeDir directory structure")
	flag.Parse()

	if *root == "" {
		log.Warn.Printf("please specify %q flag, use -h to see help menu\n", "root")
		os.Exit(1)
	}

	cms.Out = log.InfoHandler
	cms.OutErr = log.ErrorHandler

	builder, err := cms.New(*root)
	if err != nil {
		log.Error.Fatal(err)
		os.Exit(1)
	}

	err = builder.Build()
	if err != nil {
		log.Error.Fatal(err)
	}

	log.Info.Println("successfully generated zebedee file system")
	log.Info.Printf("add the following to zebedee/run.sh\n\nexport zebedee_root=%q\n", *root)
}

package main

import (
	"flag"
	"os"

	"github.com/ONSdigital/dp-zebedee-utils/content/cms"
	"github.com/ONSdigital/dp-zebedee-utils/content/scripts"
	"github.com/ONSdigital/log.go/log"
)

func main() {
	log.Namespace = "zebedee-content-generator"
	root := flag.String("r", "", "the root directory in which to build zebedee directory structure and unpack the default content")
	zebDir := flag.String("zeb-dir", "", "the root directory path of your zebedee project")
	enableCMD := flag.Bool("enable_cmd", false, "enabled or disabled the CMD features in Zebedee")
	flag.Parse()

	if *root == "" {
		log.Event(nil, "please specify a root dir, use -h to see the help menu")
		os.Exit(1)
	}

	if *zebDir == "" {
		log.Event(nil, "please specify the path to the root of you zebedee project")
		os.Exit(1)
	}

	generateCMSContent(*root, *enableCMD, *zebDir)
}

func generateCMSContent(root string, enableCMD bool, zebDir string) {
	builder, err := cms.New(root, enableCMD)
	if err != nil {
		errorAndExit(err)
	}

	err = builder.GenerateCMSContent()
	if err != nil {
		errorAndExit(err)
	}

	t := builder.GetRunTemplate()

	var file string
	file, err = scripts.GenerateCMSRunScript(t)
	if err != nil {
		errorAndExit(err)
	}

	scriptLocation, err := scripts.CopyToProjectDir(zebDir, file)
	if err != nil {
		errorAndExit(err)
	}
	log.Event(nil, "successfully generated zebedee file structure and default content you can use the generated run-cms.sh file to run the application", log.Data{
		"run_script_location":      scriptLocation,
		cms.EnableCMDEnv:           t.EnableDatasetImport,
		cms.DatasetAPIAuthTokenEnv: t.DatasetAPIAuthToken,
		cms.DatasetAPIURLEnv:       t.DatasetAPIURL,
		cms.ServiceAuthTokenEnv:    t.ServiceAuthToken,
	})
}

func errorAndExit(err error) {
	log.Event(nil, "unexpected error", log.Error(err))
	os.Exit(1)
}

package scripts

import (
	"bytes"
	"github.com/ONSdigital/dp-zebedee-utils/content/cms"
	"github.com/ONSdigital/dp-zebedee-utils/content/files"
	"github.com/ONSdigital/dp-zebedee-utils/content/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"text/template"
)

var (
	templateFile = "templates/run.cms.template.txt"
	output       = "generated"
	cmsRunFile   = output + "/run-cms.sh"
)

func GenerateCMSRunScript(builder *cms.Builder) (string, error) {
	if err := houseKeeping(); err != nil {
		return "", err
	}

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", errors.Wrap(err, "error loading template file")
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, builder.GetRunScriptArgs())
	if err != nil {
		log.Error.Fatal(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(cmsRunFile, buf.Bytes(), 0644)
	if err != nil {
		return "", errors.Wrap(err, "error generating run.sh from template")
	}
	return cmsRunFile, nil
}

func houseKeeping() error {
	exists, err := files.Exists(output)
	if err != nil {
		return err
	}

	if !exists {
		if err := os.Mkdir(output, 0700); err != nil {
			return errors.Wrap(err, "error creating generated dir")
		}
	}

	exists, err = files.Exists(cmsRunFile)
	if err != nil {
		return err
	}

	if exists {
		log.Info.Println("removing existing run-cms.sh file")
		if err := os.Remove(cmsRunFile); err != nil {
			return err
		}
	}

	return nil
}

package scripts

import (
	"bytes"
	"github.com/ONSdigital/dp-zebedee-utils/content/files"
	"github.com/ONSdigital/dp-zebedee-utils/content/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"text/template"
)

var (
	templateFile = "templates/run.cms.template.txt"
	cmsRunFile   = "generated/run-cms.sh"
)

func GenerateCMSRunScript(rootDir string) (string, error) {
	exists, err := files.Exists(cmsRunFile)
	if err != nil {
		return "", err
	}

	if exists {
		log.Info.Println("removing existing run-cms.sh file")
		err := os.Remove(cmsRunFile)
		if err != nil {
			return "", err
		}
	}

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", errors.Wrap(err, "error loading template file")
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{"ZebedeeRoot": rootDir})
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

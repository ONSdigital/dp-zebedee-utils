package scripts

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/ONSdigital/dp-zebedee-utils/content/cms"
	"github.com/ONSdigital/dp-zebedee-utils/content/files"
	"github.com/ONSdigital/log.go/log"
	"github.com/pkg/errors"
)

var (
	templateFile = "templates/run.cms.template.txt"
	output       = "generated"
	cmsRunFile   = output + "/run-cms.sh"
)

func GenerateCMSRunScript(t *cms.RunTemplate) (string, error) {
	if err := houseKeeping(); err != nil {
		return "", err
	}

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", errors.Wrap(err, "error loading template file")
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, t)
	if err != nil {
		log.Event(nil, "error generating run.sh from template", log.Error(err))
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
		log.Event(nil, "removing existing run-cms.sh file")
		if err := os.Remove(cmsRunFile); err != nil {
			return err
		}
	}

	return nil
}

func CopyToProjectDir(zebDir string, generatedFile string) (string, error) {
	target := path.Join(zebDir, "run-cms.sh")
	log.Event(nil, "copying run-cms.sh to zebedee project dir", log.Data{"target": target})
	b, err := ioutil.ReadFile(generatedFile)
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(target, b, 0700); err != nil {
		return "", err
	}
	return target, nil
}

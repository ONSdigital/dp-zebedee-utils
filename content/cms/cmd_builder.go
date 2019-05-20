package cms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ONSdigital/log.go/log"
	"github.com/pkg/errors"
)

// Build creates the Zebedee CMS directory structure
func (b *Builder) GenerateCMSContent() error {
	log.Event(nil, "generating CMS file structure and content", log.Data{
		"root":       b.zebedeeDir,
		"enable_cmd": b.enableCMD,
	})

	if err := b.createDirs(); err != nil {
		return err
	}

	err := b.copyContentZipToMaster()
	if err != nil {
		return err
	}

	err = b.unzipContentInMaster()
	if err != nil {
		return err
	}

	err = b.removeContentZipFromMaster()
	if err != nil {
		return err
	}

	err = b.createServiceAccount()
	if err != nil {
		return err
	}

	b.setDatasetAPIAuthToken()
	b.datasetAPIURL = "http://localhost:22000"
	return nil
}

func (b *Builder) createDirs() error {
	log.Event(nil, "creating zebedee directories")
	for _, dir := range b.dirs() {
		cmd := exec.Command("mkdir", dir)
		cmd.Stderr = b.Out
		cmd.Stdout = b.OutErr

		if err := cmd.Run(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("error while attempting to create zebedee directory: %s", dir))
		}
	}

	log.Event(nil, "successfully created zebedee directories", log.Data{
		"dirs": b.dirs(),
	})
	return nil
}

func (b *Builder) copyContentZipToMaster() error {
	log.Event(nil, "copying default content zip to master dir", log.Data{
		"master": b.masterDir,
	})
	cmd := newCommand("cp", "", defaultContentZip, b.masterDir)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "error executing command: copyContentZipToMaster")
	}
	return nil
}

func (b *Builder) unzipContentInMaster() error {
	log.Event(nil, "unzipping default content into master", log.Data{
		"master": b.masterDir,
	})
	cmd := newCommand("unzip", b.masterDir, "-q", defaultContentZip)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "error executing command: unzipContentInMaster")
	}
	return nil
}

func (b *Builder) removeContentZipFromMaster() error {
	log.Event(nil, "cleaning up default content zip")
	cmd := newCommand("rm", b.masterDir, defaultContentZip)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "error executing command: removeContentZipFromMaster")
	}
	return nil
}

func (b *Builder) createServiceAccount() error {
	serviceAuthToken, err := getServiceTokenID()
	if err != nil {
		return err
	}

	b.serviceAccountID = serviceAuthToken
	log.Event(nil, "generating CMD service account", log.Data{
		"serviceAccountID": b.serviceAccountID,
	})

	jsonB, err := json.Marshal(map[string]interface{}{"id": "Weyland-Yutani Corporation"})
	if err != nil {
		return errors.Wrap(err, "error marshaling service account JSON")
	}

	filename := filepath.Join(b.servicesDir, b.serviceAccountID+".json")
	err = ioutil.WriteFile(filename, jsonB, 0644)
	if err != nil {
		return errors.Wrap(err, "error writing service account JSON to file")
	}

	log.Event(nil, "successfully generated service account", log.Data{
		"serviceAccountID": b.serviceAccountID,
	})
	return nil
}

func getServiceTokenID() (string, error) {
	if serviceAuthToken := os.Getenv(ServiceAuthTokenEnv); serviceAuthToken != "" {
		log.Event(nil, fmt.Sprintf("found existing environment variable for %s using this id for generated service account", ServiceAuthTokenEnv))
		return serviceAuthToken, nil
	}

	log.Event(nil, fmt.Sprintf("no existing environment variable %s found, generating new ID for generated service account", ServiceAuthTokenEnv))

	return newRandomID(64), nil
}

func (b *Builder) setDatasetAPIAuthToken() {
	if datasetAPIAuthToken := os.Getenv(DatasetAPIAuthTokenEnv); datasetAPIAuthToken != "" {
		log.Event(nil, fmt.Sprintf("found existing environment variable for %s using this token value for generated run script", DatasetAPIAuthTokenEnv))
		b.datasetAPIAuthToken = datasetAPIAuthToken
	} else {
		log.Event(nil, fmt.Sprintf("no existing environment variable %s found generating new token for generated run script", DatasetAPIAuthTokenEnv))
		b.datasetAPIAuthToken = "FD0108EA-825D-411C-9B1D-41EF7727F465"
	}
}

func (b *Builder) setDatasetAPIURL() {
	if datasetAPIURL := os.Getenv(DatasetAPIURLEnv); datasetAPIURL != "" {
		log.Event(nil, fmt.Sprintf("found existing environment variable for %q using this value for generated run script", DatasetAPIURLEnv))
		b.datasetAPIURL = datasetAPIURL
	} else {
		log.Event(nil, fmt.Sprintf("no existing environment variable %s found generating new for generated run script", DatasetAPIURLEnv))
		b.datasetAPIURL = "http://localhost:22000"
	}
}

func (b *Builder) dirs() []string {
	return []string{
		b.zebedeeDir,
		b.masterDir,
		b.collectionsDir,
		b.publishLogDir,
		b.usersDir,
		b.sessionsDir,
		b.permissionsDir,
		b.teamsDir,
		b.launchPadDir,
		b.appKeysDir,
		b.servicesDir,
	}
}

func newCommand(name string, dir string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

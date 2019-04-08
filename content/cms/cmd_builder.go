package cms

import (
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/dp-zebedee-utils/content/files"
	"github.com/ONSdigital/dp-zebedee-utils/content/log"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const (
	Zebedee           = "zebedee"
	Master            = "master"
	Collections       = "collections"
	PublishLog        = "publish-log"
	Users             = "users"
	Sessions          = "sessions"
	Services          = "services"
	Permissions       = "permissions"
	Teams             = "teams"
	LaunchPad         = "launchpad"
	AppKeys           = "application-keys"
	defaultContentZip = "default-content.zip"

	datasetAPIAuthTokenEnv = "DATASET_API_AUTH_TOKEN"
	serviceAuthTokenEnv    = "SERVICE_AUTH_TOKEN"
	datasetAPIURLEnv       = "DATASET_API_URL"
)

var (
	Out    io.Writer
	OutErr io.Writer
)

type Builder struct {
	Out                 io.Writer
	OutErr              io.Writer
	rootDir             string
	zebedeeDir          string
	masterDir           string
	collectionsDir      string
	publishLogDir       string
	usersDir            string
	sessionsDir         string
	servicesDir         string
	permissionsDir      string
	teamsDir            string
	launchPadDir        string
	appKeysDir          string
	isCMD               bool
	serviceAccountID    string
	datasetAPIAuthToken string
	datasetAPIURL       string
}

// New construct a new cmd.Builder
func New(root string, isCMD bool) (*Builder, error) {
	zebedeeDir := filepath.Join(root, Zebedee)
	exists, err := files.Exists(zebedeeDir)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("cannot generate directory structure as a zebedee a dir already exists at the root location provided")
	}

	b := &Builder{
		rootDir:             root,
		zebedeeDir:          zebedeeDir,
		masterDir:           filepath.Join(zebedeeDir, Master),
		collectionsDir:      filepath.Join(zebedeeDir, Collections),
		publishLogDir:       filepath.Join(zebedeeDir, PublishLog),
		usersDir:            filepath.Join(zebedeeDir, Users),
		sessionsDir:         filepath.Join(zebedeeDir, Sessions),
		servicesDir:         filepath.Join(zebedeeDir, Services),
		permissionsDir:      filepath.Join(zebedeeDir, Permissions),
		teamsDir:            filepath.Join(zebedeeDir, Teams),
		launchPadDir:        filepath.Join(zebedeeDir, LaunchPad),
		appKeysDir:          filepath.Join(zebedeeDir, AppKeys),
		isCMD:               isCMD,
		datasetAPIURL:       "",
		datasetAPIAuthToken: "",
		serviceAccountID:    "",
	}
	return b, nil
}

// Build creates the Zebedee CMS directory structure
func (b *Builder) Build() error {
	log.Info.Printf("args: root=%s, cmd=%t\n", b.zebedeeDir, b.isCMD)
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

	if b.isCMD {
		err = b.createServiceAccount()
		if err != nil {
			return err
		}

		b.setDatasetAPIAuthToken()

		b.datasetAPIURL = "http://localhost:22000"

		log.Info.Printf("Your CMD config:\n\t%q: %s\n\t%q: %s\n\t%q%s", serviceAuthTokenEnv, b.serviceAccountID, datasetAPIAuthTokenEnv, b.datasetAPIAuthToken, datasetAPIURLEnv, b.datasetAPIURL)
	}

	return nil
}

func (b *Builder) GetRunScriptArgs() map[string]interface{} {
	args := make(map[string]interface{})
	args["ZEBEDEE_ROOT"] = b.rootDir
	args["ENABLE_DATASET_IMPORT"] = strconv.FormatBool(b.isCMD)
	args["DATASET_API_URL"] = b.datasetAPIURL
	args["DATASET_API_AUTH_TOKEN"] = b.datasetAPIAuthToken
	args["SERVICE_AUTH_TOKEN"] = b.serviceAccountID
	return args
}

func (b *Builder) createDirs() error {
	log.Info.Println("creating zebedee directories")
	for _, dir := range b.dirs() {
		cmd := exec.Command("mkdir", dir)
		cmd.Stderr = b.Out
		cmd.Stdout = b.OutErr

		if err := cmd.Run(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("error while attempting to create zebedee directory: %s", dir))
		}
		log.Info.Printf("created: %s\n", dir)
	}

	log.Info.Println("successfully created zebedee directories")
	return nil
}

func (b *Builder) copyContentZipToMaster() error {
	log.Info.Printf("copying default content zip to master: %s\n", b.masterDir)
	cmd := newCommand("cp", "", defaultContentZip, b.masterDir)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "error executing command: copyContentZipToMaster")
	}
	return nil
}

func (b *Builder) unzipContentInMaster() error {
	log.Info.Printf("unzipping default content into master: %s\n", b.masterDir)
	cmd := newCommand("unzip", b.masterDir, "-q", defaultContentZip)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "error executing command: unzipContentInMaster")
	}
	return nil
}

func (b *Builder) removeContentZipFromMaster() error {
	log.Info.Println("cleaning up default content zip")
	cmd := newCommand("rm", b.masterDir, defaultContentZip)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "error executing command: removeContentZipFromMaster")
	}
	return nil
}

func (b *Builder) createServiceAccount() error {
	if !b.isCMD {
		log.Info.Printf("skipping create service account step")
		return nil
	}

	serviceAuthToken, err := getServiceTokenID()
	if err != nil {
		return err
	}

	b.serviceAccountID = serviceAuthToken
	log.Info.Printf("generating CMD service account: ID: %s\n", b.serviceAccountID)

	jsonB, err := json.Marshal(map[string]interface{}{"id": "Wayne Enterprises"})
	if err != nil {
		return errors.Wrap(err, "error marshaling service account JSON")
	}

	filename := filepath.Join(b.servicesDir, b.serviceAccountID+".json")
	err = ioutil.WriteFile(filename, jsonB, 0644)
	if err != nil {
		return errors.Wrap(err, "error writing service account JSON to file")
	}

	log.Info.Printf("successfully generated service account make")
	return nil
}

func getServiceTokenID() (string, error) {
	if serviceAuthToken := os.Getenv(serviceAuthTokenEnv); serviceAuthToken != "" {
		log.Info.Printf("found existing environment variable for %q using this id for generated service account", serviceAuthTokenEnv)
		return serviceAuthToken, nil
	}

	log.Info.Printf("no existing environment variable %q found, generating new ID for generated service account", serviceAuthTokenEnv)
	id, err := uuid.NewV4()
	if err != nil {
		return "", errors.Wrap(err, "error generating UUID for service account")
	}

	return id.String(), nil
}

func (b *Builder) setDatasetAPIAuthToken() {
	if datasetAPIAuthToken := os.Getenv(datasetAPIAuthTokenEnv); datasetAPIAuthToken != "" {
		log.Info.Printf("found existing environment variable for %q using this token value for generated run script", datasetAPIAuthTokenEnv)
		b.datasetAPIAuthToken = datasetAPIAuthToken
	} else {
		log.Info.Printf("no existing environment variable %q found generating new token for generated run script", datasetAPIAuthTokenEnv)
		b.datasetAPIAuthToken = "FD0108EA-825D-411C-9B1D-41EF7727F465"
	}
}

func (b *Builder) setDatasetAPIURL() {
	if datasetAPIURL := os.Getenv(datasetAPIURLEnv); datasetAPIURL != "" {
		log.Info.Printf("found existing environment variable for %q using this value for generated run script", datasetAPIURLEnv)
		b.datasetAPIURL = datasetAPIURL
	} else {
		log.Info.Printf("no existing environment variable %q found generating new for generated run script", datasetAPIURLEnv)
		b.datasetAPIURL = "http://localhost:22000"
	}
}

func (b *Builder) dirs() []string {
	dirs := []string{
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
	}

	if b.isCMD {
		dirs = append(dirs, b.servicesDir)
	}

	return dirs
}

func newCommand(name string, dir string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = Out
	cmd.Stderr = OutErr
	return cmd
}

package cms

import (
	"errors"
	"io"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/ONSdigital/dp-zebedee-utils/content/files"
)

const (
	Zebedee                = "zebedee"
	Master                 = "master"
	Collections            = "collections"
	PublishLog             = "publish-log"
	Users                  = "users"
	Sessions               = "sessions"
	Services               = "services"
	Permissions            = "permissions"
	Teams                  = "teams"
	LaunchPad              = "launchpad"
	AppKeys                = "application-keys"
	defaultContentZip      = "default-content.zip"
	EnableCMDEnv           = "ENABLE_DATASET_IMPORT"
	DatasetAPIAuthTokenEnv = "DATASET_API_AUTH_TOKEN"
	ServiceAuthTokenEnv    = "SERVICE_AUTH_TOKEN"
	DatasetAPIURLEnv       = "DATASET_API_URL"
)

var (
	serviceIDChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	r              *rand.Rand
)

func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r = rand.New(s1)
}

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
	enableCMD           bool
	serviceAccountID    string
	datasetAPIAuthToken string
	datasetAPIURL       string
}

type RunTemplate struct {
	ZebedeeRoot         string
	EnableDatasetImport bool
	DatasetAPIURL       string
	DatasetAPIAuthToken string
	ServiceAuthToken    string
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
		enableCMD:           isCMD,
		datasetAPIURL:       "",
		datasetAPIAuthToken: "",
		serviceAccountID:    "",
	}
	return b, nil
}

func (b *Builder) GetRunTemplate() *RunTemplate {
	return &RunTemplate{
		ZebedeeRoot:         b.rootDir,
		EnableDatasetImport: b.enableCMD,
		DatasetAPIURL:       b.datasetAPIURL,
		DatasetAPIAuthToken: b.datasetAPIAuthToken,
		ServiceAuthToken:    b.serviceAccountID,
	}
}

func newRandomID(size int) string {
	id := ""
	for len(id) < size {
		id += string(serviceIDChars[r.Intn(len(serviceIDChars))])
	}

	return id
}

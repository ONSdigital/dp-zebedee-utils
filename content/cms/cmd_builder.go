package cms

import (
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/dp-zebedee-utils/content/log"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
)

var (
	Out    io.Writer
	OutErr io.Writer
)

type Builder struct {
	Out            io.Writer
	OutErr         io.Writer
	zebedeeDir     string
	masterDir      string
	collectionsDir string
	publishLogDir  string
	usersDir       string
	sessionsDir    string
	servicesDir    string
	permissionsDir string
	teamsDir       string
	launchPadDir   string
	appKeysDir     string
	serviceAccount bool
}

// New construct a new cmd.Builder
func New(root string, createServiceAccount bool) (*Builder, error) {
	zebedeeDir := filepath.Join(root, Zebedee)
	exists, err := exists(zebedeeDir)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("cannot generate directory structure as a zebedee a dir already exists at the root location provided")
	}

	c := &Builder{
		zebedeeDir:     zebedeeDir,
		masterDir:      filepath.Join(zebedeeDir, Master),
		collectionsDir: filepath.Join(zebedeeDir, Collections),
		publishLogDir:  filepath.Join(zebedeeDir, PublishLog),
		usersDir:       filepath.Join(zebedeeDir, Users),
		sessionsDir:    filepath.Join(zebedeeDir, Sessions),
		servicesDir:    filepath.Join(zebedeeDir, Services),
		permissionsDir: filepath.Join(zebedeeDir, Permissions),
		teamsDir:       filepath.Join(zebedeeDir, Teams),
		launchPadDir:   filepath.Join(zebedeeDir, LaunchPad),
		appKeysDir:     filepath.Join(zebedeeDir, AppKeys),
		serviceAccount: createServiceAccount,
	}
	return c, nil
}

// Build creates the Zebedee CMS directory structure
func (b *Builder) Build() error {
	log.Info.Printf("args: root=%s, cmd=%t\n", b.zebedeeDir, b.serviceAccount)
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

	if b.serviceAccount {
		err := b.createServiceAccount()
		if err != nil {
			return err
		}
	}
	return nil
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
	id, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "error generating UUID")
	}

	log.Info.Printf("generating CMD service account: ID: %s\n", id.String())

	jsonB, err := json.Marshal(map[string]interface{}{"id": "Wayne Enterprises"})
	if err != nil {
		return errors.Wrap(err, "error marshaling service account JSON")
	}

	filename := filepath.Join(b.servicesDir, id.String()+".json")
	err = ioutil.WriteFile(filename, jsonB, 0644)
	if err != nil {
		return errors.Wrap(err, "error writing service account JSON to file")
	}
	return nil
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

	if b.serviceAccount {
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

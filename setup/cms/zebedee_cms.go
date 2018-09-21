package cms

import (
	"github.com/ONSdigital/dp-zebedee-utils/setup/log"
	"github.com/pkg/errors"
	"io"
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
	Permissions       = "permissions"
	Teams             = "teams"
	LaunchPad         = "launchpad"
	AppKeys           = "application-keys"
	defaultContentZip = "DefaultContent.zip"
)

var (
	Out    io.Writer
	OutErr io.Writer
)

type CMS struct {
	Out            io.Writer
	OutErr         io.Writer
	zebedeeDir     string
	masterDir      string
	collectionsDir string
	publishLogDir  string
	usersDir       string
	sessionsDir    string
	permissionsDir string
	teamsDir       string
	launchPadDir   string
	appKeysDir     string
}

func New(root string) (*CMS, error) {
	zebedeeDir := filepath.Join(root, Zebedee)
	exists, err := exists(zebedeeDir)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("cannot generate directory structure as a zebedee a dir already exists at the root location provided")
	}

	c := &CMS{
		zebedeeDir:     zebedeeDir,
		masterDir:      filepath.Join(zebedeeDir, Master),
		collectionsDir: filepath.Join(zebedeeDir, Collections),
		publishLogDir:  filepath.Join(zebedeeDir, PublishLog),
		usersDir:       filepath.Join(zebedeeDir, Users),
		sessionsDir:    filepath.Join(zebedeeDir, Sessions),
		permissionsDir: filepath.Join(zebedeeDir, Permissions),
		teamsDir:       filepath.Join(zebedeeDir, Teams),
		launchPadDir:   filepath.Join(zebedeeDir, LaunchPad),
		appKeysDir:     filepath.Join(zebedeeDir, AppKeys),
	}
	return c, nil
}

func (c *CMS) Initialize() error {
	if err := c.createDirs(); err != nil {
		return err
	}

	err := c.copyContentZipToMaster()
	if err != nil {
		return err
	}

	err = c.unzipContentInMaster()
	if err != nil {
		return err
	}

	err = c.removeContentZipFromMaster()
	if err != nil {
		return err
	}

	return nil
}

func (c *CMS) createDirs() error {
	log.Info.Println("creating zebedee directories")
	for _, dir := range c.dirs() {
		cmd := exec.Command("mkdir", dir)
		cmd.Stderr = c.Out
		cmd.Stdout = c.OutErr

		if err := cmd.Run(); err != nil {
			return err
		}
		log.Info.Printf("created: %s\n", dir)
	}

	log.Info.Println("successfully created zebedee directories")
	return nil
}

func (c *CMS) copyContentZipToMaster() error {
	log.Info.Printf("copying default content zip to master: %s\n", c.masterDir)
	cmd := newCommand("cp", "", defaultContentZip, c.masterDir)
	return cmd.Run()
}

func (c *CMS) unzipContentInMaster() error {
	log.Info.Printf("unzipping default content into master: %s\n", c.masterDir)
	cmd := newCommand("unzip", c.masterDir, "-q", defaultContentZip)
	return cmd.Run()
}

func (c *CMS) removeContentZipFromMaster() error {
	log.Info.Println("cleaning up default content zip")
	cmd := newCommand("rm", c.masterDir, defaultContentZip)
	return cmd.Run()
}

func (c *CMS) dirs() []string {
	return []string{
		c.zebedeeDir,
		c.masterDir,
		c.collectionsDir,
		c.publishLogDir,
		c.usersDir,
		c.sessionsDir,
		c.permissionsDir,
		c.teamsDir,
		c.launchPadDir,
		c.appKeysDir,
	}
}

func (c *CMS) GetZebedeeRoot() string {
	return c.zebedeeDir
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

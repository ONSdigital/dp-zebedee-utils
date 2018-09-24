package cms

import (
	"github.com/ONSdigital/dp-zebedee-utils/content/log"
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
	permissionsDir string
	teamsDir       string
	launchPadDir   string
	appKeysDir     string
}

func New(root string) (*Builder, error) {
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
		permissionsDir: filepath.Join(zebedeeDir, Permissions),
		teamsDir:       filepath.Join(zebedeeDir, Teams),
		launchPadDir:   filepath.Join(zebedeeDir, LaunchPad),
		appKeysDir:     filepath.Join(zebedeeDir, AppKeys),
	}
	return c, nil
}

func (b *Builder) Build() error {
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

	return nil
}

func (b *Builder) createDirs() error {
	log.Info.Println("creating zebedee directories")
	for _, dir := range b.dirs() {
		cmd := exec.Command("mkdir", dir)
		cmd.Stderr = b.Out
		cmd.Stdout = b.OutErr

		if err := cmd.Run(); err != nil {
			return err
		}
		log.Info.Printf("created: %s\n", dir)
	}

	log.Info.Println("successfully created zebedee directories")
	return nil
}

func (b *Builder) copyContentZipToMaster() error {
	log.Info.Printf("copying default content zip to master: %s\n", b.masterDir)
	cmd := newCommand("cp", "", defaultContentZip, b.masterDir)
	return cmd.Run()
}

func (b *Builder) unzipContentInMaster() error {
	log.Info.Printf("unzipping default content into master: %s\n", b.masterDir)
	cmd := newCommand("unzip", b.masterDir, "-q", defaultContentZip)
	return cmd.Run()
}

func (b *Builder) removeContentZipFromMaster() error {
	log.Info.Println("cleaning up default content zip")
	cmd := newCommand("rm", b.masterDir, defaultContentZip)
	return cmd.Run()
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
	}
}

func (b *Builder) GetZebedeeRoot() string {
	return b.zebedeeDir
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

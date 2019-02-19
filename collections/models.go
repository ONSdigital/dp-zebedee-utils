package collections

import (
	"fmt"
	"github.com/ONSdigital/log.go/log"
	"github.com/satori/go.uuid"
	"path"
)

type Metadata struct {
	Name           string
	CollectionRoot string
	InProgress     string
	Complete       string
	Reviewed       string
	CollectionJSON string
}

type Collection struct {
	Metadata              *Metadata     `json:"-"`
	ApprovalStatus        string        `json:"approvalStatus"`
	PublishComplete       bool          `json:"publishComplete"`
	IsEncrypted           bool          `json:"isEncrypted"`
	CollectionOwner       string        `json:"collectionOwner"`
	TimeSeriesImportFiles []string      `json:"timeseriesImportFiles"`
	ID                    string        `json:"id"`
	Name                  string        `json:"name"`
	Type                  string        `json:"type"`
	Teams                 []interface{} `json:"teams"`
}

func New(rootPath string, name string) *Collection {
	newID, _ := uuid.NewV4()
	id := fmt.Sprintf("%s-%s", name, newID.String())
	metadata := NewMetadata(rootPath, name)

	return &Collection{
		Metadata:        metadata,
		ApprovalStatus:  "NOT_STARTED",
		CollectionOwner: "PUBLISHING_SUPPORT",
		IsEncrypted:     false,
		PublishComplete: false,
		Type:            "manual",
		ID:              id,
		Name:            name,
		TimeSeriesImportFiles: []string{},
		Teams: []interface{}{},
	}
}

func NewMetadata(collectionsDirPath string, name string) *Metadata {
	collectionRoot := path.Join(collectionsDirPath, name)

	return &Metadata{
		CollectionRoot: collectionRoot,
		Name:           name,
		CollectionJSON: path.Join(collectionsDirPath, name+".json"),
		InProgress:     path.Join(collectionRoot, "inprogress"),
		Complete:       path.Join(collectionRoot, "complete"),
		Reviewed:       path.Join(collectionRoot, "reviewed"),
	}
}

func (c *Collection) GetName() string {
	return c.Name
}

func (c *Collection) GetRootPath() string {
	return c.Metadata.CollectionRoot
}

func (c *Collection) GetInProgress() string {
	return c.Metadata.InProgress
}

func (c *Collection) GetComplete() string {
	return c.Metadata.Complete
}

func (c *Collection) GetReviewed() string {
	return c.Metadata.Reviewed
}

func (c *Collection) Contains(uri string) bool {
	if Exists(path.Join(c.Metadata.InProgress, uri)) {
		log.Event(nil, "collection contains uri", log.Data{
			"uri":        uri,
			"collection": c.Name,
			"dir":        c.Metadata.InProgress,
		})
		return true
	}

	if Exists(path.Join(c.Metadata.Complete, uri)) {
		log.Event(nil, "collection contains uri", log.Data{
			"uri":        uri,
			"collection": c.Name,
			"dir":        c.Metadata.Complete,
		})
		return true
	}

	if Exists(path.Join(c.Metadata.Reviewed, uri)) {
		log.Event(nil, "collection contains uri", log.Data{
			"uri":        uri,
			"collection": c.Name,
			"dir":        c.Metadata.Reviewed,
		})
		return true
	}
	return false
}

func (c *Collection) getDirs() []string {
	return []string{
		c.Metadata.CollectionRoot,
		c.Metadata.InProgress,
		c.Metadata.Complete,
		c.Metadata.Reviewed,
	}
}

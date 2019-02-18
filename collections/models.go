package collections

import (
	"fmt"
	"github.com/satori/go.uuid"
	"path"
)

type CollectionMetadata struct {
	name           string
	collectionRoot string
	inProgress     string
	complete       string
	reviewed       string
	collectionJSON string
}

type Collection struct {
	Metadata              *CollectionMetadata `json:"-"`
	ApprovalStatus        string              `json:"approvalStatus"`
	PublishComplete       bool                `json:"publishComplete"`
	IsEncrypted           bool                `json:"isEncrypted"`
	CollectionOwner       string              `json:"collectionOwner"`
	TimeSeriesImportFiles []string            `json:"timeseriesImportFiles"`
	ID                    string              `json:"id"`
	Name                  string              `json:"name"`
	Type                  string              `json:"type"`
	Teams                 []interface{}       `json:"teams"`
}

func New(rootPath string, name string) *Collection {
	collectionRoot := path.Join(rootPath, name)

	newID, _ := uuid.NewV4()
	id := fmt.Sprintf("%s-%s", name, newID.String())

	metadata := &CollectionMetadata{
		collectionRoot: collectionRoot,
		name:           name,
		collectionJSON: path.Join(rootPath, name+".json"),
		inProgress:     path.Join(collectionRoot, "inprogress"),
		complete:       path.Join(collectionRoot, "complete"),
		reviewed:       path.Join(collectionRoot, "reviewed"),
	}

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

func (c *Collection) GetName() string {
	return c.Name
}

func (c *Collection) GetRootPath() string {
	return c.Metadata.collectionRoot
}

func (c *Collection) GetInProgress() string {
	return c.Metadata.inProgress
}

func (c *Collection) GetComplete() string {
	return c.Metadata.complete
}

func (c *Collection) GetReviewed() string {
	return c.Metadata.reviewed
}

func (c *Collection) getDirs() []string {
	return []string{
		c.Metadata.collectionRoot,
		c.Metadata.inProgress,
		c.Metadata.complete,
		c.Metadata.reviewed,
	}
}

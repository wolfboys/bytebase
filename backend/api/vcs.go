package api

import (
	"context"
	"encoding/json"
)

type VCSType string

const (
	GITLAB_SELF_HOST VCSType = "GITLAB_SELF_HOST"
)

func (e VCSType) String() string {
	switch e {
	case GITLAB_SELF_HOST:
		return "GITLAB_SELF_HOST"
	}
	return "UNKNOWN"
}

type VCS struct {
	ID int `jsonapi:"primary,vcs"`

	// Standard fields
	CreatorId int
	Creator   *Principal `jsonapi:"attr,creator"`
	CreatedTs int64      `jsonapi:"attr,createdTs"`
	UpdaterId int
	Updater   *Principal `jsonapi:"attr,updater"`
	UpdatedTs int64      `jsonapi:"attr,updatedTs"`

	// Domain specific fields
	Name string `jsonapi:"attr,name"`
	// Identify oauth callback endpoint
	Uuid          string  `jsonapi:"attr,uuid"`
	Type          VCSType `jsonapi:"attr,type"`
	InstanceURL   string  `jsonapi:"attr,instanceURL"`
	ApiURL        string  `jsonapi:"attr,apiURL"`
	ApplicationId string  `jsonapi:"attr,applicationId"`
	Secret        string  `jsonapi:"attr,secret"`
	AccessToken   string  `jsonapi:"attr,accessToken"`
}

type VCSCreate struct {
	// Standard fields
	// Value is assigned from the jwt subject field passed by the client.
	CreatorId int

	// Domain specific fields
	Name          string  `jsonapi:"attr,name"`
	Type          VCSType `jsonapi:"attr,type"`
	InstanceURL   string  `jsonapi:"attr,instanceURL"`
	ApiURL        string
	ApplicationId string `jsonapi:"attr,applicationId"`
	Secret        string `jsonapi:"attr,secret"`
}

type VCSFind struct {
	ID *int

	// Domain specific fields
	Uuid *string
}

func (find *VCSFind) String() string {
	str, err := json.Marshal(*find)
	if err != nil {
		return err.Error()
	}
	return string(str)
}

type VCSPatch struct {
	ID int

	// Standard fields
	// Value is assigned from the jwt subject field passed by the client.
	UpdaterId int

	// Domain specific fields
	AccessToken *string `jsonapi:"attr,accessToken"`
}

type VCSDelete struct {
	ID int

	// Standard fields
	// Value is assigned from the jwt subject field passed by the client.
	DeleterId int
}

type VCSService interface {
	CreateVCS(ctx context.Context, create *VCSCreate) (*VCS, error)
	FindVCSList(ctx context.Context, find *VCSFind) ([]*VCS, error)
	FindVCS(ctx context.Context, find *VCSFind) (*VCS, error)
	PatchVCS(ctx context.Context, patch *VCSPatch) (*VCS, error)
	DeleteVCS(ctx context.Context, delete *VCSDelete) error
}

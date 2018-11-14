package model

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/materials-commons/mc/pkg/mc"
	"github.com/pkg/errors"
)

type GlobusUploadSchema struct {
	ModelSimpleNoID
	ID               string `db:"id" json:"id" r:"id"`
	Owner            string `db:"owner" json:"owner" r:"owner"`
	Path             string `db:"path" json:"path" r:"path"`
	ProjectID        string `db:"project_id" json:"project_id" r:"project_id"`
	GlobusAclID      string `db:"globus_acl_id" json:"globus_acl_id" r:"globus_acl_id"`
	GlobusEndpointID string `db:"globus_endpoint_id" json:"globus_endpoint_id" r:"globus_endpoint_id"`
	GlobusIdentityID string `db:"globus_identity_id" json:"globus_identity_id" r:"globus_identity_id"`
}

type AddGlobusUploadModel struct {
	ID               string
	Owner            string
	Path             string
	ProjectID        string
	GlobusAclID      string
	GlobusEndpointID string
	GlobusIdentityID string
}

func (m AddGlobusUploadModel) Validate() error {
	err := validation.ValidateStruct(&m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.Owner, validation.Required, validation.By(IsEmail)),
		validation.Field(&m.Path, validation.Required),
		validation.Field(&m.ProjectID, validation.Required, is.UUIDv4),
		validation.Field(&m.GlobusAclID, validation.Required),
		validation.Field(&m.GlobusEndpointID, validation.Required),
		validation.Field(&m.GlobusIdentityID, validation.Required))

	if err != nil {
		return errors.WithMessage(mc.ErrValidation, err.Error())
	}

	return nil
}

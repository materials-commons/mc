package model

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/materials-commons/mc/pkg/mc"
	"github.com/pkg/errors"
)

type FileLoadSchema struct {
	ModelSimple
	ProjectID string   `db:"project_id" json:"project_id" r:"project_id"`
	Path      string   `db:"path" json:"path" r:"path"`
	Owner     string   `db:"owner" json:"owner" r:"path"`
	Exclude   []string `json:"exclude" r:"exclude"`
	Loading   bool     `db:"loading" json:"loading" r:"loading"`
}

type AddFileLoadModel struct {
	ProjectID string
	Path      string
	Owner     string
	Exclude   []string
}

func (m AddFileLoadModel) Validate() error {
	err := validation.ValidateStruct(&m,
		validation.Field(&m.ProjectID, validation.Required),
		validation.Field(&m.Path, validation.Required),
		validation.Field(&m.Owner, validation.Required, validation.By(IsEmail)))
	if err != nil {
		return errors.WithMessage(mc.ErrValidation, err.Error())
	}

	return nil
}

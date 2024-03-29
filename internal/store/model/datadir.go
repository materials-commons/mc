package model

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/materials-commons/mc/pkg/mc"
	"github.com/pkg/errors"
)

type DatadirSchema struct {
	Model
	Parent   string `db:"parent" json:"parent" r:"parent"`
	Project  string `db:"project" json:"project" r:"project"`
	Shortcut bool   `db:"shortcut" json:"shortcut" r:"shortcut"`
}

type DatadirEntryModel struct {
	DatadirSimpleModel
	Directories []DatadirSimpleModel  `json:"directories" r:"directories"`
	Files       []DatafileSimpleModel `json:"files" r:"files"`
}

type DatadirSimpleModel struct {
	ID   string `db:"id" json:"id" r:"id"`
	Name string `db:"name" json:"name" r:"name"`
}

type AddDatadirModel struct {
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	Parent    string
	ProjectID string
}

func (d AddDatadirModel) Validate() error {
	err := validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required, validation.Length(1, 250)),
		validation.Field(&d.Owner, validation.Required),
		validation.Field(&d.ProjectID, validation.Required),
	)

	if err != nil {
		return errors.WithMessage(mc.ErrValidation, err.Error())
	}

	return nil
}

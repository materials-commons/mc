package model

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/materials-commons/mc/pkg/mc"
	"github.com/pkg/errors"
)

type DatafileSchema struct {
	Model
	Checksum    string            `db:"checksum" json:"checksum" r:"checksum"`
	Current     bool              `db:"current" json:"current" r:"current"`
	Description string            `db:"description" json:"description" r:"description"`
	MediaType   DatafileMediaType `json:"mediatype" r:"mediatype"`
	Parent      string            `db:"parent" json:"parent" r:"parent"`
	Size        int64             `db:"size" json:"size" r:"size"`
	Uploaded    int               `db:"uploaded" json:"uploaded" r:"uploaded"`
	UsesID      string            `db:"usesid" json:"usesid" r:"usesid"`
}

type DatafileMediaType struct {
	Description string `db:"description" json:"description" r:"description"`
	Mime        string `db:"mime" json:"mime" r:"mime"`
}

type DatafileSimpleModel struct {
	ID   string `json:"id" r:"id"`
	Name string `json:"name" r:"name"`
}

type AddDatafileModel struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Checksum    string
	Description string `json:"description"`
	UsesID      string
	Parent      string
	Size        int64
	ProjectID   string
	DatadirID   string
	MediaType   DatafileMediaType
}

func (d AddDatafileModel) Validate() error {
	err := validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&d.Owner, validation.Required, validation.By(IsEmail)),
		validation.Field(&d.Description, validation.Length(0, 300)))

	if err != nil {
		return errors.WithMessage(mc.ErrValidation, err.Error())
	}

	return nil
}

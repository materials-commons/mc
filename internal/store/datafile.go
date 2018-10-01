package store

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
)

type DatafileSchema struct {
	Model
	Checksum    string            `db:"checksum" json:"checksum" r:"checksum"`
	Current     bool              `db:"current" json:"current" r:"current"`
	Description string            `db:"description" json:"description" r:"description"`
	MediaType   DatafileMediaType `json:"mediatype" r:"mediatype"`
	Parent      string            `db:"parent" json:"parent" r:"parent"`
	Size        int               `db:"size" json:"size" r:"size"`
	Uploaded    int               `db:"uploaded" json:"uploaded" r:"uploaded"`
	UsesID      string            `db:"usesid" json:"usesid" r:"usesid"`
}

type DatafileMediaType struct {
	Description string `db:"description" json:"description" r:"description"`
	Mime        string `db:"mime" json:"mime" r:"mime"`
}

type AddDatafileModel struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Checksum    string
	Description string `json:"description"`
	UsesID      string
	Parent      string
	Size        int
	ProjectID   string
	DatadirID   string
}

func (d AddDatafileModel) Validate() error {
	err := validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&d.Owner, validation.Required, is.Email),
		validation.Field(&d.Description, validation.Required, validation.Length(0, 300)),
		validation.Field(&d.Size, validation.Required, validation.Min(1)))

	if err != nil {
		return errors.WithMessage(ErrValidation, err.Error())
	}

	return nil
}

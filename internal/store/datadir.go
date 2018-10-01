package store

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

type DatadirSchema struct {
	Model
	Parent   string `db:"parent" json:"parent" r:"parent"`
	Project  string `db:"project" json:"project" r:"project"`
	Shortcut bool   `db:"shortcut" json:"shortcut" r:"shortcut"`
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
		return errors.WithMessage(ErrValidation, err.Error())
	}

	return nil
}

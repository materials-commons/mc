package store

import (
	"time"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
)

type UserSchema struct {
	ModelSimple
	Admin           bool      `db:"admin" json:"admin" r:"admin"`
	APIKey          string    `db:"apikey" json:"apikey" r:"apikey"`
	BetaUser        bool      `db:"beta_user" json:"beta_user" r:"beta_user"`
	DemoInstalled   bool      `db:"demo_installed" json:"demo_installed" r:"demo_installed"`
	Email           string    `db:"email" json:"email" r:"email"`
	Fullname        string    `db:"fullname" json:"fullname" r:"fullname"`
	Password        string    `db:"password" json:"-" r:"password"`
	IsTemplateAdmin bool      `db:"is_template_admin" json:"is_template_admin" r:"is_template_admin"`
	LastLogin       time.Time `db:"last_login" json:"last_login" r:"last_login"`
}

type AddUserModel struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
}

func (u AddUserModel) Validate() error {
	err := validation.ValidateStruct(&u,
		validation.Field(&u.Fullname, validation.Required, validation.Length(1, 40)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(1, 100)))
	if err != nil {
		return errors.WithMessage(ErrValidation, err.Error())
	}

	return nil
}

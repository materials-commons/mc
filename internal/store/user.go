package store

import "time"

type UserModel struct {
	Model
	Admin           bool      `db:"admin" json:"admin"`
	APIKey          string    `db:"apikey" json:"apikey"`
	BetaUser        bool      `db:"beta_user" json:"beta_user"`
	DemoInstalled   bool      `db:"demo_installed" json:"demo_installed"`
	EMail           string    `db:"email" json:"email"`
	Fullname        string    `db:"fullname" json:"fullname"`
	Name            string    `db:"name" json:"name"`
	OType           string    `db:"otype" json:"otype"`
	Password        string    `db:"-" json:"-"`
	IsTemplateAdmin bool      `db:"is_template_admin" json:"is_template_admin"`
	LastLogin       time.Time `db:"last_login" json:"last_login"`
}

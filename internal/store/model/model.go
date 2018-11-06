package model

import "time"

type Model struct {
	ID        string    `db:"id" json:"id" r:"id,omitempty"`
	Birthtime time.Time `db:"birthime" json:"birthtime" r:"birthtime"`
	MTime     time.Time `db:"mtime" json:"mtime" r:"mtime"`
	Name      string    `db:"name" json:"name" r:"name"`
	Owner     string    `db:"owner" json:"owner" r:"owner"`
	OType     string    `db:"otype" json:"otype" r:"otype"`
}

type ModelSimple struct {
	ID        string    `db:"id" json:"id" r:"id,omitempty"`
	Birthtime time.Time `db:"birthtime" json:"birthtime" r:"birthtime"`
	MTime     time.Time `db:"mtime" json:"mtime" r:"mtime"`
	OType     string    `db:"otype" json:"otype" r:"otype"`
}

type ModelSimpleNoID struct {
	Birthtime time.Time `db:"birthtime" json:"birthtime" r:"birthtime"`
	MTime     time.Time `db:"mtime" json:"mtime" r:"mtime"`
	OType     string    `db:"otype" json:"otype" r:"otype"`
}

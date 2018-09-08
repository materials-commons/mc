package store

import "time"

type Model struct {
	ID        string    `db:"id" json:"id"`
	Birthtime time.Time `db:"birthtime" json:"birthtime"`
	MTime     time.Time `db:"mtime" json:"mtime"`
}

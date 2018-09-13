package store

import "time"

type Model struct {
	ID        string    `db:"id" json:"id" r:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at" r:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" r:"updated_at"`
	Name      string    `db:"name" json:"name" r:"name"`
	Owner     string    `db:"owner" json:"owner" r:"owner"`
	OType     string    `db:"otype" json:"otype" r:"otype"`
}

type ModelSimple struct {
	ID        string    `db:"id" json:"id" r:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at" r:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" r:"updated_at"`
	OType     string    `db:"otype" json:"otype" r:"otype"`
}

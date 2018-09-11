package store

import "time"

type SampleSchema struct {
	ID          string    `db:"id" json:"id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	OType       string    `db:"otype" json:"otype"`
	Owner       string    `db:"owner" json:"owner"`
}

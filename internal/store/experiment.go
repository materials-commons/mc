package store

import "time"

type ExperimentSchema struct {
	ID          string    `db:"id" json:"id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Owner       string    `db:"owner" json:"owner"`
	OType       string    `db:"otype" json:"otype"`
	Status      string    `db:"status" json:"status"`
}

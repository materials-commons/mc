package model

type SampleSchema struct {
	Model
	Description string `db:"description" json:"description"`
}

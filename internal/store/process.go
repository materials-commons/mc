package store

import "time"

type ProcessSchema struct {
	ID            string    `db:"id" json:"id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	OType         string    `db:"otype" json:"otype"`
	Name          string    `db:"name" json:"name"`
	DoesTransform bool      `db:"does_transform" json:"does_transform"`
	Owner         string    `db:"owner" json:"owner"`
	Note          string    `db:"note" json:"note"`
	ProcessType   string    `db:"process_type" json:"process_type"`
	TemplateID    string    `db:"template_id" json:"template_id"`
	TemplateName  string    `db:"template_name" json:"template_name"`
}

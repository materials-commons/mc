package model

import "time"

type ProjectAccessSchema struct {
	ID        string    `json:"id" r:"id"`
	UserID    string    `json:"user_id" r:"user_id"`
	Birthtime time.Time `json:"birthtime" r:"birthtime"`
	ProjectID string    `json:"project_id" r:"project_id"`
}

type ProjectUserAccessModel struct {
	ProjectAccessSchema
	Fullname string `json:"fullname" r:"fullname"`
}

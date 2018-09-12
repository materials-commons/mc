package store

import "time"

type ProjectSchema struct {
	ID          string    `db:"id" json:"id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Owner       string    `db:"owner" json:"owner"`
}

type ProjectExtendedModel struct {
	ProjectSchema
	FilesCount    int                  `json:"files_count"`
	Users         []ProjectUserModel   `json:"users"`
	Samples       []SampleSchema       `json:"samples"`
	Processes     []ProcessSchema      `json:"processes"`
	Experiments   []ExperimentSchema   `json:"experiments"`
	Relationships ProjectRelationships `json:"relationships"`
}

type ProjectRelationships struct {
	Process2Sample    []Process2Sample    `json:"process2sample"`
	Experiment2Sample []Experiment2Sample `json:"experiment2sample"`
}

type ProjectUserModel struct {
	ID        string    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	BetaUser  bool      `db:"beta_user" json:"beta_user"`
	Fullname  string    `db:"fullname" json:"fullname"`
	ProjectID string    `db:"project_id" json:"project_id"`
	UserID    string    `db:"user_id" json:"user_id"`
}

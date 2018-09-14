package store

type AccessSchema struct {
	ModelSimple
	Permissions string `db:"permissions" json:"permissions" r:"permissions"`
	ProjectID   string `db:"project_id" json:"project_id" r:"project_id"`
	ProjectName string `db:"project_name" json:"project_name" r:"project_name"` // TODO: Is this needed?
	UserID      string `db:"user_id" json:"user_id" r:"user_id"`
}

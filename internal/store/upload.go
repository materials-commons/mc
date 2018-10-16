package store

type UploadSchema struct {
	ModelSimple
	ProjectID string   `db:"project_id" json:"project_id" r:"project_id"`
	Path      string   `db:"path" json:"path" r:"path"`
	Owner     string   `db:"owner" json:"owner" r:"path"`
	Exclude   []string `json:"exclude" r:"exclude"`
}

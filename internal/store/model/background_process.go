package model

type BackgroundProcessSchema struct {
	ModelSimple
	UserID                string `db:"user_id" json:"user_id" r:"user_id"`
	ProjectID             string `db:"project_id" json:"project_id" r:"project_id"`
	BackgroundProcessID   string `db:"background_process_id" json:"background_process_id" r:"background_process_id"`
	BackgroundProcessType string `db:"background_process_type" json:"background_process_type" r:"background_process_type"`
	Status                string `db:"status" json:"status" r:"status"`
	Message               string `db:"message" json:"message" r:"message"`
	IsFinished            bool   `db:"is_finished" json:"is_finished" r:"is_finished"`
	IsOk                  bool   `db:"is_ok" json:"is_ok" r:"is_ok"`
}

type AddBackgroundProcessModel struct {
	UserID                string `db:"user_id" json:"user_id" r:"user_id"`
	ProjectID             string `db:"project_id" json:"project_id" r:"project_id"`
	BackgroundProcessID   string `db:"background_process_id" json:"background_process_id" r:"background_process_id"`
	BackgroundProcessType string `db:"background_process_type" json:"background_process_type" r:"background_process_type"`
	Status                string `db:"status" json:"status" r:"status"`
	Message               string `db:"message" json:"message" r:"message"`
}

type UpdateBackgroundProcessModel struct {
	Status  string `db:"status" json:"status" r:"status"`
	Message string `db:"message" json:"message" r:"message"`
}

type MarkOKBackgroundProcessModel struct {
	IsOk bool `db:"is_ok" json:"is_ok" r:"is_ok"`
}

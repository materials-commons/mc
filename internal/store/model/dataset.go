package model

type DatasetSchema struct {
	ModelSimple
	SelectionID   string        `json:"selection_id" r:"selection_id"`
	Title         string        `json:"title" r:"title"`
	Owner         string        `json:"owner" r:"owner"`
	DOI           string        `json:"doi" r:"doi"`
	Published     bool          `json:"published" r:"published"`
	FileSelection FileSelection `json:"file_selection" r:"file_selection"`
}

type FileSelection struct {
	ID           string   `json:"id" r:"id"`
	IncludeFiles []string `json:"include_files" r:"include_files"`
	ExcludeFiles []string `json:"exclude_files" r:"exclude_files"`
	IncludeDirs  []string `json:"include_dirs" r:"include_dirs"`
	ExcludeDirs  []string `json:"exclude_dirs" r:"exclude_dirs"`
}

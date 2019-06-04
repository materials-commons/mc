package model

type DatasetSchema struct {
	ModelSimple
	Title     string `json:"title"`
	Owner     string `json:"owner"`
	DOI       string `json:"doi"`
	Published bool   `json:"published"`
}

type FileSelection struct {
	ID            string   `json:"id"`
	IncludedFiles []string `json:"included_files"`
	ExcludedFiles []string `json:"excluded_files"`
	IncludedDirs  []string `json:"included_dirs"`
	ExcludedDirs  []string `json:"excluded_dirs"`
}

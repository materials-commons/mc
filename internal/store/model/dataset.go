package model

type DatasetSchema struct {
	ModelSimple
	Title     string `json:"title"`
	Owner     string `json:"owner"`
	DOI       string `json:"doi"`
	Published bool   `json:"published"`
}

type FileSelection struct {
	ID           string   `json:"id"`
	IncludeFiles []string `json:"include_files"`
	ExcludeFiles []string `json:"exclude_files"`
	IncludeDirs  []string `json:"include_dirs"`
	ExcludeDirs  []string `json:"exclude_dirs"`
}

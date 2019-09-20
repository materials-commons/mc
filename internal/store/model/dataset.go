package model

type DatasetSchema struct {
	ModelSimple
	SelectionID   string            `json:"selection_id" r:"selection_id"`
	Title         string            `json:"title" r:"title"`
	Owner         string            `json:"owner" r:"owner"`
	DOI           string            `json:"doi" r:"doi"`
	Published     bool              `json:"published" r:"published"`
	FileSelection FileSelection     `json:"file_selection" r:"file_selection"`
	Zip           DatasetZipDetails `json:"zip" r:"zip""`
}

type DatasetZipDetails struct {
	Size     int64  `json:"size" r:"size"`
	Filename string `json:"filename" r:"filename"`
}

type FileSelection struct {
	ID           string   `json:"id" r:"id"`
	IncludeFiles []string `json:"include_files" r:"include_files"`
	ExcludeFiles []string `json:"exclude_files" r:"exclude_files"`
	IncludeDirs  []string `json:"include_dirs" r:"include_dirs"`
	ExcludeDirs  []string `json:"exclude_dirs" r:"exclude_dirs"`
}

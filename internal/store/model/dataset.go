package model

type DatasetSchema struct {
	ModelSimple
	Title     string `json:"title"`
	Owner     string `json:"owner"`
	DOI       string `json:"doi"`
	Published bool   `json:"published"`
}

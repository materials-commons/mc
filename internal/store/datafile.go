package store

type DatafileSchema struct {
	Model
	Checksum    string            `db:"checksum" json:"checksum" r:"checksum"`
	Current     bool              `db:"current" json:"current" r:"current"`
	Description string            `db:"description" json:"description" r:"description"`
	MediaType   DatafileMediaType `json:"mediatype" r:"mediatype"`
	Parent      string            `db:"parent" json:"parent" r:"parent"`
	Size        int               `db:"size" json:"size" r:"size"`
	Uploaded    int               `db:"uploaded" json:"uploaded" r:"uploaded"`
	UsesID      string            `db:"usesid" json:"usesid" r:"usesid"`
}

type DatafileMediaType struct {
	Description string `db:"description" json:"description" r:"description"`
	Mime        string `db:"mime" json:"mime" r:"mime"`
}

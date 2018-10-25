package model

import "time"

type PublishedDatasetSchema struct {
	ModelSimple
	Authors       []DatasetAuthor       `json:"authors"`
	Description   string                `json:"description"`
	DOI           string                `json:"doi"`
	Funding       string                `json:"funding"`
	Institution   string                `json:"institution"`
	Keywords      []string              `json:"keywords"`
	License       License               `json:"license"`
	Papers        []DatasetPaper        `json:"papers"`
	Published     bool                  `json:"published"`
	PublishedDate time.Time             `json:"published_date"`
	Title         string                `json:"title"`
	Owner         string                `json:"owner"`
	FileCount     int                   `json:"file_count"`
	Stats         PublishedDatasetStats `json:"stats"`
	Zip           DatasetZip            `json:"zip"`
}

type License struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type PublishedDatasetStats struct {
	UniqueViewCount PublishedDatasetUniqueViewCount `json:"unique_view_count"`
	CommentCount    int                             `json:"comment_count"`
	InterestedUsers int                             `json:"interested_users"`
}

type PublishedDatasetUniqueViewCount struct {
	Total int `json:"total"`
	// also, eventually, 'anonymous': items with user_ids that are not users
	//   and 'authenticated': items with user_ids that are users
	// Add these items at a later point
}

type DatasetZip struct {
	Filename string `json:"filename"`
	Size     int    `json:"size"`
}

type DatasetAuthor struct {
	Affiliation string `json:"affiliation"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
}

type DatasetPaper struct {
	Abstract string `json:"abstract"`
	Authors  string `json:"authors"`
	DOI      string `json:"doi"`
	Link     string `json:"link"`
	Title    string `json:"title"`
}

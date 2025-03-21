package models

// Song represents song data model
type Song struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group"`
	SongName    string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// SongRequest represents song creation request
type SongRequest struct {
	GroupName string `json:"group"`
	SongName  string `json:"song"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// SongText represents lyrics response
type SongText struct {
	Text string `json:"text"`
}

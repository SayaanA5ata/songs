package song

type SongCreateRequest struct {
	Group string `json:"group" validate:"required"`
	Name  string `json:"song" validate:"required"`
	Date  string `json:"releaseDate" validate:"required"`
	Text  string `json:"text" validate:"required"`
	Link  string `json:"link" validate:"required,http_url"`
	Hash  string `json:"hash" gorm:"unigueIndex"`
}

type SongCreateResponse struct {
	Token string `json:"token"`
}

type SongUpdateRequest struct {
	Group string `json:"group"`
	Name  string `json:"song"`
	Date  string `json:"releaseDate"`
	Text  string `json:"text"`
	Link  string `json:"link"`
	Hash  string `json:"hash"`
}

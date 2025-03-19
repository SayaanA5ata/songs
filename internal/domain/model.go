package domain

import "gorm.io/gorm"

// SongModel представляет модель песни в базе данных.
// swagger:model
type SongModel struct {
	// swagger:ignore
	gorm.Model
	Group string `json:"group"`
	Name  string `json:"song"`
	Date  string `json:"releaseDate"`
	Text  string `json:"text"`
	Link  string `json:"link"`
}

func NewSong(group, name, date, text, link string) *SongModel {
	song := &SongModel{
		Group: group,
		Name:  name,
		Date:  date,
		Text:  text,
		Link:  link,
	}

	return song
}

package domain

import (
	"crypto/rand"
	"math/big"

	"gorm.io/gorm"
)

type SongModel struct {
	gorm.Model
	Group string `json:"group"`
	Name  string `json:"song"`
	Date  string `json:"releaseDate"`
	Text  string `json:"text"`
	Link  string `json:"link"`
	Hash  string `json:"hash" gorm:"uniqueIndex"`
}

func NewSong(group, name, date, text, link string) *SongModel {
	song := &SongModel{
		Group: group,
		Name:  name,
		Date:  date,
		Text:  text,
		Link:  link,
		Hash:  GenerateStringHash(6),
	}
	song.GenerateHash()
	return song
}

func (song *SongModel) GenerateHash() {
	song.Hash = GenerateStringHash(6)
}

func GenerateStringHash(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	charsetRunes := []rune(charset)
	result := make([]rune, length)

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charsetRunes))))
		if err != nil {
			// Обработка ошибки генерации хеша
		}
		result[i] = charsetRunes[randomIndex.Int64()]
	}

	return string(result)
}

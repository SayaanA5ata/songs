package song

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
	Hash  string `json:"hash" gorm:"unigueIndex"`
}

func NewSong(createRequest *SongCreateRequest) *SongModel {
	song := &SongModel{
		Group: createRequest.Group,
		Name:  createRequest.Name,
		Date:  createRequest.Date,
		Text:  createRequest.Text,
		Link:  createRequest.Link,
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
		}
		result[i] = charsetRunes[randomIndex.Int64()]
	}

	return string(result)
}

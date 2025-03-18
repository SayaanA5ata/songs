package domain

import (
	"crypto/rand"
	"log"
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
	log.Printf("Creating a new song with group=%s, name=%s, date=%s", group, name, date)

	song := &SongModel{
		Group: group,
		Name:  name,
		Date:  date,
		Text:  text,
		Link:  link,
		Hash:  GenerateStringHash(6),
	}
	log.Printf("Initial hash generated for the song: %s", song.Hash)

	song.GenerateHash()
	log.Printf("Final hash generated for the song: %s", song.Hash)

	return song
}

func (song *SongModel) GenerateHash() {
	log.Println("Generating a new hash for the song...")
	song.Hash = GenerateStringHash(6)
	log.Printf("New hash generated: %s", song.Hash)
}

func GenerateStringHash(length int) string {
	log.Printf("Generating a random string hash of length %d...", length)

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetRunes := []rune(charset)
	result := make([]rune, length)

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charsetRunes))))
		if err != nil {
			log.Printf("Error generating random index for hash: %v", err)
			// В случае ошибки используем фиксированное значение для избежания паники
			randomIndex = big.NewInt(0)
		}
		result[i] = charsetRunes[randomIndex.Int64()]
	}

	hash := string(result)
	log.Printf("Generated hash: %s", hash)
	return hash
}

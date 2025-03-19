package service

import (
	"dinushc/gorutines/internal/domain"
	"dinushc/gorutines/internal/interfaces"
	"dinushc/gorutines/internal/payload"
	"log"

	"gorm.io/gorm"
)

type SongService struct {
	Repo interfaces.SongRepositoryInterface
}

func NewSongService(repo interfaces.SongRepositoryInterface) *SongService {
	return &SongService{
		Repo: repo,
	}
}

// Получение данных библиотеки с фильтрацией и пагинацией
func (service *SongService) GetSongs(filter map[string]interface{}, page, pageSize int) ([]payload.SongResponse, int64, error) {
	log.Printf("Fetching songs with filter: %+v, page: %d, pageSize: %d", filter, page, pageSize)

	songs, total, err := service.Repo.GetSongs(filter, page, pageSize)
	if err != nil {
		log.Printf("Error fetching songs from repository: %v", err)
		return nil, 0, err
	}
	log.Printf("Successfully fetched %d songs from the repository", len(songs))

	// Преобразуем доменные модели в DTO
	var response []payload.SongResponse
	for _, song := range songs {
		response = append(response, payload.SongResponse{
			ID:    song.ID,
			Group: song.Group,
			Name:  song.Name,
			Date:  song.Date,
			Text:  song.Text,
			Link:  song.Link,
		})
	}
	log.Printf("Converted %d songs to DTO format", len(response))

	return response, total, nil
}

// Получение текста песни с пагинацией по куплетам
func (service *SongService) GetSongVerses(songID uint, page, pageSize int) ([]payload.VerseResponse, int, error) {
	log.Printf("Fetching verses for song ID: %d, page: %d, pageSize: %d", songID, page, pageSize)

	verses, totalVerses, err := service.Repo.GetSongVerses(songID, page, pageSize)
	if err != nil {
		log.Printf("Error fetching verses from repository: %v", err)
		return nil, 0, err
	}
	log.Printf("Successfully fetched %d verses from the repository", len(verses))

	// Преобразуем куплеты в DTO
	var response []payload.VerseResponse
	for i, verse := range verses {
		response = append(response, payload.VerseResponse{
			VerseNumber: i + 1,
			Text:        verse,
		})
	}
	log.Printf("Converted %d verses to DTO format", len(response))

	return response, totalVerses, nil
}

func (service *SongService) CreateSong(createRequest *payload.SongCreateRequest) (*domain.SongModel, error) {
	song := domain.NewSong(createRequest.Group, createRequest.Name, createRequest.Date, createRequest.Text, createRequest.Link)

	createdSong, err := service.Repo.Create(song)
	if err != nil {
		log.Printf("Error creating song in repository: %v", err)
		return nil, err
	}
	log.Printf("Song created successfully with ID: %d", createdSong.ID)

	return createdSong, nil
}

func (service *SongService) UpdateSong(id uint, updateRequest *payload.SongUpdateRequest) (*domain.SongModel, error) {
	log.Printf("Updating song with ID: %d", id)

	song := &domain.SongModel{
		Model: gorm.Model{ID: id},
		Group: updateRequest.Group,
		Name:  updateRequest.Name,
		Date:  updateRequest.Date,
		Text:  updateRequest.Text,
		Link:  updateRequest.Link,
	}

	updatedSong, err := service.Repo.Update(song)
	if err != nil {
		log.Printf("Error updating song in repository: %v", err)
		return nil, err
	}
	log.Printf("Song updated successfully with ID: %d", updatedSong.ID)

	return updatedSong, nil
}

func (service *SongService) DeleteSong(id uint) error {
	log.Printf("Deleting song with ID: %d", id)

	err := service.Repo.Delete(id)
	if err != nil {
		log.Printf("Error deleting song from repository: %v", err)
		return err
	}
	log.Printf("Song deleted successfully with ID: %d", id)

	return nil
}

func (service *SongService) GetSongById(id uint) (*domain.SongModel, error) {
	song, err := service.Repo.GetById(id)
	if err != nil {
		log.Printf("Error fetching song by id from repository: %v", err)
		return nil, err
	}

	return song, nil
}

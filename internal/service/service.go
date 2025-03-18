package service

import (
	"dinushc/gorutines/internal/domain"
	"dinushc/gorutines/internal/interfaces"
	"dinushc/gorutines/internal/payload"

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
	songs, total, err := service.Repo.GetSongs(filter, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

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
			Hash:  song.Hash,
		})
	}

	return response, total, nil
}

// Получение текста песни с пагинацией по куплетам
func (service *SongService) GetSongVerses(songID uint, page, pageSize int) ([]payload.VerseResponse, int, error) {
	verses, totalVerses, err := service.Repo.GetSongVerses(songID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Преобразуем куплеты в DTO
	var response []payload.VerseResponse
	for i, verse := range verses {
		response = append(response, payload.VerseResponse{
			VerseNumber: i + 1,
			Text:        verse,
		})
	}

	return response, totalVerses, nil
}

func (service *SongService) CreateSong(createRequest *payload.SongCreateRequest) (*domain.SongModel, error) {
	song := domain.NewSong(createRequest.Group, createRequest.Name, createRequest.Date, createRequest.Text, createRequest.Link)
	for {
		existedSong, _ := service.Repo.GetByHash(song.Hash)
		if existedSong == nil {
			break
		}
		song.GenerateHash()
	}
	return service.Repo.Create(song)
}

func (service *SongService) UpdateSong(id uint, updateRequest *payload.SongUpdateRequest) (*domain.SongModel, error) {
	song := &domain.SongModel{
		Model: gorm.Model{ID: id},
		Group: updateRequest.Group,
		Name:  updateRequest.Name,
		Date:  updateRequest.Date,
		Text:  updateRequest.Text,
		Link:  updateRequest.Link,
		Hash:  updateRequest.Hash,
	}
	return service.Repo.Update(song)
}

func (service *SongService) DeleteSong(id uint) error {
	return service.Repo.Delete(id)
}

func (service *SongService) GetSongByHash(hash string) (*domain.SongModel, error) {
	return service.Repo.GetByHash(hash)
}

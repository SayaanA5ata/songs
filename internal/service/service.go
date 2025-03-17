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

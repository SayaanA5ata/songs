package interfaces

import "dinushc/gorutines/internal/domain"

type SongRepositoryInterface interface {
	GetSongs(filter map[string]interface{}, page, pageSize int) ([]domain.SongModel, int64, error)
	GetSongVerses(songID uint, page, pageSize int) ([]string, int, error)
	Create(song *domain.SongModel) (*domain.SongModel, error)
	Update(song *domain.SongModel) (*domain.SongModel, error)
	GetById(id uint) (*domain.SongModel, error)
	Delete(id uint) error
}

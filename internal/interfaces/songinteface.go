package interfaces

import "dinushc/gorutines/internal/domain"

type SongRepositoryInterface interface {
	Create(song *domain.SongModel) (*domain.SongModel, error)
	GetByHash(hash string) (*domain.SongModel, error)
	Update(song *domain.SongModel) (*domain.SongModel, error)
	GetById(id uint) (*domain.SongModel, error)
	Delete(id uint) error
}

package implementation

import (
	"dinushc/gorutines/internal/domain"
	"dinushc/gorutines/pkg/db"

	"gorm.io/gorm/clause"
)

type SongRepository struct {
	Database *db.Db
}

func NewSongRepository(database *db.Db) *SongRepository {
	return &SongRepository{
		Database: database,
	}
}

func (repo *SongRepository) Create(song *domain.SongModel) (*domain.SongModel, error) {
	result := repo.Database.DB.Create(song)
	if result.Error != nil {
		return nil, result.Error
	}
	return song, nil
}

func (repo *SongRepository) GetByHash(hash string) (*domain.SongModel, error) {
	var song domain.SongModel
	result := repo.Database.DB.First(&song, "hash=?", hash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &song, nil
}

func (repo *SongRepository) Update(song *domain.SongModel) (*domain.SongModel, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(song)
	if result.Error != nil {
		return nil, result.Error
	}
	return song, nil
}

func (repo *SongRepository) GetById(id uint) (*domain.SongModel, error) {
	var song domain.SongModel
	result := repo.Database.DB.First(&song, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &song, nil
}

func (repo *SongRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&domain.SongModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

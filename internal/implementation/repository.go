package implementation

import (
	"dinushc/gorutines/internal/domain"
	"dinushc/gorutines/pkg/db"
	"strings"

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

// Получение данных библиотеки с фильтрацией и пагинацией
func (repo *SongRepository) GetSongs(filter map[string]interface{}, page, pageSize int) ([]domain.SongModel, int64, error) {
	var songs []domain.SongModel
	var total int64

	query := repo.Database.DB.Model(&domain.SongModel{})

	// Применяем фильтры
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	// Подсчитываем общее количество записей
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Применяем пагинацию
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&songs).Error; err != nil {
		return nil, 0, err
	}

	return songs, total, nil
}

// Получение текста песни с пагинацией по куплетам
func (repo *SongRepository) GetSongVerses(songID uint, page, pageSize int) ([]string, int, error) {
	var song domain.SongModel
	if err := repo.Database.DB.First(&song, songID).Error; err != nil {
		return nil, 0, err
	}

	// Разделяем текст на куплеты
	verses := splitTextIntoVerses(song.Text)
	totalVerses := len(verses)

	// Применяем пагинацию
	offset := (page - 1) * pageSize
	end := offset + pageSize
	if end > totalVerses {
		end = totalVerses
	}

	return verses[offset:end], totalVerses, nil
}

// Вспомогательная функция для разделения текста на куплеты
func splitTextIntoVerses(text string) []string {
	// Разделяем текст по символу новой строки
	verses := strings.Split(text, "\n\n")
	return verses
}

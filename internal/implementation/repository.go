package implementation

import (
	"dinushc/gorutines/internal/domain"
	"dinushc/gorutines/pkg/db"
	"log"
	"strings"

	"gorm.io/gorm/clause"
)

type SongRepository struct {
	Database *db.Db
}

func NewSongRepository(database *db.Db) *SongRepository {
	log.Println("Initializing SongRepository...")
	return &SongRepository{
		Database: database,
	}
}

func (repo *SongRepository) Create(song *domain.SongModel) (*domain.SongModel, error) {
	log.Printf("Creating a new song in the database: %+v", song)
	result := repo.Database.DB.Create(song)
	if result.Error != nil {
		log.Printf("Error creating song: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Song created successfully with ID: %d", song.ID)
	return song, nil
}

func (repo *SongRepository) GetByHash(hash string) (*domain.SongModel, error) {
	log.Printf("Fetching song by hash: %s", hash)
	var song domain.SongModel
	result := repo.Database.DB.First(&song, "hash=?", hash)
	if result.Error != nil {
		log.Printf("Error fetching song by hash: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Song fetched successfully by hash: %s", hash)
	return &song, nil
}

func (repo *SongRepository) Update(song *domain.SongModel) (*domain.SongModel, error) {
	log.Printf("Updating song with ID: %d", song.ID)
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(song)
	if result.Error != nil {
		log.Printf("Error updating song: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Song updated successfully with ID: %d", song.ID)
	return song, nil
}

func (repo *SongRepository) GetById(id uint) (*domain.SongModel, error) {
	log.Printf("Fetching song by ID: %d", id)
	var song domain.SongModel
	result := repo.Database.DB.First(&song, id)
	if result.Error != nil {
		log.Printf("Error fetching song by ID: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Song fetched successfully by ID: %d", id)
	return &song, nil
}

func (repo *SongRepository) Delete(id uint) error {
	log.Printf("Deleting song with ID: %d", id)
	result := repo.Database.DB.Delete(&domain.SongModel{}, id)
	if result.Error != nil {
		log.Printf("Error deleting song: %v", result.Error)
		return result.Error
	}
	log.Printf("Song deleted successfully with ID: %d", id)
	return nil
}

// Получение данных библиотеки с фильтрацией и пагинацией
func (repo *SongRepository) GetSongs(filter map[string]interface{}, page, pageSize int) ([]domain.SongModel, int64, error) {
	log.Printf("Fetching songs with filter: %+v, page: %d, pageSize: %d", filter, page, pageSize)

	var songs []domain.SongModel
	var total int64

	query := repo.Database.DB.Model(&domain.SongModel{})

	// Применяем фильтры
	for key, value := range filter {
		log.Printf("Applying filter: %s = %v", key, value)
		query = query.Where(key+" = ?", value)
	}

	// Подсчитываем общее количество записей
	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting songs: %v", err)
		return nil, 0, err
	}
	log.Printf("Total songs found: %d", total)

	// Применяем пагинацию
	offset := (page - 1) * pageSize
	log.Printf("Applying pagination with offset: %d, limit: %d", offset, pageSize)
	if err := query.Offset(offset).Limit(pageSize).Find(&songs).Error; err != nil {
		log.Printf("Error fetching paginated songs: %v", err)
		return nil, 0, err
	}

	log.Printf("Fetched %d songs successfully", len(songs))
	return songs, total, nil
}

// Получение текста песни с пагинацией по куплетам
func (repo *SongRepository) GetSongVerses(songID uint, page, pageSize int) ([]string, int, error) {
	log.Printf("Fetching verses for song ID: %d, page: %d, pageSize: %d", songID, page, pageSize)

	var song domain.SongModel
	if err := repo.Database.DB.First(&song, songID).Error; err != nil {
		log.Printf("Error fetching song by ID: %v", err)
		return nil, 0, err
	}
	log.Printf("Song fetched successfully by ID: %d", songID)

	// Разделяем текст на куплеты
	verses := splitTextIntoVerses(song.Text)
	totalVerses := len(verses)
	log.Printf("Split song text into %d verses", totalVerses)

	// Применяем пагинацию
	offset := (page - 1) * pageSize
	end := offset + pageSize
	if end > totalVerses {
		end = totalVerses
	}
	log.Printf("Applying pagination to verses with offset: %d, limit: %d", offset, pageSize)

	return verses[offset:end], totalVerses, nil
}

// Вспомогательная функция для разделения текста на куплеты
func splitTextIntoVerses(text string) []string {
	log.Println("Splitting song text into verses...")
	verses := strings.Split(text, "\n\n")
	log.Printf("Split text into %d verses", len(verses))
	return verses
}

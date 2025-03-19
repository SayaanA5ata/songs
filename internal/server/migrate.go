package server

import (
	"dinushc/gorutines/internal/domain"
	"dinushc/gorutines/pkg/db"
	"log"

	"gorm.io/gorm"
)

func MigrateIfNotExists(dbase db.Db) {
	// Проверка существования таблицы перед выполнением миграции
	log.Println("Checking if the table exists...")
	tableExists := dbase.Migrator().HasTable(&domain.SongModel{})
	if !tableExists {
		log.Println("Table does not exist. Running database migrations...")
		err := dbase.AutoMigrate(&domain.SongModel{})
		if err != nil {
			log.Fatalf("Error running database migrations: %v", err)
		}
		log.Println("Database migrations completed successfully.")
	} else {
		log.Println("Table already exists. Skipping migrations.")
	}

	// Проверка наличия данных в таблице
	var songCount int64
	if err := dbase.Model(&domain.SongModel{}).Count(&songCount).Error; err != nil {
		log.Fatalf("Error counting records in the table: %v", err)
	}

	if songCount == 0 {
		log.Println("Table is empty. Populating with test data...")
		seedTestData(dbase.DB)
		log.Println("Test data populated successfully.")
	} else {
		log.Println("Table already contains data. Skipping seeding.")
	}
}

// Функция для наполнения таблицы тестовыми данными
func seedTestData(db *gorm.DB) {
	testData := []domain.SongModel{
		{
			Group: "Muse",
			Name:  "Supermassive Black Hole",
			Date:  "16.07.2006",
			Text:  "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?",
			Link:  "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		},
		{
			Group: "Coldplay",
			Name:  "Yellow",
			Date:  "26.06.2000",
			Text:  "Look at the stars,\nLook how they shine for you,\nAnd everything you do,\nYeah, they were all yellow.",
			Link:  "https://www.youtube.com/watch?v=yKNxeF4KMsY",
		},
		{
			Group: "Radiohead",
			Name:  "Creep",
			Date:  "21.09.1992",
			Text:  "When you were here before,\nCouldn't look you in the eye,\nYou're just like an angel,\nYour skin makes me cry.",
			Link:  "https://www.youtube.com/watch?v=XFkzRNyygfk",
		},
	}

	// Добавление тестовых данных в таблицу
	for _, song := range testData {
		if err := db.Create(&song).Error; err != nil {
			log.Printf("Error inserting test data: %v", err)
		}
	}
}

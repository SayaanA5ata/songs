package main

import (
	"dinushc/gorutines/internal/domain"
	"dinushc/gorutines/pkg/dsn"

	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Загрузка .env файла
	log.Println("Loading .env file...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Successfully loaded .env file.")

	// Подключение к базе данных
	log.Println("Connecting to the database...")
	db, err := gorm.Open(postgres.Open(dsn.GetPureDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	log.Println("Successfully connected to the database.")

	// Выполнение миграций
	log.Println("Running database migrations...")
	err = db.AutoMigrate(&domain.SongModel{})
	if err != nil {
		log.Fatalf("Error running database migrations: %v", err)
	}
	log.Println("Database migrations completed successfully.")
}

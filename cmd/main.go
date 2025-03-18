// @title Test API
// @version 1.0
// @description It's a simple API for songs.
// @host localhost:8081
// @BasePath /

package main

import (
	"dinushc/gorutines/configs"
	"dinushc/gorutines/internal/server"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка конфигурации
	log.Println("Loading .env file...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Successfully loaded .env file.")

	conf := configs.LoadConfig()
	if conf == nil {
		log.Fatal("Failed to load configuration")
	}
	log.Println("Configuration loaded successfully.")

	// Инициализация сервера
	srv := server.NewServer(conf)
	srv.Run()
}

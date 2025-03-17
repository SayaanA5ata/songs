package main

import (
	"context"
	"dinushc/gorutines/configs"
	"dinushc/gorutines/internal/song"
	"dinushc/gorutines/pkg/db"
	"dinushc/gorutines/pkg/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка конфигурации
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	conf := configs.LoadConfig()
	if conf == nil {
		log.Fatal("Failed to load configuration")
	}

	dbase := db.NewDb(conf)
	if dbase == nil {
		log.Fatal("Failed to initialize database")
	}

	// Создаем Chi-роутер
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.CORS)
	router.Use(middleware.Logging)

	// Repositories
	songRepository := song.NewSongRepository(dbase)

	// Handlers
	song.NewSongHandler(router, song.SongHandlerDeps{
		SongRepo: songRepository,
	})

	// Определение порта
	port := ":" + os.Getenv("SERVER_PORT")
	if port == ":" {
		port = "8080" // Используем порт по умолчанию
	}

	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	// Обработка сигналов завершения
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Запуск сервера в горутине
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Гарантированное завершение
	defer func() {
		<-ctx.Done() // Ждем сигнала завершения

		// Корректное завершение HTTP-сервера
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	log.Printf("Server started on port %s", port)
	<-ctx.Done()
	log.Println("Application shutdown gracefully.")
}

package server

import (
	"context"
	"dinushc/gorutines/configs"
	"dinushc/gorutines/internal/handlers"
	"dinushc/gorutines/internal/implementation"
	"dinushc/gorutines/internal/service"
	"dinushc/gorutines/pkg/db"
	"dinushc/gorutines/pkg/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Config *configs.Config
}

func NewServer(config *configs.Config) *Server {
	return &Server{
		Config: config,
	}
}

func (s *Server) Run() {
	log.Println("Initializing database connection...")
	dbase := db.NewDb(s.Config)
	if dbase == nil {
		log.Fatal("Failed to initialize database")
	}
	log.Println("Database connection initialized successfully.")

	// Создаем Chi-роутер
	log.Println("Setting up Chi router...")
	router := chi.NewRouter()

	// Middleware
	log.Println("Applying middleware...")
	router.Use(middleware.CORS)
	router.Use(middleware.Logging)

	// Repositories
	log.Println("Initializing song repository...")
	songRepository := implementation.NewSongRepository(dbase)

	// Service
	log.Println("Initializing song service...")
	songService := service.NewSongService(songRepository)

	// Handlers
	log.Println("Registering song handlers...")
	handlers.NewSongHandler(router, songService)

	// Определение порта
	port := ":" + os.Getenv("SERVER_PORT")
	if port == ":" {
		port = ":8080" // Используем порт по умолчанию
		log.Println("No SERVER_PORT found in .env, using default port 8080.")
	} else {
		log.Printf("Using SERVER_PORT from .env: %s", port)
	}

	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	// Обработка сигналов завершения
	log.Println("Setting up signal handling for graceful shutdown...")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Запуск сервера в горутине
	log.Printf("Starting HTTP server on port %s...", port)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
		log.Println("HTTP server stopped.")
	}()

	// Гарантированное завершение
	log.Println("Waiting for shutdown signal...")
	defer func() {
		<-ctx.Done() // Ждем сигнала завершения
		log.Println("Shutdown signal received. Initiating graceful shutdown...")

		// Корректное завершение HTTP-сервера
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		} else {
			log.Println("HTTP server shutdown completed successfully.")
		}
	}()

	log.Printf("Application is running on port %s. Press Ctrl+C to exit.", port)
	<-ctx.Done()
	log.Println("Application shutdown gracefully.")
}

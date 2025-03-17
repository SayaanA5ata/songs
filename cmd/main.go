package main

import (
	"dinushc/gorutines/configs"
	"dinushc/gorutines/internal/song"
	"dinushc/gorutines/pkg/db"
	"dinushc/gorutines/pkg/middleware"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("error")
	}
	conf := configs.LoadConfig()
	dbase := db.NewDb(conf)
	// Создаем Chi-роутер
	router := chi.NewRouter()
	// Middleware
	router.Use(middleware.CORS)
	router.Use(middleware.Logging)
	//repositories
	songRepository := song.NewSongRepository(dbase)
	//handler
	song.NewSongHandler(router, song.SongHandlerDeps{
		SongRepo: songRepository,
	})
	port := os.Getenv("SERVER_PORT")
	port = ":" + port
	server := http.Server{
		Addr:    port,
		Handler: router,
	}
	server.ListenAndServe()
}
